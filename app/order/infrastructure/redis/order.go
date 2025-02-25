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
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

type orderCache struct {
	client *redis.Client
}

func NewOrderCache(client *redis.Client) repository.Cache {
	return &orderCache{client: client}
}

func (cache *orderCache) SetPaymentResultRecord(ctx context.Context, orderID int64, data []byte, expire time.Duration) error {
	key := getKey(orderID)
	if err := cache.client.Set(ctx, key, data, expire); err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed set kv to redis, err: %v", err))
	}
	return nil
}
func (cache *orderCache) GetPaymentResultRecord(ctx context.Context, orderID int64) ([]byte, bool, error) {
	key := getKey(orderID)
	data, err := cache.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed get kv from redis, err: %v", err))
	}
	return data, true, nil
}

func (cache *orderCache) DelPaymentResultRecord(ctx context.Context, orderID int64) error {
	key := getKey(orderID)
	if err := cache.client.Del(ctx, key).Err(); err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed del kv from redis, err: %v", err))
	}
	return nil
}

func getKey(orderID int64) string {
	return fmt.Sprintf(constants.OrderID2PaymentStatusFormat, orderID)
}
