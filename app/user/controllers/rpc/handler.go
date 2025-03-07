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
	"github.com/west2-online/DomTok/app/user/domain/model"
	"github.com/west2-online/DomTok/app/user/usecase"
	"github.com/west2-online/DomTok/kitex_gen/user"
	"github.com/west2-online/DomTok/pkg/base"
)

// UserHandler 实现 idl 中定义的 RPC 接口
type UserHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (r *user.RegisterResponse, err error) {
	r = new(user.RegisterResponse)
	u := &model.User{
		UserName: req.Username,
		Password: req.Password,
		Email:    req.Email,
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

	u := &model.User{
		UserName: req.Username,
		Password: req.Password,
	}

	ans, err := h.useCase.Login(ctx, u)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(nil)
	r.User = pack.BuildUser(ans)
	return
}

func (h *UserHandler) GetAddress(ctx context.Context, req *user.GetAddressRequest) (r *user.GetAddressResponse, err error) {
	r = new(user.GetAddressResponse)
	ans, err := h.useCase.GetAddress(ctx, req.AddressId)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(nil)
	r.Address = pack.BuildAddress(ans)
	return
}

func (h *UserHandler) AddAddress(ctx context.Context, req *user.AddAddressRequest) (r *user.AddAddressResponse, err error) {
	r = new(user.AddAddressResponse)
	address := &model.Address{
		Province: req.Province,
		City:     req.City,
		Detail:   req.Detail,
	}
	var addressID int64
	if addressID, err = h.useCase.AddAddress(ctx, address); err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(nil)
	r.AddressID = addressID
	return
}

func (h *UserHandler) BanUser(ctx context.Context, req *user.BanUserReq) (r *user.BanUserResp, err error) {
	r = new(user.BanUserResp)
	err = h.useCase.BanUser(ctx, req.Uid)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(nil)
	return
}

func (h *UserHandler) Logout(ctx context.Context, req *user.LogoutReq) (r *user.LogoutResp, err error) {
	r = new(user.LogoutResp)
	err = h.useCase.LogoutUser(ctx)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(nil)
	return
}

func (h *UserHandler) LiftBandUser(ctx context.Context, req *user.LiftBanUserReq) (r *user.LiftBanUserResp, err error) {
	r = new(user.LiftBanUserResp)
	err = h.useCase.LiftUser(ctx, req.Uid)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}
	r.Base = base.BuildBaseResp(nil)
	return
}
