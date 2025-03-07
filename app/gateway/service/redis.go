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
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
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

func (svc *RedisService) IsUserBanned(ctx context.Context, userId int64) (bool, error) {
	err := svc.client.Exists(ctx, svc.GetUserBanedKey(userId)).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalRedisErrorCode, "RedisService.IsUserBanned failed: %v", err)
	}
	return true, nil
}

func (svc *RedisService) IsUserLogout(ctx context.Context, userId int64) (bool, error) {
	err := svc.client.Exists(ctx, svc.GetUserLogoutKey(userId)).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalRedisErrorCode, "RedisService.IsUserLogout failed: %v", err)
	}
	return true, nil
}

func (svc *RedisService) GetUserBanedKey(userId int64) string {
	return fmt.Sprintf(constants.RedisUserBanedKey+"%d", userId)
}

func (svc *RedisService) GetUserLogoutKey(userId int64) string {
	return fmt.Sprintf(constants.RedisUserLogoutKey+"%d", userId)
}
