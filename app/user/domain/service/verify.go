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
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

type UserVerifyOps func() error

// Verify 通过传来的参数进行一系列的校验
func (svc *UserService) Verify(opts ...UserVerifyOps) error {
	for _, opt := range opts {
		if err := opt(); err != nil {
			return err
		}
	}
	return nil
}

// VerifyEmail 返回一个校验 email 格式的函数, 不应单独使用, 应结合 Verify
func (svc *UserService) VerifyEmail(email string) UserVerifyOps {
	return func() error {
		if !svc.emailRe.MatchString(email) {
			return errno.NewErrNo(errno.ParamVerifyErrorCode, "wrong email format")
		}
		return nil
	}
}

// VerifyPassword 返回一个校验 password 长度的函数, 不应单独使用, 应结合 Verify
func (svc *UserService) VerifyPassword(pw string) UserVerifyOps {
	return func() error {
		bpw := len([]byte(pw))
		if bpw > constants.UserMaximumPasswordLength || bpw < constants.UserMinimumPasswordLength {
			return errno.NewErrNo(errno.ParamVerifyErrorCode, "password length should be greater than 5 and less than 72")
		}
		return nil
	}
}
