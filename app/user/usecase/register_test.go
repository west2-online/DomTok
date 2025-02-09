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

package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/west2-online/DomTok/app/user/entities"
	"github.com/west2-online/DomTok/app/user/usecase"
	"github.com/west2-online/DomTok/app/user/usecase/mocks"
)

// 全部由 ai 生成, 写法待改进

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockDB := new(mocks.PersistencePort)
	ctx := context.Background()
	username := "existing_user"

	// 模拟IsUserExist返回用户存在
	mockDB.On("IsUserExist", ctx, username).Return(true, nil)

	u := &usecase.UseCase{DB: mockDB}
	entity := &entities.User{UserName: username}
	uid, err := u.RegisterUser(ctx, entity)

	// 断言错误为用户已存在
	assert.Equal(t, int64(0), uid)
	assert.Error(t, err)

	mockDB.AssertExpectations(t)
}

func TestRegisterUser_DBCheckError(t *testing.T) {
	mockDB := new(mocks.PersistencePort)
	ctx := context.Background()
	username := "testuser"
	expectedErr := fmt.Errorf("database error")

	// 模拟IsUserExist返回错误
	mockDB.On("IsUserExist", ctx, username).Return(false, expectedErr)

	u := &usecase.UseCase{DB: mockDB}
	entity := &entities.User{UserName: username}
	uid, err := u.RegisterUser(ctx, entity)

	// 断言错误被正确包装
	assert.Equal(t, int64(0), uid)
	assert.ErrorContains(t, err, "check user exist failed")
	assert.ErrorIs(t, err, expectedErr)
	mockDB.AssertExpectations(t)
}

func TestRegisterUser_InvalidEmail(t *testing.T) {
	mockDB := new(mocks.PersistencePort)
	ctx := context.Background()
	entity := &entities.User{
		UserName: "testuser",
		Email:    "invalid-email", // 无效邮箱
	}

	mockDB.On("IsUserExist", ctx, entity.UserName).Return(false, nil)

	u := &usecase.UseCase{DB: mockDB}
	uid, err := u.RegisterUser(ctx, entity)

	// 断言邮箱无效错误
	assert.Equal(t, int64(0), uid)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestRegisterUser_EncryptPasswordError(t *testing.T) {
	mockDB := new(mocks.PersistencePort)
	ctx := context.Background()
	entity := &entities.User{
		UserName: "testuser",
		Email:    "valid@example.com",
		Password: "", // 空密码导致加密失败
	}

	mockDB.On("IsUserExist", ctx, entity.UserName).Return(false, nil)

	u := &usecase.UseCase{DB: mockDB}
	uid, err := u.RegisterUser(ctx, entity)

	// 断言密码加密错误
	assert.Equal(t, int64(0), uid)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestRegisterUser_CreateUserError(t *testing.T) {
	mockDB := new(mocks.PersistencePort)
	ctx := context.Background()
	entity := &entities.User{
		UserName: "testuser",
		Email:    "valid@example.com",
		Password: "Ppassword1.",
	}
	expectedErr := fmt.Errorf("create error")

	// 模拟IsUserExist和CreateUser
	mockDB.On("IsUserExist", ctx, entity.UserName).Return(false, nil)
	mockDB.On("CreateUser", ctx, entity).Return(expectedErr)

	u := &usecase.UseCase{DB: mockDB}
	uid, err := u.RegisterUser(ctx, entity)

	// 断言创建用户错误
	assert.Equal(t, int64(0), uid)
	// assert.ErrorContains(t, err, "create user failed")
	assert.ErrorIs(t, err, expectedErr)
	mockDB.AssertExpectations(t)
}

func TestRegisterUser_Success(t *testing.T) {
	mockDB := new(mocks.PersistencePort)
	ctx := context.Background()
	entity := &entities.User{
		UserName: "testuser",
		Email:    "valid@example.com",
		Password: "Password1.",
	}
	expectedUID := int64(123)

	// 模拟成功流程
	mockDB.On("IsUserExist", ctx, entity.UserName).Return(false, nil)
	mockDB.On("CreateUser", ctx, mock.AnythingOfType("*entities.User")).Return(nil).Run(func(args mock.Arguments) {
		entity.Uid = expectedUID // 模拟数据库生成UID
	})

	u := &usecase.UseCase{DB: mockDB}
	uid, err := u.RegisterUser(ctx, entity)

	// 断言成功返回UID
	assert.NoError(t, err)
	assert.Equal(t, expectedUID, uid)
	mockDB.AssertExpectations(t)
}
