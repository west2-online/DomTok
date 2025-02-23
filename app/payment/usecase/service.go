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

package usecase

import (
	"context"
	"fmt"

	"github.com/west2-online/DomTok/app/payment/domain/model"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
)

// CreatePayment 这里定义一些具体的方法和函数，比如校验密码，加密密码，创建用户之类的
func (uc *paymentUseCase) CreatePayment(ctx context.Context, orderID int64) (*model.PaymentOrder, error) {
	return nil, nil
}

func (uc *paymentUseCase) GetPaymentToken(ctx context.Context, orderID int64) (token string, expTime int64, err error) {
	// 1. 检查订单是否存在
	// TODO 这个要向order模块要一个RPC接口然后再来填充
	/*var orderInfo bool
	orderInfo, err = uc.svc.CheckOrderExist(ctx, orderID)
	if err != nil {
		return "", 0, fmt.Errorf("check order existed failed:%w", err)
	}
	if orderInfo == paymentStatus.OrderNotExist {
		return "", 0, errno.NewErrNo(errno.PaymentOrderNotExist, "order does not exist")
	}
	*/
	// 2. 获取用户id,无需检查用户是否存在
	// 获取用户id
	var uid int64
	uid, err = uc.svc.GetUserID(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("get user id failed:%w", err)
	}
	/*
		// 3. 检查订单支付信息
		var paymentInfo bool
		paymentInfo, err = uc.db.CheckPaymentExist(ctx, orderID)
		if err != nil {
			return "", 0, fmt.Errorf("check payment existed failed:%w", err)
		}
		if paymentInfo == paymentStatus.PaymentNotExist { // 如果订单不存在
			// 创建支付订单
			// TODO 待完善
			_, err := uc.svc.CreatePaymentInfo(ctx, orderID)
			if err != nil {
				return "", 0, fmt.Errorf("create payment info failed:%w", err)
			}
		} else if paymentInfo == paymentStatus.PaymentExist { // 如果订单存在
			// 获取订单的支付状态
			payStatus, err := uc.db.GetPaymentInfo(ctx, orderID)
			// 如果订单正在支付或者已经支付完成，则拒绝进行接下来的生成令牌的活动
			if payStatus == paymentStatus.PaymentStatusSuccessCode || payStatus == paymentStatus.PaymentStatusProcessingCode {
				return "", 0, fmt.Errorf("payment is processing or has already done:%w", err)
			}
		}
	*/
	// 4. HMAC生成支付令牌
	token, expTime, err = uc.svc.GeneratePaymentToken(ctx, orderID)
	if err != nil {
		logger.Errorf("Error generating payment token: orderID:%d,err:%v", orderID, err)
		return "", 0, fmt.Errorf("generate payment token failed:%w", err)
	}
	var redisStatus bool
	// 5. 存储令牌到 Redis
	redisStatus, err = uc.svc.StorePaymentToken(ctx, token, expTime, uid, orderID)
	if err != nil && redisStatus != paymentStatus.RedisStoreSuccess {
		logger.Errorf("Error store payment token: orderID:%d,userID:%d,err:%v", orderID, uid, err)
		return "", 0, fmt.Errorf("store payment token failed:%w", err)
	}
	logger.Infof("Success generating payment token: orderID:%d,token:%s", orderID, token)
	return token, expTime, nil
}

// GetRefundInfo 获取退款信息 TODO
func (uc *paymentUseCase) GetRefundInfo(ctx context.Context, orderID int64) (refundID int64, err error) {
	/*
		// 1. 检查订单是否存在
		orderExists, err := uc.svc.CheckOrderExist(ctx, orderID)
		if err != nil {
			return "", 0, fmt.Errorf("check order existence failed: %w", err)
		}
		if !orderExists {
			return "", 0, errno.NewErrNo(errno.PaymentOrderNotExist, "order does not exist")
		}
	*/
	// 2. 获取用户ID
	uid, err := uc.svc.GetUserID(ctx)
	if err != nil {
		return 0, fmt.Errorf("get user id failed: %w", err)
	}
	// 3. Redis 限流检查
	var frequencyInfo bool
	var timeInfo bool
	frequencyInfo, timeInfo, err = uc.svc.CheckRedisRateLimiting(ctx, uid, orderID)
	if err != nil {
		return 0, fmt.Errorf("check redis rate limiting failed: %w", err)
	}
	if frequencyInfo != paymentStatus.RedisValid {
		return 0, fmt.Errorf("too many refund requests in a short time")
	}
	if timeInfo != paymentStatus.RedisValid {
		return 0, fmt.Errorf("refund already requested for this order in the last 24 hours")
	}

	// 4. 创建退款信息
	refundID, err = uc.svc.CreateRefundInfo(ctx, orderID)
	if err != nil {
		return 0, fmt.Errorf("create refund info failed: %w", err)
	}
	/*
		// 5. 生成退款令牌
		token, expTime, err = uc.svc.GenerateRefundToken(ctx, orderID)
		if err != nil {
			return "", 0, fmt.Errorf("generate refund token failed: %w", err)
		}

		// 6. 存储令牌到 Redis
		var redisStatus bool
		redisStatus, err = uc.svc.StoreRefundToken(ctx, token, expTime, uid, orderID)
		if err != nil || redisStatus != paymentStatus.RedisStoreSuccess {
			return "", 0, fmt.Errorf("store refund token failed: %w", err)
		}

		return token, expTime, nil
	*/
	return refundID, nil
}
