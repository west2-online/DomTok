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

package usecase

import (
	"context"

	"github.com/west2-online/DomTok/app/user/domain/model"
	"github.com/west2-online/DomTok/app/user/domain/repository"
	"github.com/west2-online/DomTok/app/user/domain/service"
)

// UserUseCase 接口应该不应该定义在 domain 中，这属于 use case 层
type UserUseCase interface {
	RegisterUser(ctx context.Context, user *model.User) (uid int64, err error)
	Login(ctx context.Context, user *model.User) (*model.User, error)
	GetAddress(ctx context.Context, addressID int64) (*model.Address, error)
	AddAddress(ctx context.Context, address *model.Address) (addressID int64, err error)
	BanUser(ctx context.Context, uid int64) error
	LiftUser(ctx context.Context, uid int64) error
	LogoutUser(ctx context.Context) error
}

// useCase 实现了 domain.UserUseCase
// 只会以接口的形式被调用, 所以首字母小写改为私有类型
type useCase struct {
	db    repository.UserDB
	svc   *service.UserService
	cache repository.UserCache
}

func NewUserCase(db repository.UserDB, svc *service.UserService, re repository.UserCache) *useCase {
	return &useCase{
		db:    db,
		svc:   svc,
		cache: re,
	}
}
