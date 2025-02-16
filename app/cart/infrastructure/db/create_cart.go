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

package db

import (
	"context"

	"github.com/west2-online/DomTok/pkg/errno"
)

// CreateCart 创建购物车
func (c *DBAdapter) CreateCart(ctx context.Context, uid int64, cart string) error {
	model := Cart{
		UserId:  uid,
		SkuJson: cart,
	}
	if err := c.client.WithContext(ctx).Create(&model).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "db.CreateCart error: %v", err)
	}
	return nil
}
