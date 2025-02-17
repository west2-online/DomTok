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
)

// sf可以生成id,详见user/domain/service/service.go
// TODO 这个函数的逻辑不知道要怎么写，我只知道大概要包括生成订单信息、往sql里存信息、返回支付id这三步
func (svc *PaymentService) CreatePaymentInfo(ctx context.Context, orderID int64) (paymentID int64, err error) {
	// TODO 2025.02.17 00：22 先来解决你
	// 1. 生成支付 ID（雪花算法）
	paymentID, err = svc.sf.NextVal()
	if err != nil {
		return 0, fmt.Errorf("failed to create payment information order: %w", err)
	}

	// 2. 构造支付订单对象
	paymentOrder := &model.PaymentOrder{
		ID:      paymentID,
		OrderID: orderID,
		Status:  paymentStatus.PaymentStatusPending, // 设定初始状态
	}

	// 3. 存入数据库
	err = svc.db.CreatePayment(ctx, paymentOrder)
	if err != nil {
		return 0, fmt.Errorf("failed to create payment order: %w", err)
	}

	// 4. 返回支付 ID
	return paymentID, nil
}

// TODO 这个也要向User模块发起数据库查询申请
func (svc *PaymentService) CheckUserExist(ctx context.Context, uid int64) (userInfo bool, err error) {
	return paymentStatus.UserNotExist, nil
}

// GetUserID 等User模块完成了再写这个，从ctx里获取userID
func (svc *PaymentService) GetUserID(ctx context.Context) (uid int64, err error) {
	uid, err = loginData.GetLoginData(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get login data: %w", err)
	}
	return uid, nil
}

// TODO 后面完善这个接口，要发起RPC请求向order模块申请数据库的查询，所以后面再来写
func (svc *PaymentService) CheckOrderExist(ctx context.Context, orderID int64) (orderInfo bool, err error) {
	return paymentStatus.OrderNotExist, nil
}

// GeneratePaymentToken HMAC生成支付令牌
func (svc *PaymentService) GeneratePaymentToken(ctx context.Context, orderID int64) (string, int64, error) {
	// 1. 设定过期时间为15分钟后, 即现在时间加上15分钟之后的秒级时间戳
	expirationTime := time.Now().Add(paymentStatus.ExpirationDuration).Unix()
	// 2. 获取 HMAC 密钥（可以从环境变量或配置文件获取）
	secretKey := []byte(paymentStatus.PaymentSecretKey)

	// 3. 计算 HMAC-SHA256 哈希
	h := hmac.New(sha256.New, secretKey)
	_, err := h.Write([]byte(fmt.Sprintf("%d:%d", orderID, expirationTime)))
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate HMAC: %w", err)
	}

	// 4. 生成十六进制编码的 HMAC 签名
	token := hex.EncodeToString(h.Sum(nil))

	// 5. 返回令牌和过期时间
	// TODO
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
		return paymentStatus.RedisStoreFailed, fmt.Errorf("cannot store token: expiration time has already passed")
	}
	// 2. 构造 Redis Key
	redisKey := fmt.Sprintf("payment_token:%d:%d", userID, orderID)
	// 3. 存储到 Redis
	err := svc.redis.SetPaymentToken(ctx, redisKey, token, expirationDuration)
	if err != nil {
		return paymentStatus.RedisStoreFailed, fmt.Errorf("failed to store payment token in redis: %w", err)
	}
	// 4. 返回成功状态码
	return paymentStatus.RedisStoreSuccess, nil

}
