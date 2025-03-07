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
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
)

type RedisService struct {
	client redis.Client
}

func NewRedisService() *RedisService {
	cli, err := client.NewRedisClient(constants.RedisDBGateWay)
	if err != nil {
		panic(err)
	}

	return &RedisService{*cli}
}

func (svc *RedisService) IsUserBanned(ctx context.Context, userId int64) bool {
	return svc.client.Exists(ctx, svc.GetUserBanedKey(userId)).Val() == 1
}

func (svc *RedisService) IsUserLogout(ctx context.Context, userId int64) bool {
	return svc.client.Exists(ctx, svc.GetUserLogoutKey(userId)).Val() != 1
}

func (svc *RedisService) GetAllBanedUser(ctx context.Context) []int64 {
	keys, err := svc.client.Keys(ctx, constants.RedisUserBanedKey+"*").Result()
	if err != nil {
		logger.Fatalf("get all baned user failed: %v", err)
	}
	var res []int64
	for _, key := range keys {
		var userId int64
		_, _ = fmt.Sscanf(key, constants.RedisUserBanedKey+"%d", &userId)
		res = append(res, userId)
	}
	return res
}

func (svc *RedisService) GetUserBanedKey(userId int64) string {
	return fmt.Sprintf(constants.RedisUserBanedKey+"%d", userId)
}

func (svc *RedisService) GetUserLogoutKey(userId int64) string {
	return fmt.Sprintf(constants.RedisUserLogoutKey+"%d", userId)
}
