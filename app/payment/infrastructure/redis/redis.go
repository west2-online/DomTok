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

package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/DomTok/app/payment/domain/repository"
)

type paymentRedis struct {
	client *redis.Client
}

func NewPaymentRedis(client *redis.Client) repository.PaymentRedis {
	return &paymentRedis{client: client}
}

func (p *paymentRedis) SetPaymentToken(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return p.client.Set(ctx, key, value, expiration).Err()
}

func (p *paymentRedis) GetPaymentToken(ctx context.Context, key string) (string, error) {
	return p.client.Get(ctx, key).Result()
}
