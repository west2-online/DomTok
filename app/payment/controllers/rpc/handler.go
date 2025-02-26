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

	"github.com/west2-online/DomTok/app/payment/controllers/rpc/pack"
	"github.com/west2-online/DomTok/app/payment/usecase"
	"github.com/west2-online/DomTok/kitex_gen/payment"
	"github.com/west2-online/DomTok/pkg/base"
)

type PaymentHandler struct {
	useCase usecase.PaymentUseCase
}

func NewPaymentHandler(useCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

// ProcessPayment 任何一个error必须是errno类型的
// code路径：pkg/errno/code_service.go
// TODO 待完善
func (handler *PaymentHandler) ProcessPayment(ctx context.Context, req *payment.PaymentRequest) (r *payment.PaymentResponse, err error) {
	r = new(payment.PaymentResponse)
	p, err := handler.useCase.CreatePayment(ctx, req.GetOrderID())
	r.Status = p.Status
	if err != nil {
		return r, err
	}
	return r, err
}

func (handler *PaymentHandler) RequestPaymentToken(ctx context.Context, req *payment.PaymentTokenRequest) (r *payment.PaymentTokenResponse, err error) {
	r = new(payment.PaymentTokenResponse)
	var token string
	var expTime int64
	// 传入ctx（包含uid）和orderID,获取令牌和令牌过期时间
	token, expTime, err = handler.useCase.GetPaymentToken(ctx, req.OrderID)
	if err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	r.TokenInfo = pack.BuildTokenInfo(token, expTime)
	return
}

func (handler *PaymentHandler) ProcessRefund(ctx context.Context, req *payment.RefundReviewRequest) (r *payment.RefundReviewResponse, err error) {
	return nil, err
}

func (handler *PaymentHandler) RequestRefund(ctx context.Context, req *payment.RefundRequest) (r *payment.RefundResponse, err error) {
	r = new(payment.RefundResponse)
	/*var token string
	var expTime int64
	*/
	var refundID int64
	var refundStatus int64
	// 传入ctx（包含uid）和orderID,获取退款令牌和退款令牌过期时间
	refundStatus, refundID, err = handler.useCase.CreateRefund(ctx, req.OrderID)
	if err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	r.RefundInfo = pack.BuildRefundTokenInfo(refundID, refundStatus)
	return
}
