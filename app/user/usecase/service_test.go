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
	contextPkg "github.com/west2-online/DomTok/pkg/base/context"
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

func TestUseCase_ListAddress(t *testing.T) {
	mockDB := new(mocks.UserDB)
	mockCache := new(redis.Client)
	c := cache.NewUserCache(mockCache)
	mockSf, _ := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
	mockService := service.NewUserService(mockDB, mockSf, c)
	uc := usecase.NewUserCase(mockDB, mockService, c)

	uid := int64(1001)
	ctx := contextPkg.WithLoginData(context.Background(), uid)

	addresses := []*model.Address{
		{AddressID: 1, Uid: uid, Province: "P1", City: "C1", Detail: "D1"},
		{AddressID: 2, Uid: uid, Province: "P2", City: "C2", Detail: "D2"},
	}

	mockDB.On("ListAddress", mock.Anything, uid, 1, 10).Return(addresses, nil)

	res, err := uc.ListAddress(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(1), res[0].AddressID)

	mockDB.AssertExpectations(t)
}

func TestUseCase_DeleteAddress(t *testing.T) {
	mockDB := new(mocks.UserDB)
	mockCache := new(redis.Client)
	c := cache.NewUserCache(mockCache)
	mockSf, _ := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
	mockService := service.NewUserService(mockDB, mockSf, c)
	uc := usecase.NewUserCase(mockDB, mockService, c)

	uid := int64(1001)
	ctx := contextPkg.WithLoginData(context.Background(), uid)

	addressID := int64(1)
	address := &model.Address{AddressID: addressID, Uid: uid}

	// Success case
	mockDB.On("GetAddressInfo", mock.Anything, addressID).Return(address, nil).Once()
	mockDB.On("DeleteAddress", mock.Anything, addressID).Return(nil).Once()

	err := uc.DeleteAddress(ctx, addressID)
	assert.NoError(t, err)

	// Fail case (permission denied)
	otherUid := int64(2002)
	otherAddress := &model.Address{AddressID: 2, Uid: otherUid}
	mockDB.On("GetAddressInfo", mock.Anything, int64(2)).Return(otherAddress, nil).Once()
	// DeleteAddress should NOT be called

	err = uc.DeleteAddress(ctx, 2)
	assert.Error(t, err) // Should be permission denied

	mockDB.AssertExpectations(t)
}
