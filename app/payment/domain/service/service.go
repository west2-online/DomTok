/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/west2-online/DomTok/app/payment/domain/model"
	loginData "github.com/west2-online/DomTok/pkg/base/context"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

// CreatePaymentInfo sf可以生成id,详见user/domain/service/service.go
func (svc *PaymentService) CreatePaymentInfo(ctx context.Context, orderID int64) (paymentID int64, err error) {
	// 1. 生成支付 ID（雪花算法）
	paymentID, err = svc.sf.NextVal()
	if err != nil {
		return 0, errno.Errorf(errno.ServiceCreatePaymentIDFailed, "failed to create payment ID: %v", err)
	}

	// 2. 构造支付订单对象
	paymentOrder := &model.PaymentOrder{
		ID:      paymentID,
		OrderID: orderID,
		Status:  paymentStatus.PaymentStatusPendingCode, // 设定初始状态
	}

	// 3. 存入数据库
	err = svc.db.CreatePayment(ctx, paymentOrder)
	if err != nil {
		return 0, errno.Errorf(errno.ServiceCreatePaymentOrderFailed, "failed to create payment order: %v", err)
	}

	// 4. 返回支付 ID
	return paymentID, nil
}

// GetUserID 等User模块完成了再写这个，从ctx里获取userID
func (svc *PaymentService) GetUserID(ctx context.Context) (uid int64, err error) {
	uid, err = loginData.GetLoginData(ctx)
	if err != nil {
		return 0, errno.Errorf(errno.ServiceGetLoginDataFailed, "failed to get login data: %v", err)
	}
	return uid, nil
}

// CheckOrderExist 检查订单是否存在（调用Order模块的接口）
func (svc *PaymentService) CheckOrderExist(ctx context.Context, orderID int64) (orderInfo bool, err error) {
	userInfo, err := svc.rpc.PaymentIsOrderExist(ctx, orderID)
	if err != nil {
		return false, errno.Errorf(errno.ServiceCheckOrderExistFailed, "failed to check order existence: %v", err)
	}
	return userInfo, nil
}

// GeneratePaymentToken HMAC生成支付令牌
func (svc *PaymentService) GeneratePaymentToken(ctx context.Context, orderID int64) (token string, expirationTime int64, err error) {
	// 1. 设定过期时间为15分钟后, 即现在时间加上15分钟之后的秒级时间戳
	expirationTime = time.Now().Add(paymentStatus.ExpirationDuration).Unix()
	logger.Infof("Generating payment token, orderID: %d, expirationTime: %d", orderID, expirationTime)
	// 2. 获取 HMAC 密钥（可以从环境变量或配置文件获取）
	secretKey := []byte(paymentStatus.PaymentSecretKey)

	// 3. 计算 HMAC-SHA256 哈希
	h := hmac.New(sha256.New, secretKey)
	_, err = h.Write([]byte(fmt.Sprintf("%d:%d", orderID, expirationTime)))
	if err != nil {
		logger.Infof("failed to generate payment HMAC token, orderID: %d, expirationTime: %d", orderID, expirationTime)
		return "", 0, errno.Errorf(errno.ServiceGenerateHMACFailed, "failed to generate payment HMAC token: %v", err)
	}

	// 4. 生成十六进制编码的 HMAC 签名
	token = hex.EncodeToString(h.Sum(nil))
	// logger.Infof("Generated payment HAMC token successfully, orderID: %d, expirationTime: %d", orderID, expirationTime)
	// 5. 返回令牌和过期时间
	return token, expirationTime, nil
}

// StorePaymentToken 这里的返回值还没有想好，是返回状态码还是消息字段？
func (svc *PaymentService) StorePaymentToken(ctx context.Context, token string, expTime int64, userID int64, orderID int64) (bool, error) {
	// 1. 计算剩余过期时间
	// 这个expiration是expTime减去当前时间，得到的是过期剩余时间（如900s）
	// 这样可以防止“直接用paymentStatus.ExpirationTime存redis的参数的话，
	// 如果StorePaymentToken执行时expTime早就过期了，仍然会存15min”的bug
	// TODO 我不知道是不是这样的，因为我感觉两个函数执行时间基本上只差几十毫秒，不可能出现这样的情况吧，但想想又有道理
	expirationDuration := time.Until(time.Unix(expTime, 0))
	if expirationDuration <= 0 {
		logger.Warnf("Token expiration time has already passed: orderID: %d, userID: %d", orderID, userID)
		return paymentStatus.RedisStoreFailed, errno.Errorf(errno.ServiceStorePaymentRedisTokenFailed, "cannot store token: expiration time has already passed")
	}
	// 2. 构造 Redis Key
	redisKey := fmt.Sprintf("payment_token:%d:%d", userID, orderID)
	// 3. 存储到 Redis
	err := svc.redis.SetPaymentToken(ctx, redisKey, token, expirationDuration)
	if err != nil {
		logger.Infof("failed to store payment token in redis, orderID: %d, userID: %d", orderID, userID)
		return paymentStatus.RedisStoreFailed, errno.Errorf(errno.InternalRedisErrorCode, "failed to store payment token in redis: %v", err)
	}
	// 4. 返回成功状态码
	return paymentStatus.RedisStoreSuccess, nil
}

