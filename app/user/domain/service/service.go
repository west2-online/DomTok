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

	"golang.org/x/crypto/bcrypt"

	"github.com/west2-online/DomTok/app/user/domain/model"
	"github.com/west2-online/DomTok/config"
	metadata "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/utils"
)

func (svc *UserService) EncryptPassword(pwd string) (string, error) {
	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(pwd), constants.UserDefaultEncryptPasswordCost)
	if err != nil {
		return "", errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("encrypt password failed, pwd: %s, err: %v", pwd, err))
	}
	return string(passwordDigest), nil
}

func (svc *UserService) CheckPassword(passwordDigest, password string) error {
	if bcrypt.CompareHashAndPassword([]byte(passwordDigest), []byte(password)) != nil {
		return errno.NewErrNo(errno.ServiceWrongPassword, "wrong password")
	}
	return nil
}

func (svc *UserService) CreateUser(ctx context.Context, u *model.User) (int64, error) {
	uid, err := svc.db.CreateUser(ctx, u)
	if err != nil {
		return 0, fmt.Errorf("create user failed: %w", err)
	}
	return uid, nil
}

func (svc *UserService) GetAddress(ctx context.Context, addressID int64) (address *model.Address, err error) {
	address, err = svc.db.GetAddressInfo(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("domain.svc.GetAddress failed: %w", err)
	}
	return address, nil
}

func (svc *UserService) AddAddress(ctx context.Context, address *model.Address) (addressID int64, err error) {
	addressID, err = svc.db.CreateAddress(ctx, address)
	if err != nil {
		return 0, fmt.Errorf("domain.svc.AddAddress failed: %w", err)
	}
	return addressID, nil
}

func (svc *UserService) UserLogin(ctx context.Context, uid int64) error {
	key := svc.cache.UserLogOutKey(uid)
	var token string
	var err error
	exist := svc.cache.IsExist(ctx, key)
	if exist {
		oldToken, err := svc.cache.GetToken(ctx, key)
		if err != nil {
			return fmt.Errorf("domain.svc.UserLogOut failed: %w", err)
		}
		_, _, err = utils.CheckToken(oldToken)
		if err != nil {
			return fmt.Errorf("domain.svc.UserLogOut failed: %w", err)
		}
		return nil
	} else {
		token, err = utils.CreateToken(constants.TypeUserLoginToken, uid)
		if err != nil {
			return errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("create token failed, err: %v", err))
		}
	}

	err = svc.cache.SetUserLogOut(ctx, key, token)
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "domain.svc.UserLogin failed: %v", err)
	}
	return nil
}

func (svc *UserService) UserBaned(ctx context.Context, uid int64) error {
	u, err := svc.db.GetUserById(ctx, uid)
	if err != nil {
		return fmt.Errorf("domain.svc.UserBaned failed: %w", err)
	}

	me, err := metadata.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("domain.svc.UserBaned failed: %w", err)
	}

	if me == uid {
		return errno.NewErrNo(errno.ParamVerifyErrorCode, "can not do this at self")
	}

	myInfo, err := svc.db.GetUserById(ctx, me)
	if err != nil {
		return fmt.Errorf("domain.svc.UserBaned failed: %w", err)
	}

	if myInfo.Role != constants.UserAdministrator {
		return errno.NewErrNo(errno.AuthNoOperatePermissionCode, "permission denied")
	}

	if u.Role == constants.UserAdministrator {
		return errno.NewErrNo(errno.AuthNoOperatePermissionCode, "domain.svc.UserBaned failed: role is administrator")
	}

	key := svc.cache.UserBanedKey(uid)
	exist := svc.cache.IsExist(ctx, key)
	if !exist {
		err = svc.cache.SetUserBaned(ctx, key)
		if err != nil {
			return fmt.Errorf("domain.svc.UserBaned failed: %w", err)
		}
		return nil
	} else {
		return errno.NewErrNo(errno.RepeatedOperation, "domain.svc.UserBaned failed, already banned user")
	}
}

func (svc *UserService) LiftUserBaned(ctx context.Context, uid int64) error {
	_, err := svc.db.GetUserById(ctx, uid)
	if err != nil {
		return fmt.Errorf("domain.svc.LiftUserBaned failed: %w", err)
	}

	me, err := metadata.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("domain.svc.LiftUserBaned failed: %w", err)
	}

	myInfo, err := svc.db.GetUserById(ctx, me)
	if err != nil {
		return fmt.Errorf("domain.svc.LiftUserBaned failed: %w", err)
	}

	if myInfo.Role != constants.UserAdministrator {
		return errno.NewErrNo(errno.AuthNoOperatePermissionCode, "permission denied")
	}

	key := svc.cache.UserBanedKey(uid)
	exist := svc.cache.IsExist(ctx, key)
	if exist {
		err = svc.cache.DeleteUserBaned(ctx, key)
		if err != nil {
			return fmt.Errorf("domain.svc.LiftUserBaned failed: %w", err)
		}
		return nil
	} else {
		return errno.NewErrNo(errno.RepeatedOperation, "domain.svc.LiftUserBaned failed, already normal user")
	}
}

func (svc *UserService) Logout(ctx context.Context) error {
	uid, err := metadata.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("domain.svc.Logout failed: %w", err)
	}
	key := svc.cache.UserLogOutKey(uid)
	exist := svc.cache.IsExist(ctx, key)
	if exist {
		err = svc.cache.DeleteUserLogOut(ctx, key)
		if err != nil {
			return fmt.Errorf("domain.svc.Logout failed: %w", err)
		}
		return nil
	} else {
		return errno.NewErrNo(errno.RepeatedOperation, "domain.svc.Logout failed, already logout")
	}
}

func (svc *UserService) SetAdministrator(ctx context.Context, uid int64, password []byte, action int) error {
	err := bcrypt.CompareHashAndPassword([]byte(config.Administrator.Secret), password)
	if err != nil {
		return errno.NewErrNo(errno.AuthNoOperatePermissionCode, "wrong secret")
	}

	_, err = svc.db.GetUserById(ctx, uid)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "domain.svc.SetAdministrator failed: %v", err)
	}

	if err = svc.Verify(svc.VerifyAction(action)); err != nil {
		return errno.NewErrNo(errno.ParamVerifyErrorCode, "action type error")
	}

	err = svc.db.UpdateUser(ctx, &model.User{
		Uid:  uid,
		Role: action,
	})
	if err != nil {
		return errno.NewErrNo(errno.InternalServiceErrorCode, "domain.svc.SetAdministrator failed")
	}
	return nil
}

func (svc *UserService) IsBaned(ctx context.Context, uid int64) (bool, error) {
	key := svc.cache.UserBanedKey(uid)
	exist := svc.cache.IsExist(ctx, key)
	return exist, nil
}
