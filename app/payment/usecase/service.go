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

// GetPaymentToken 这里要怎么让他一次只返回两个参数呢
func (uc *paymentUseCase) GetPaymentToken(ctx context.Context, p *model.PaymentOrder) (token string, expTime int64, err error) {
	// 1. 检查订单是否存在
	pid, err := uc.db.GetOrderByID(ctx, p)
	// 这里直接return就可以吗？
	if err != nil {
		return paymentStatus.PaymentOrderNotExistToken, paymentStatus.PaymentOrderNotExistExpirationTime, fmt.Errorf("check payment order existed failed:%w", err)
	}
	if pid == paymentStatus.PaymentOrderNotExist {
		return paymentStatus.PaymentOrderNotExistToken, paymentStatus.PaymentOrderNotExistExpirationTime, errno.NewErrNo(errno.PaymentOrderNotExist, "payment order does not exist")
	}

	// 2. 检查用户是否存在
	uid, err := uc.db.GetUserByID(ctx, p)
	if err != nil {
		return paymentStatus.UserNotExistToken, paymentStatus.UserNotExistExpirationTime, fmt.Errorf("check user existed failed:%w", err)
	}
	if uid == paymentStatus.UserNotExist {
		return paymentStatus.UserNotExistToken, paymentStatus.UserNotExistExpirationTime, errno.NewErrNo(errno.UserNotExist, "user does not exist")
	}

	// 3. 检查订单支付信息
	// 这里用int还是int8？
	var paymentInfo int
	paymentInfo, err = uc.db.GetPaymentInfo(ctx, p)
	if err != nil {
		return paymentStatus.PaymentOrderNotExistToken, paymentStatus.PaymentOrderNotExistExpirationTime, fmt.Errorf("check payment information failed:%w", err)
	}
	if paymentInfo == paymentStatus.PaymentStatusSuccess || paymentInfo == paymentStatus.PaymentStatusProcessing {
		return paymentStatus.HavePaidToken, paymentStatus.HavePaidExpirationTime, fmt.Errorf("payment is processing or has already done:%w", err)
	} else {
		// 创建支付订单
		_, err := uc.svc.CreatePaymentInfo(ctx, p)
		if err != nil {
			return paymentStatus.ErrorToken, paymentStatus.ErrorExpirationTime, fmt.Errorf("create payment info failed:%w", err)
		}
	}

	// 4. 生成支付令牌

	// 感觉这里一次返回三个值非常非常非常不优雅，但是不知道要怎么写得更优雅
	token, expTime, err = uc.svc.GeneratePaymentToken(ctx, p)
	if err != nil {
		return paymentStatus.ErrorToken, paymentStatus.ErrorExpirationTime, fmt.Errorf("generate payment token failed:%w", err)
	}

	// 5. 存储令牌到 Redis
	err = uc.svc.StorePaymentToken(ctx, p)
	if err != nil {
		return paymentStatus.ErrorToken, paymentStatus.ErrorExpirationTime, fmt.Errorf("store payment token failed:%w", err)
	}
	return token, expTime, nil
}

// 这里没有直接调用 db.CreateUser 是因为 svc.CreateUser 包含了一点业务逻辑, 这些细节不需要被 useCase 知道
// if err = uc.svc.CreateUser(ctx, u); err != nil {
// return
// }

// return u.Uid, nil
// }
