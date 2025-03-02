package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/west2-online/DomTok/app/user/domain/model"
	"github.com/west2-online/DomTok/app/user/domain/service"
	"github.com/west2-online/DomTok/app/user/usecase"
	"github.com/west2-online/DomTok/app/user/usecase/mocks"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

func TestUseCase_RegisterUser(t *testing.T) {
	mockDB := new(mocks.UserDB)
	mockSf, _ := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
	mockService := service.NewUserService(mockDB, mockSf)
	uc := usecase.NewUserCase(mockDB, mockService)

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
