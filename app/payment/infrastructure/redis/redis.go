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
	"github.com/west2-online/DomTok/pkg/errno"
)

type paymentRedis struct {
	client *redis.Client
}

func NewPaymentRedis(client *redis.Client) repository.PaymentRedis {
	cli := paymentRedis{client: client}
	err := cli.loadScript()
	if err != nil {
		panic(err)
	}
	return &cli
}

func (p *paymentRedis) SetPaymentToken(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return p.client.Set(ctx, key, value, expiration).Err()
}

// IncrRedisKey 自增 Redis 键，并在第一次设置过期时间
func (p *paymentRedis) IncrRedisKey(ctx context.Context, key string, expiration int) (int64, error) {
	// 自增键值
	count, err := p.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, errno.Errorf(errno.InternalRedisErrorCode, "failed to increment key %s: %v", key, err)
		// return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return count, nil
}

func (p *paymentRedis) CheckRedisDayKey(ctx context.Context, key string) (bool, error) {
	// exists返回1表示key存在，返回0表示key不存在
	exists, err := p.client.Exists(ctx, key).Result()
	if err != nil {
		return false, errno.Errorf(errno.InternalRedisErrorCode, "failed to check 24h existence of key %s: %v", key, err)
	}
	return exists == 1, nil
}

// SetRedisDayKey 设置 Redis 键值，并指定过期时间
func (p *paymentRedis) SetRedisDayKey(ctx context.Context, key string, value string, expiration int) error {
	err := p.client.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "failed to set key %s in Redis: %v", key, err)
	}
	return nil
}

func (p *paymentRedis) SetRefundToken(ctx context.Context, key string, value string, expiration time.Duration) error {
	return p.client.Set(ctx, key, value, expiration).Err()
}

// func (p *paymentRedis) GetPaymentToken(ctx context.Context, key string) (string, error) {
// 	   return p.client.Get(ctx, key).Result()
// }

// CheckAndDelPaymentToken 检查并删除退款令牌
func (p *paymentRedis) CheckAndDelPaymentToken(ctx context.Context, key string, value string) (bool, error) {
	// 执行脚本
	result, err := p.execScript(ctx, CheckAndDelScript, []string{key}, value)
	if err != nil {
		return false, fmt.Errorf("failed to check and delete refund token: %w", err)
	}
	exist, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("failed to convert result to int64")
	}
	return exist == 1, nil
}

func (p *paymentRedis) GetTTLAndDelPaymentToken(ctx context.Context, key string, value string) (bool, time.Duration, error) {
	// 执行脚本
	result, err := p.execScript(ctx, GetTTLAndDelScript, []string{key}, value)
	if err != nil {
		return false, -1, fmt.Errorf("failed to get ttl and delete refund token: %w", err)
	}
	res, ok := result.([]interface{})
	if !ok || len(res) != 2 {
		return false, -1, fmt.Errorf("failed to convert result to [2]interface{}")
	}
	redisTTL, ok := res[0].(int64)
	if !ok {
		return false, -1, fmt.Errorf("failed to convert ttl to int64")
	}
	redisExist, ok := res[1].(int64)
	if !ok {
		return false, -1, fmt.Errorf("failed to convert exist to int64")
	}
	return redisExist == 1, time.Duration(redisTTL) * time.Second, nil
}
