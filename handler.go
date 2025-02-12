package main

import (
	"context"
	payment "github.com/west2-online/DomTok/kitex_gen/payment"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// ProcessPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) ProcessPayment(ctx context.Context, request *payment.PaymentRequest) (resp *payment.PaymentResponse, err error) {
	// TODO: Your code here...
	return
}

// RequestPaymentToken implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) RequestPaymentToken(ctx context.Context, request *payment.PaymentTokenRequest) (resp *payment.PaymentTokenResponse, err error) {
	// TODO: Your code here...
	return
}

// ProcessRefund implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) ProcessRefund(ctx context.Context, request *payment.RefundRequest) (resp *payment.RefundResponse, err error) {
	// TODO: Your code here...
	return
}

// RequestRefundToken implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) RequestRefundToken(ctx context.Context, request *payment.RefundTokenRequest) (resp *payment.RefundTokenResponse, err error) {
	// TODO: Your code here...
	return
}
