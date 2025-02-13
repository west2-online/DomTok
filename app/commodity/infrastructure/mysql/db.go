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

package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	kmodel "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

// commodityDB impl domain.CommodityDB defined domain
type commodityDB struct {
	client *gorm.DB
}

func NewCommodityDB(client *gorm.DB) repository.CommodityDB {
	return &commodityDB{client: client}
}

func (d *commodityDB) IsCategoryExist(ctx context.Context, name string) (int64, error) {
	var category model.Category
	err := d.client.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query category: %v", err)
	}
	return category.CreatorId, nil
}

func (d *commodityDB) CreateCategory(ctx context.Context, entity *model.Category) error {
	model := Category{
		Id:        entity.Id,
		Name:      entity.Name,
		CreatorId: entity.CreatorId,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
	if err := d.client.WithContext(ctx).Create(model).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create category: %v", err)
	}
	return nil
}

func (d *commodityDB) DeleteCategory(ctx context.Context, category *model.Category) error {
	if err := d.client.WithContext(ctx).Delete(Category{Id: category.Id}).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete category: %v", err)
	}
	return nil
}

func (d *commodityDB) UpdateCategory(ctx context.Context, category *model.Category) error {
	if err := d.client.WithContext(ctx).Model(&model.Category{}).Where("id = ?", category.Id).Updates(category).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update category: %v", err)
	}
	return nil
}

func (u *commodityDB) ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*kmodel.CategoryInfo, err error) {
	offset := (pageNum - 1) * pageSize
	if err := u.client.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&resp).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to list categories: %v", err)
	}
	return resp, nil
}
