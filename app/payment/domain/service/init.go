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
	"github.com/west2-online/DomTok/app/payment/domain/repository"
	"github.com/west2-online/DomTok/pkg/utils"
)

type PaymentService struct {
	db    repository.PaymentDB
	sf    *utils.Snowflake
	redis repository.PaymentRedis
	rpc   repository.PaymentRPC
}

func NewPaymentService(db repository.PaymentDB, sf *utils.Snowflake, redis repository.PaymentRedis, rpc repository.PaymentRPC) *PaymentService {
	if db == nil {
		panic("paymentService`s db should not be nil")
	}
	if sf == nil {
		panic("paymentService`s sf should not be nil")
	}
	if redis == nil {
		panic("paymentService`s redis should not be nil")
	}
	if rpc == nil {
		panic("paymentService`s rpc should not be nil")
	}
	svc := &PaymentService{
		db:    db,
		sf:    sf,
		redis: redis,
		rpc:   rpc,
	}
	return svc
}
