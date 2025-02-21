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

	entities "github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (svc *CommodityService) DeleteCategory(ctx context.Context, category *entities.Category) (err error) {
	// 判断是否存在
	exist, err := svc.db.IsCategoryExist(ctx, category.Id)
	if err != nil {
		return fmt.Errorf("check category exist failed: %w", err)
	}
	if !exist {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "category does not exist")
	}
	err = svc.db.DeleteCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("delete category failed: %w", err)
	}
	return nil
}
func (svc *CommodityService) nextID() int64 {
	id, _ := svc.sf.NextVal()
	return id
}

func (svc *CommodityService) CreateCategory(ctx context.Context, category *entities.Category) error {
	category.Id = svc.nextID()
	if err := svc.db.CreateCategory(ctx, category); err != nil {
		return fmt.Errorf("create category failed: %w", err)
	}
	return nil
}
