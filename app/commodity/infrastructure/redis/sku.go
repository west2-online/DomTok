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

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

func (c commodityCache) SetSku(ctx context.Context, key string, sku *model.Sku) {
	dataJSON, err := sonic.Marshal(sku)
	if err != nil {
		logger.Errorf("commodityCache.SetSku marshal data failed: %v", err)
	}

	err = c.client.Set(ctx, key, dataJSON, constants.RedisSkuExpireTime).Err()
	if err != nil {
		logger.Errorf("commodity.SetSku set data failed: %v", err)
	}
}

func (c commodityCache) GetSku(ctx context.Context, key string) (*model.Sku, error) {
	dataJSON, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, errno.Errorf(errno.InternalRedisErrorCode, "commodityCache.GetSku faile: %v", err)
	}
	ret := new(model.Sku)
	err = sonic.Unmarshal([]byte(dataJSON), &ret)
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "commodityCache.GetSku Unmarshal failed: %v", err)
	}
	return ret, nil
}

func (c commodityCache) DeleteSku(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "commodityCache.DeleteSku failed: %v", err)
	}
	return nil
}
