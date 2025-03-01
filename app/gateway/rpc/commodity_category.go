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

	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func CreateCategory(ctx context.Context, req *commodity.CreateCategoryReq) (r *commodity.CreateCategoryResp, err error) {
	resp, err := commodityClient.CreateCategory(ctx, req)
	if err != nil {
		logger.Errorf("rpc.CreateCategory CreateCategory failed, err  %v", err)
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp, nil
}

func DeleteCategory(ctx context.Context, req *commodity.DeleteCategoryReq) (r *commodity.DeleteCategoryResp, err error) {
	resp, err := commodityClient.DeleteCategory(ctx, req)
	if err != nil {
		logger.Errorf("rpc.DeleteCategory DeleteCategory failed, err  %v", err)
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp, nil
}

func ViewCategory(ctx context.Context, req *commodity.ViewCategoryReq) (r *commodity.ViewCategoryResp, err error) {
	resp, err := commodityClient.ViewCategory(ctx, req)
	if err != nil {
		logger.Errorf("rpc.ViewCategory ViewCategory failed, err  %v", err)
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp, nil
}

func UpdateCategory(ctx context.Context, req *commodity.UpdateCategoryReq) (r *commodity.UpdateCategoryResp, err error) {
	resp, err := commodityClient.UpdateCategory(ctx, req)
	if err != nil {
		logger.Errorf("rpc.UpdateCategory UpdateCategory failed, err  %v", err)
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp, nil
}
