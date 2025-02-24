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

package repository

import (
	"context"
	"time"

	"github.com/west2-online/DomTok/app/payment/domain/model"
)

type PaymentDB interface {
	// CheckPaymentExist CheckUserExist(ctx context.Context, uid int64) (userInfo interface{}, err error)
	CheckPaymentExist(ctx context.Context, orderID int64) (paymentInfo bool, err error)
	GetPaymentInfo(ctx context.Context, orderID int64) (payStatus interface{}, err error)
	CreatePayment(ctx context.Context, order *model.PaymentOrder) error
	CreateRefund(ctx context.Context, order *model.PaymentRefund) error
}
type PaymentRedis interface {
	SetPaymentToken(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	IncrRedisKey(ctx context.Context, key string, expiration int) (int64, error)
	CheckRedisDayKey(ctx context.Context, key string) (bool, error)
	SetRedisDayKey(ctx context.Context, key string, value string, expiration int) error
	SetRefundToken(ctx context.Context, key string, token string, duration time.Duration) error
	// GetPaymentToken(ctx context.Context, key string) (string, error)
}

type RPC interface {
	PaymentIsOrderExist(ctx context.Context, orderID int64) (bool, error)
}
