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

	"go.uber.org/zap"

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
	fmt.Println("GetPaymentToken")
	// TODO 这个要向order模块要一个RPC接口然后再来填充
	/*var orderInfo bool
	orderInfo, err = uc.svc.CheckOrderExist(ctx, orderID)
	if err != nil {
		return "", 0, fmt.Errorf("check order existed failed:%w", err)
	}
	if orderInfo == paymentStatus.OrderNotExist {
		return "", 0, errno.NewErrNo(errno.PaymentOrderNotExist, "order does not exist")
	}

	// 2. 获取用户id,并检查用户是否存在
	// 获取用户id
	var uid int64
	uid, err = uc.svc.GetUserID(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("get user id failed:%w", err)
	}

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
		if payStatus == paymentStatus.PaymentStatusSuccess || payStatus == paymentStatus.PaymentStatusProcessing {
			return "", 0, fmt.Errorf("payment is processing or has already done:%w", err)
		}
	}
	*/
	// 4. HMAC生成支付令牌
	//log.Println("GetPaymentToken called with orderID:", orderID)
	logger.Info("GetPaymentToken called", zap.Int64("orderID", orderID))
	token, expTime, err = uc.svc.GeneratePaymentToken(ctx, orderID)
	if err != nil {
		//log.Printf("Error generating payment token: %v", err)
		logger.Error("Error generating payment token",
			zap.Int64("orderID", orderID),
			zap.Error(err),
		)
		return "", 0, fmt.Errorf("generate payment token failed:%w", err)
	}
	//log.Printf("Generated token: %s, expires at: %d", token, expTime)
	logger.Info("Generated payment token",
		zap.String("token", token),
		zap.Int64("expTime", expTime),
	)

	var redisStatus bool
	// 5. 存储令牌到 Redis
	// TODO
	uid := int64(123)
	logger.Info("Storing token in Redis",
		zap.Int64("userID", uid),
		zap.Int64("orderID", orderID),
	)

	//log.Printf("Storing token in Redis for userID: %d, orderID: %d", uid, orderID)
	redisStatus, err = uc.svc.StorePaymentToken(ctx, token, expTime, uid, orderID)
	if err != nil && redisStatus != paymentStatus.RedisStoreSuccess {
		//log.Printf("Error storing payment token in Redis: %v", err)
		logger.Error("Error storing payment token in Redis",
			zap.Int64("orderID", orderID),
			zap.Int64("userID", uid),
			zap.Error(err),
		)
		return "", 0, fmt.Errorf("store payment token failed:%w", err)
	}
	//log.Println("Payment token stored successfully")
	logger.Info("Payment token stored successfully",
		zap.Int64("orderID", orderID),
		zap.Int64("userID", uid),
	)
	return token, expTime, nil
}