func (svc *PaymentService) CheckRedisRateLimiting(ctx context.Context, uid int64, orderID int64) (frequencyInfo bool, timeInfo bool, err error) {
	minuteKey := fmt.Sprintf("userID:%d:refund:minute", uid)
	dayKey := fmt.Sprintf("orderID:%d:refund:day", orderID)

	// 检查 1 分钟内的申请次数
	count, err := svc.redis.IncrRedisKey(ctx, minuteKey, paymentStatus.RedisMinute)
	if err != nil {
		return false, false, errno.Errorf(errno.InternalRedisErrorCode, "check refund request limit failed: %v", err)
	}
	if count > paymentStatus.RedisCheckTimesInMinute {
		return false, false, errno.Errorf(errno.ServiceRedisTimeLimited, "too many refund requests in a short time")
	}

	// 检查 24 小时内是否已申请过退款
	// TODO 这下面的返回值的全用false吗？
	exists, err := svc.redis.CheckRedisDayKey(ctx, dayKey)
	if err != nil {
		return false, false, errno.Errorf(errno.InternalRedisErrorCode, "check refund request history failed: %v", err)
	}
	if exists {
		return true, false, errno.Errorf(errno.ServiceRedisTimeLimited, "refund already requested for this order in the last 24 hours")
	}
	// 记录订单退款请求，设置 24 小时过期
	err = svc.redis.SetRedisDayKey(ctx, dayKey, paymentStatus.RedisDayPlaceholder, paymentStatus.RedisDay)
	if err != nil {
		return false, false, errno.Errorf(errno.InternalRedisErrorCode, "record refund request failed: %v", err)
	}
	return true, true, nil
}

func (svc *PaymentService) CreateRefundInfo(ctx context.Context, orderID int64) (refundID int64, err error) {
	// 1. 生成退款 ID（雪花算法）
	refundID, err = svc.sf.NextVal()
	if err != nil {
		return 0, errno.Errorf(errno.ServiceCreateRefundIDFailed, "failed to create refund information order: %v", err)
	}

	// 2. 构造退款订单对象
	refundOrder := &model.PaymentRefund{
		ID:      refundID,
		OrderID: orderID,
		Status:  paymentStatus.RefundStatusPendingCode, // 设定初始状态
	}

	// 3. 存入数据库
	err = svc.db.CreateRefund(ctx, refundOrder)
	if err != nil {
		return 0, errno.Errorf(errno.ServiceCreateRefundFailed, "failed to create refund order: %v", err)
	}

	// 4. 返回退款 ID
	return refundID, nil
}

func (svc *PaymentService) GetPaymentStatusMsg(code int8) string {
	return paymentStatus.GetPaymentStatus(code)
}

func (svc *PaymentService) GetRefundStatusMsg(code int8) string {
	return paymentStatus.GetRefundStatus(code)
}

