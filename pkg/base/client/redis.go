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

package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

func InitRedis(db int) (*redis.Client, error) {
	if config.Redis == nil {
		return nil, errors.New("redis config is nil")
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Password,
		DB:           db,
		PoolSize:     constants.RedisPoolSize,           // 连接池大小
		MinIdleConns: constants.RedisMinIdleConnections, // 最小空闲连接数
		DialTimeout:  constants.RedisDialTimeout,        // 连接超时
	})

	// 添加日志 Hook
	l := logger.GetRedisLogger()
	redis.SetLogger(l)
	client.AddHook(l)

	// 使用超时的 Ping
	ctx, cancel := context.WithTimeout(context.Background(), constants.PingTime*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("client.NewRedisClient: ping redis failed: %v", err))
	}

	return client, nil
}

func InitRedSync(client *redis.Client) *redsync.Redsync {
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)
	return rs
}
