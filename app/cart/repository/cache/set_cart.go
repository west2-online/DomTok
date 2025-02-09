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

package cache

import (
	"context"
	"strconv"

	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

// SetCartCache 将购物车存入redis
func (c *CacheAdapter) SetCartCache(ctx context.Context, uid int64, cart string) error {
	if err := c.client.Set(ctx, strconv.FormatInt(uid, 10), cart, constants.RedisCartExpireTime).Err(); err != nil {
		logger.Errorf("cache.SetCartCache error:%v", err)
		return errno.Errorf(errno.InternalRedisErrorCode, "cache.SetCartCache error:%v", err.Error())
	}
	return nil
}
