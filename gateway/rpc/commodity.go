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
	"github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func InitCommodityRPC() {
	c, err := client.InitCommodityRPC()
	if err != nil {
		logger.Fatalf("api.rpc.commodity InitUserRPC failed, err is %v", err)
	}
	commodityClient = *c
}

func CreateCategoryRPC(ctx context.Context, req *commodity.CreateCategoryReq) (int64, error) {
	resp, err := commodityClient.CreateCategory(ctx, req)
	if err != nil {
		logger.Errorf("CreateCategoryRPC: RPC called failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.CategoryID, nil
}

func DeleteCategoryRPC(ctx context.Context, req *commodity.DeleteCategoryReq) (*model.BaseResp, error) {
	resp, err := commodityClient.DeleteCategory(ctx, req)
	if err != nil {
		logger.Errorf("DeleteCategoryRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Base, nil
}

func UpdateCategoryRPC(ctx context.Context, req *commodity.UpdateCategoryReq) (*model.BaseResp, error) {
	resp, err := commodityClient.UpdateCategory(ctx, req)
	if err != nil {
		logger.Errorf("UpdateCategoryRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Base, nil
}

func ViewCategoryRPC(ctx context.Context, req *commodity.ViewCategoryReq) ([]*model.CategoryInfo, error) {
	resp, err := commodityClient.ViewCategory(ctx, req)
	if err != nil {
		logger.Errorf("ViewCategoryRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.CategoryInfo, nil
}
