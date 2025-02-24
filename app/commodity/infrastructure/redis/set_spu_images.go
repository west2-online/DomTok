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
	_ "github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

func (c *commodityCache) SetSpuImages(ctx context.Context, key string, images *model.SpuImages) {
	dataJSON, err := sonic.Marshal(images)
	if err != nil {
		logger.Errorf("commodityCache.SetSpuImages marshal data failed: %v", err)
	}

	err = c.client.Set(ctx, key, dataJSON, constants.RedisSpuImageExpireTime).Err()
	if err != nil {
		logger.Errorf("commodity.SetSpuImages set data failed: %v", err)
	}
}
