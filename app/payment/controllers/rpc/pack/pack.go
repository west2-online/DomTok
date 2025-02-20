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

<<<<<<<< HEAD:app/payment/domain/repository/interface.go
package repository

import (
	"context"
)

type PaymentDB interface {
	GetOrderByToken(ctx context.Context, paramToken string) (int64, error)
	GetUserByToken(ctx context.Context, paramToken string) (int64, error)
	GetPaymentInfo(ctx context.Context, paramToken string) (int, error)
========
package pack

import (
	"github.com/west2-online/DomTok/kitex_gen/model"
)

// BuildPaymentOrder BuildUser 将 entities 定义的 Payment 实体转换成 idl 定义的 RPC 交流实体，类似 dto

func BuildTokenInfo(token string, expTime int64) *model.PaymentTokenInfo {
	return &model.PaymentTokenInfo{
		PaymentToken:               token,
		PaymentTokenExpirationTime: expTime,
	}
>>>>>>>> upstream/main:app/payment/controllers/rpc/pack/pack.go
}
