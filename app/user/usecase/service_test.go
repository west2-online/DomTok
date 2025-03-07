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
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/west2-online/DomTok/app/user/domain/model"
	"github.com/west2-online/DomTok/app/user/domain/service"
	"github.com/west2-online/DomTok/app/user/infrastructure/cache"
	"github.com/west2-online/DomTok/app/user/usecase"
	"github.com/west2-online/DomTok/app/user/usecase/mocks"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

func TestUseCase_RegisterUser(t *testing.T) {
	mockDB := new(mocks.UserDB)
	mockCache := new(redis.Client)
	c := cache.NewUserCache(mockCache)
	mockSf, _ := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
	mockService := service.NewUserService(mockDB, mockSf, c)
	uc := usecase.NewUserCase(mockDB, mockService, c)

	user := &model.User{
		UserName: "testuser",
		Password: "password",
		Email:    "test@example.com",
	}

	mockDB.On("IsUserExist", mock.Anything, user.UserName).Return(false, nil)
	mockDB.On("CreateUser", mock.Anything, mock.Anything).Return(int64(1), nil)

	uid, err := uc.RegisterUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), uid)

	mockDB.AssertExpectations(t)
}
