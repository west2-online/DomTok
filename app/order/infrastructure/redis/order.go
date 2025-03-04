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

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

type orderCache struct {
	client                 *redis.Client
	expire                 time.Duration
	updatePaymentStatusLua string
}

func NewOrderCache(client *redis.Client) repository.Cache {
	c := &orderCache{client: client}
	c.loadUpdateLUAScript()
	c.expire = constants.OrderPaymentStatusExpireTime
	return c
}

func (cache *orderCache) SetPaymentStatus(ctx context.Context, s *model.CachePaymentStatus) error {
	exK, sK := getExpireKey(s.OrderID), getStatusKey(s.OrderID)
	if err := cache.client.Set(ctx, exK, s.OrderExpire, cache.expire).Err(); err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed set key: %s to %v, err: %v", exK, sK, err))
	}
	if err := cache.client.Set(ctx, sK, s.PaymentStatus, cache.expire).Err(); err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed set key: %s to %v, err: %v", exK, sK, err))
	}
	return nil
}

func (cache *orderCache) GetPaymentStatus(ctx context.Context, orderID int64) (*model.CachePaymentStatus, bool, error) {
	exK, sK := getExpireKey(orderID), getStatusKey(orderID)
	var ex, s int64
	var err error

	if ex, err = cache.client.Get(ctx, exK).Int64(); err != nil {
		if errors.Is(err, redis.Nil) {
			return &model.CachePaymentStatus{}, false, nil
		}
		return nil, false, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed get key: %s,err: %v", exK, err))
	}
	if s, err = cache.client.Get(ctx, sK).Int64(); err != nil {
		if errors.Is(err, redis.Nil) {
			return &model.CachePaymentStatus{}, false, nil
		}
		return nil, false, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed get key: %s,err: %v", exK, err))
	}

	return &model.CachePaymentStatus{OrderExpire: ex, PaymentStatus: int8(s)}, true, nil
}

// UpdatePaymentStatus 使用 lua 脚本保证了过程的原子性
func (cache *orderCache) UpdatePaymentStatus(ctx context.Context, s *model.CachePaymentStatus) (bool, error) {
	exK, sK := getExpireKey(s.OrderID), getStatusKey(s.OrderID)
	result, err := cache.client.EvalSha(ctx, cache.updatePaymentStatusLua, []string{exK, sK}, s.OrderExpire, s.PaymentStatus, cache.expire).Result()
	if err != nil {
		return false, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed to execute lua script: %v", err))
	}
	rel, ok := result.(int64)
	if !ok {
		return false, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed to execute lua script: %v", err))
	}

	return rel == constants.OrderCacheLuaKeyExistFlag, nil
}

func (cache *orderCache) DeletePaymentStatus(ctx context.Context, orderID int64) error {
	exK, sK := getExpireKey(orderID), getStatusKey(orderID)
	if err := cache.client.Del(ctx, exK, sK).Err(); err != nil {
		return errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("failed to delete key: %s, %s, err: %v", exK, sK, err))
	}
	return nil
}

func (cache *orderCache) loadUpdateLUAScript() {
	sha1, err := cache.client.ScriptLoad(context.Background(), constants.OrderUpdatePaymentStatusLuaScript).Result()
	if err != nil {
		logger.Fatalf("failed to load lua script: %v", err)
	}
	cache.updatePaymentStatusLua = sha1
}

func getExpireKey(orderID int64) string {
	return fmt.Sprintf(constants.OrderCacheOrderExpireFormat, orderID)
}

func getStatusKey(orderID int64) string {
	return fmt.Sprintf(constants.OrderCachePaymentStatusFormat, orderID)
}
