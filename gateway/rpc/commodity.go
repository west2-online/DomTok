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
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func InitCommodityRPC() {
	c, err := client.InitCommodityRPC()
	if err != nil {
		logger.Fatalf("api.rpc.Commodity InitCommodityRPC failed, err  %v", err)
	}
	commodityClient = *c
}

func CreateSpuRPC(ctx context.Context, req *commodity.CreateSpuReq) (id int64, err error) {
	resp, err := commodityClient.CreateSpu(ctx, req)
	if err != nil {
		logger.Errorf("GetDownloadUrlRPC: RPC called failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.SpuID, nil
}
