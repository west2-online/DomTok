package rpc

import (
	"context"
	"github.com/west2-online/DomTok/app/payment/usecase"
	"github.com/west2-online/DomTok/kitex_gen/payment"
)

type PaymentHandler struct {
	useCase usecase.PaymentUseCase
}

// usecase只管逻辑，不涉及具体实现
func (handler *PaymentHandler) ProcessPayment(ctx context.Context, req *payment.PaymentRequest) (r *payment.PaymentResponse, err error) {
	r = new(payment.PaymentResponse)
	p, err := handler.useCase.ProcessPayment(ctx, req.GetOrderID())
	if err != nil {
		return r, err
	}
	r.Base = p.Base
	r.PaymentID = p.PaymentID
	r.Status = p.Status
	return r, err
}

func (handler *PaymentHandler) RequestPaymentToken(ctx context.Context, req *payment.PaymentTokenRequest) (r *payment.PaymentTokenResponse, err error) {
	return nil, err
}

func (handler *PaymentHandler) ProcessRefund(ctx context.Context, req *payment.RefundRequest) (r *payment.RefundResponse, err error) {
	return nil, err
}

func (handler *PaymentHandler) RequestRefundToken(ctx context.Context, req *payment.RefundTokenRequest) (r *payment.RefundTokenResponse, err error) {
	return nil, err
}
