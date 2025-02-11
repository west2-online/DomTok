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

	"github.com/west2-online/DomTok/app/commodity/entities"
	"github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

type DBAdapter struct {
	client *gorm.DB
}

func NewDBAdapter(client *gorm.DB) *DBAdapter {
	return &DBAdapter{client: client}
}

func (d *DBAdapter) IsCategoryExist(ctx context.Context, Name string) (bool, error) {
	var category entities.Category
	err := d.client.WithContext(ctx).Where("name = ?", Name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query category: %v", err)
	}
	return true, nil
}

func (d *DBAdapter) CreateCategory(ctx context.Context, entity *entities.Category) error {
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

func (d *DBAdapter) DeleteCategory(ctx context.Context, category *entities.Category) error {
	if err := d.client.WithContext(ctx).Delete(Category{Id: category.Id}).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete category: %v", err)
	}
	return nil
}

func (d *DBAdapter) UpdateCategory(ctx context.Context, category *entities.Category) error {
	if err := d.client.WithContext(ctx).Model(&entities.Category{}).Where("id = ?", category.Id).Updates(category).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update category: %v", err)
	}
	return nil
}

func (u *DBAdapter) ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*model.CategoryInfo, err error) {
	offset := (pageNum - 1) * pageSize
	if err := u.client.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&resp).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to list categories: %v", err)
	}
	return resp, nil
}
