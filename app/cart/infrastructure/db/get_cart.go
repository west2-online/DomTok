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
	"errors"
	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

// GetCartByUserId 查询购物车
func (c *DBAdapter) GetCartByUserId(ctx context.Context, uid int64) (bool, *model.Cart, error) {
	dbModel := new(Cart)
	if err := c.client.WithContext(ctx).Where("user_id=?", uid).First(dbModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, errno.Errorf(errno.InternalDatabaseErrorCode, "db.GetCartByUserId error: %v", err)
	}
	return true, &model.Cart{
		UserId:  uid,
		SkuJson: dbModel.SkuJson,
	}, nil
}
