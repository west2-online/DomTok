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

	"github.com/west2-online/DomTok/app/user/domain"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
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

func (svc *UserService) CreateUser(ctx context.Context, u *domain.User) error {
	u.Uid = svc.nextID()

	if err := svc.db.CreateUser(ctx, u); err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}

	return nil
}

func (svc *UserService) nextID() int64 {
	id, _ := svc.sf.NextVal()
	return id
}