/*func (svc *PaymentService) GenerateRefundToken(ctx context.Context, refundID int64) (string, int64, error) {
	// 1. 设定过期时间（15 分钟后）
	expirationTime := time.Now().Add(paymentStatus.RefundExpirationDuration).Unix()
	logger.Infof("Generating refund token, refundID: %d, expirationTime: %d", refundID, expirationTime)

	// 2. 获取 HMAC 密钥
	secretKey := []byte(paymentStatus.RefundSecretKey)

	// 3. 计算 HMAC-SHA256 哈希
	h := hmac.New(sha256.New, secretKey)
	_, err := h.Write([]byte(fmt.Sprintf("%d:%d", refundID, expirationTime)))
	if err != nil {
		logger.Infof("failed to generate refund HMAC token, refundID: %d, expirationTime: %d", refundID, expirationTime)
		return "", 0, fmt.Errorf("failed to generate refund HMAC token: %w", err)
	}

	// 4. 生成十六进制编码的 HMAC 签名
	token := hex.EncodeToString(h.Sum(nil))
	// logger.Infof("Generated refund HMAC token successfully, refundID: %d, expirationTime: %d", refundID, expirationTime)

	// 5. 返回令牌和过期时间
	return token, expirationTime, nil
}

func (svc *PaymentService) StoreRefundToken(ctx context.Context, token string, expTime int64, userID int64, refundID int64) (bool, error) {
	// 1. 计算剩余过期时间
	expirationDuration := time.Until(time.Unix(expTime, 0))
	if expirationDuration <= 0 {
		logger.Warnf("Token expiration time has already passed: refundID: %d, userID: %d", refundID, userID)
		return paymentStatus.RedisStoreFailed, fmt.Errorf("cannot store refund token: expiration time has already passed")
	}

	// 2. 构造 Redis Key
	redisKey := fmt.Sprintf("refund_token:%d:%d", userID, refundID)

	// 3. 存储到 Redis
	err := svc.redis.SetRefundToken(ctx, redisKey, token, expirationDuration)
	if err != nil {
		logger.Infof("failed to store refund token in redis, refundID: %d, userID: %d", refundID, userID)
		return paymentStatus.RedisStoreFailed, fmt.Errorf("failed to store refund token in redis: %w", err)
	}

	// 4. 返回成功状态码
	return paymentStatus.RedisStoreSuccess, nil
}
*/

func (svc *PaymentService) CheckAdminPermission(_ context.Context, uid int64) (bool, error) {
	return uid == 1, nil
}

func (svc *PaymentService) CheckAndDelPaymentToken(ctx context.Context, token string, userID int64, orderID int64) (bool, error) {
	result, err := svc.redis.CheckAndDelPaymentToken(ctx, fmt.Sprintf("payment_token:%d:%d", userID, orderID), token)
	if err != nil {
		return false, fmt.Errorf("failed to check and delete payment token: %w", err)
	}
	return result, nil
}

func (svc *PaymentService) GetExpiredAtAndDelPaymentToken(ctx context.Context,
	token string, userId int64, orderID int64,
) (exist bool, exp time.Time, err error) {
	exist, ttl, err := svc.redis.GetTTLAndDelPaymentToken(ctx, fmt.Sprintf("payment_token:%d:%d", userId, orderID), token)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("failed to get and delete payment token: %w", err)
	}
	return exist, time.Now().Add(ttl), nil
}

func (svc *PaymentService) PutBackPaymentToken(ctx context.Context, token string, userID int64, orderID int64, exp time.Time) error {
	return svc.redis.SetPaymentToken(ctx, fmt.Sprintf("payment_token:%d:%d", userID, orderID), token, time.Until(exp))
}

func (svc *PaymentService) GetOrderStatus(ctx context.Context, orderID int64) (bool, bool, error) {
	exist, expire, err := svc.rpc.GetOrderStatus(ctx, orderID)
	if err != nil {
		return false, true, fmt.Errorf("failed to get order status: %w", err)
	}
	return exist, time.Now().UnixMilli() > expire, nil
}

// GetPayInfo 模拟获取支付信息
func (svc *PaymentService) GetPayInfo(_ context.Context) (int64, string, error) {
	return time.Now().UnixMilli(), paymentStatus.PaymentStyleDomTok, nil
}

// Pay 模拟支付
func (svc *PaymentService) Pay(_ context.Context) (int64, string, error) {
	return time.Now().UnixMilli(), paymentStatus.PaymentStyleDomTok, nil
}

// Refund 模拟退款
func (svc *PaymentService) Refund(_ context.Context) (int64, string, error) {
	return time.Now().UnixMilli(), paymentStatus.PaymentStyleDomTok, nil
}

func (svc *PaymentService) CancelOrder(ctx context.Context, orderID int64, paymentAt int64, paymentStyle string) error {
	return svc.rpc.OrderPaymentCancel(ctx, orderID, paymentAt, paymentStyle)
}

func (svc *PaymentService) ConfirmOrder(ctx context.Context, orderID int64, paymentAt int64, paymentStyle string) error {
	return svc.rpc.OrderPaymentSuccess(ctx, orderID, paymentAt, paymentStyle)
}
