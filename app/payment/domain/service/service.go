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

package service

import (
	"context"
	"github.com/west2-online/DomTok/app/payment/domain/model"
)

// sf可以生成id,详见user/domain/service/service.go
type PaymentService interface {
	GetOrderByID(ctx context.Context, p *model.PaymentOrder) (interface{}, error)
	GetUserByID(ctx context.Context, p *model.PaymentOrder) (interface{}, error)
	GetPaymentInfo(ctx context.Context, p *model.PaymentOrder) (int, error)
	GeneratePaymentToken(ctx context.Context, p *model.PaymentOrder) (string, int64, error)
}

func (svc *PaymentService) CreatePaymentInfo(ctx context.Context, p *model.PaymentOrder) (int64, error) {
	return 0, nil
}
