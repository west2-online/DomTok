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

package rpc

import (
	"context"

	"github.com/west2-online/DomTok/app/user/controllers/rpc/pack"
	"github.com/west2-online/DomTok/app/user/domain"
	"github.com/west2-online/DomTok/kitex_gen/user"
)

// UserHandler 实现 idl 中定义的 RPC 接口
type UserHandler struct {
	useCase domain.UserUseCase
}

func NewUserHandler(useCase domain.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (r *user.RegisterResponse, err error) {
	r = new(user.RegisterResponse)
	u := &domain.User{
		UserName: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	var uid int64
	if uid, err = h.useCase.RegisterUser(ctx, u); err != nil {
		return
	}

	r.UserID = uid
	return
}

func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (r *user.LoginResponse, err error) {
	r = new(user.LoginResponse)

	u := &domain.User{
		UserName: req.Username,
		Password: req.Password,
	}

	ans, err := h.useCase.Login(ctx, u)
	if err != nil {
		return
	}

	r.User = pack.BuildUser(ans)
	return
}
