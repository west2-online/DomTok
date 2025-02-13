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
