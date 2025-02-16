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
	"github.com/west2-online/DomTok/pkg/errno"
)

// CreatePayment 这里定义一些具体的方法和函数，比如校验密码，加密密码，创建用户之类的
func (uc *paymentUseCase) CreatePayment(ctx context.Context, orderID int64) (*model.PaymentOrder, error) {
	return nil, nil
}

func (uc *paymentUseCase) GetPaymentToken(ctx context.Context, orderID int64) (token string, expTime int64, err error) {
	// 1. 检查订单是否存在
	// TODO 这个要向order模块要一个RPC接口然后再来填充
	orderInfo, err := uc.svc.CheckOrderExist(ctx, orderID)
	if err != nil {
		return "", 0, fmt.Errorf("check order existed failed:%w", err)
	}
	if orderInfo == paymentStatus.OrderNotExist {
		return "", 0, errno.NewErrNo(errno.PaymentOrderNotExist, "order does not exist")
	}

	// 2. 获取用户id,并检查用户是否存在
	// 获取用户id
	// TODO 这个函数要等user那边写完了才能填
	uid, err := uc.svc.GetUserID(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("get user id failed:%w", err)
	}
	// 检查用户是否存在
	userInfo, err := uc.svc.CheckUserExist(ctx, uid)
	if err != nil {
		return "", 0, fmt.Errorf("check user existed failed:%w", err)
	}
	if userInfo == paymentStatus.UserNotExist {
		return "", 0, errno.NewErrNo(errno.UserNotExist, "user does not exist")
	}

	// 3. 检查订单支付信息
	paymentInfo, err := uc.db.CheckPaymentExist(ctx, orderID)
	if err != nil {
		return "", 0, fmt.Errorf("check payment existed failed:%w", err)
	}
	// 如果订单不存在
	if paymentInfo == paymentStatus.PaymentNotExist {
		// 创建支付订单
		// TODO 待完善
		_, err := uc.svc.CreatePaymentInfo(ctx, orderID)
		// TODO 为什么这里的err是绿色的？？？
		if err != nil {
			return "", 0, fmt.Errorf("create payment info failed:%w", err)
		}
		// 如果订单存在
	} else if paymentInfo == paymentStatus.PaymentExist {
		// 获取订单的支付状态
		payStatus, err := uc.db.GetPaymentInfo(ctx, orderID)
		// 如果订单正在支付或者已经支付完成，则拒绝进行接下来的生成令牌的活动
		if payStatus == paymentStatus.PaymentStatusSuccess || payStatus == paymentStatus.PaymentStatusProcessing {
			return "", 0, fmt.Errorf("payment is processing or has already done:%w", err)
		}
		// TODO 这个else有必要保留吗？
	} // else {
	// return "", 0, fmt.Errorf("check payment existed failed:%w", err)
	// }

	// 4. HMAC生成支付令牌
	token, expTime, err = uc.svc.GeneratePaymentToken(ctx, orderID)
	if err != nil {
		return "", 0, fmt.Errorf("generate payment token failed:%w", err)
	}
	var redisStatus int
	// 5. 存储令牌到 Redis
	redisStatus, err = uc.svc.StorePaymentToken(ctx, token, expTime)
	if err != nil && redisStatus != paymentStatus.RedisStoreSuccess {
		return "", 0, fmt.Errorf("store payment token failed:%w", err)
	}
	return token, expTime, nil
}
