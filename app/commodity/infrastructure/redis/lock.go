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

	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (c *commodityCache) IsHealthy(ctx context.Context) error {
	newCtx, cancel := context.WithTimeout(ctx, constants.RedisCheckoutTimeOut)
	defer cancel()
	err := c.client.Ping(newCtx).Err()
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.IsHealthy failed: %v", err)
	}
	return nil
}

func (c *commodityCache) Lock(ctx context.Context, keys []string, ttl time.Duration) error {
	maxWaitTime := constants.RedisMaxLockRetryTime
	timeout := time.After(maxWaitTime)
	for {
		select {
		case <-timeout:
			return errno.Errorf(errno.InternalRedisErrorCode, "Timeout while waiting for locks")
		default:
			// 尝试获取锁
			locked := true
			_, err := c.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
				for _, key := range keys {
					lockSuccess, err := c.client.SetNX(ctx, key, 0, ttl).Result()
					if err != nil {
						locked = false
						return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.Lock failed: %v", err)
					}
					if !lockSuccess {
						locked = false
						return errno.Errorf(errno.InternalRedisErrorCode, "Lock already held for key: %v", key)
					}
				}
				return nil
			})
			if err == nil && locked {
				// 所有锁成功获取
				return nil
			}
			// retry
			time.Sleep(constants.RedisRetryStopTime)
		}
	}
}

func (c *commodityCache) UnLock(ctx context.Context, keys []string) error {
	err := c.client.Del(ctx, keys...).Err()
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.Unlock failed: %v", err)
	}
	return nil
}
