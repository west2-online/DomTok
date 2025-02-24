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

	api "github.com/west2-online/DomTok/app/gateway/model/api/user"
	"github.com/west2-online/DomTok/app/gateway/model/model"
	"github.com/west2-online/DomTok/kitex_gen/user"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func InitUserRPC() {
	c, err := client.InitUserRPC()
	if err != nil {
		logger.Fatalf("api.rpc.user InitUserRPC failed, err is %v", err)
	}
	userClient = *c
}

func RegisterRPC(ctx context.Context, req *user.RegisterRequest) (response *api.RegisterResponse, err error) {
	resp, err := userClient.Register(ctx, req)
	// 这里的 err 是属于 RPC 间调用的错误，例如 network error
	// 而业务错误则是封装在 resp.base 当中的
	if err != nil {
		logger.Errorf("RegisterRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}

	response = &api.RegisterResponse{UID: resp.UserID}
	return response, nil
}

func LoginRPC(ctx context.Context, req *user.LoginRequest) (response *api.LoginResponse, err error) {
	resp, err := userClient.Login(ctx, req)
	if err != nil {
		logger.Errorf("LoginRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}

	response = &api.LoginResponse{
		User: &model.UserInfo{
			UserId: resp.User.UserId,
			Name:   resp.User.Name,
		},
	}

	return response, nil
}
