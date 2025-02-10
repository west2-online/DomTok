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

	"github.com/west2-online/DomTok/app/user/domain"
	"github.com/west2-online/DomTok/pkg/errno"
)

// userDB impl domain.UserDB defined domain
type userDB struct {
	client *gorm.DB
}

func NewUserDB(client *gorm.DB) domain.UserDB {
	return &userDB{client: client}
}

func (db *userDB) CreateUser(ctx context.Context, u *domain.User) error {
	// 将 entity 转换成 mysql 这边的 model
	// TODO 可以考虑整一个函数统一转化, 放在这里占了太多行, 而且这不是这个方法该做的. 这个方法应该做的是创建用户
	model := User{
		UserName: u.UserName,
		Password: u.Password,
		Email:    u.Email,
	}

	if err := db.client.WithContext(ctx).Create(model).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create user: %v", err)
	}
	return nil
}

func (db *userDB) IsUserExist(ctx context.Context, username string) (bool, error) {
	var user User
	err := db.client.WithContext(ctx).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		// 这里虽然是数据库返回的 err 不为 nil,
		// 但这显然是业务上的错误, 而不是我们服务本身的
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	return true, nil
}
