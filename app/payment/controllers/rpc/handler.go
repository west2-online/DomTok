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

package rpc

import (
	"context"
	"fmt"
	"github.com/west2-online/DomTok/app/payment/domain/model"
	"github.com/west2-online/DomTok/app/payment/usecase"
	"github.com/west2-online/DomTok/kitex_gen/payment"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
)

type PaymentHandler struct {
	useCase usecase.PaymentUseCase
}

func NewPaymentHandler(useCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

// 任何一个error必须是errno类型的
// code路径：pkg/errno/code_service.go
func (handler *PaymentHandler) ProcessPayment(ctx context.Context, req *payment.PaymentRequest) (r *payment.PaymentResponse, err error) {
	r = new(payment.PaymentResponse)
	p, err := handler.useCase.ProcessPayment(ctx, req.GetOrderID())
	if err != nil {
		return r, err
	}


	r.Base = p.

	r.PaymentID = p.PaymentID
	r.Status = p.Status
	return r, err
}

func (handler *PaymentHandler) RequestPaymentToken(ctx context.Context, req *payment.PaymentTokenRequest) (r *payment.PaymentTokenResponse, err error) {
	r = new(payment.PaymentTokenResponse)
	p:=&model.PaymentOrder{
		OrderID: req.OrderID,
		UserID: req.UserID,
	}

	// 1. 检查订单是否存在
	_,err=handler.useCase.GetOrderByID(ctx,p)
	// 这里直接return就可以吗？
	if err != nil {
		return
	}

	// 2. 检查用户是否存在
	_, err = handler.useCase.GetUserByID(ctx,p)
	if err != nil {
		return
	}

	// 3. 检查订单支付信息
	// 这里用int还是int8？
	var paymentInfo int
	paymentInfo, err = handler.useCase.GetPaymentInfo(ctx, p)
	if err != nil {
		return
	}
	if paymentInfo == paymentStatus.PaymentStatusSuccess||paymentInfo == paymentStatus.PaymentStatusProcessing {
		return
	} else {
		err:=handler.useCase.CreatePaymentInfo(ctx,p)
		if err != nil {
			return
		}
	}

	// 4. 生成支付令牌
	var token string
	var expTime int64
	// 感觉这里一次返回三个值非常非常非常不优雅，但是不知道要怎么写得更优雅
	token,expTime,err = handler.useCase.GeneratePaymentToken(ctx, p)
	if err != nil {
		return
	}
	r.PaymentToken = token
	r.ExpirationTime = expTime

	// 5. 存储令牌到 Redis
	err = handler.useCase.StorePaymentToken(ctx, p)
	if err != nil {
		return
	}

	return
}

func (handler *PaymentHandler) ProcessRefund(ctx context.Context, req *payment.RefundRequest) (r *payment.RefundResponse, err error) {
	return nil, err
}

func (handler *PaymentHandler) RequestRefundToken(ctx context.Context, req *payment.RefundTokenRequest) (r *payment.RefundTokenResponse, err error) {
	return nil, err
}
