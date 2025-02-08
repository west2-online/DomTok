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

	"github.com/west2-online/DomTok/app/user/entities"
	"github.com/west2-online/DomTok/pkg/errno"
)

// DBAdapter impl PersistencePort defined in use case package
type DBAdapter struct {
	client *gorm.DB
}

func NewDBAdapter(client *gorm.DB) *DBAdapter {
	return &DBAdapter{client: client}
}

func (d *DBAdapter) CreateUser(ctx context.Context, entity *entities.User) error {
	// 将 entity 转换成 mysql 这边的 model
	model := User{
		UserName: entity.UserName,
		Password: entity.Password,
		Email:    entity.Email,
	}
	if err := d.client.WithContext(ctx).Create(model).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create user: %v", err)
	}
	return nil
}

func (d *DBAdapter) IsUserExist(ctx context.Context, username string) (bool, error) {
	var user entities.User
	err := d.client.WithContext(ctx).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	return true, nil
}
