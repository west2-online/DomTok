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

	"github.com/west2-online/DomTok/app/payment/domain/model"
	"github.com/west2-online/DomTok/app/payment/domain/repository"
	"github.com/west2-online/DomTok/app/payment/domain/service"
)

type PaymentUseCase interface {
	CreatePayment(ctx context.Context, orderID int64) (*model.PaymentOrder, error)
	// GetPaymentToken 这里一次返回三个参数很不优雅但我不知道要怎么搞>_<
	GetPaymentToken(ctx context.Context, payment *model.PaymentOrder) (string, int64, error)
	// ProcessRefund
	// RequestRefundToken
}

type paymentUseCase struct {
	db  repository.PaymentDB
	svc *service.PaymentService
}

func NewPaymentCase(db repository.PaymentDB, svc *service.PaymentService) PaymentUseCase {
	return &paymentUseCase{
		db:  db,
		svc: svc,
	}
}
