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
		logger.Fatalf("api.rpc.Commodity InitCommodityRPC failed, err  %v", err)
	}
	commodityClient = *c
}

func InitCommodityStreamClientRPC() {
	c, err := client.InitCommodityStreamClientRPC()
	if err != nil {
		logger.Fatalf("api.rpc.Commodity InitCommodityStreamClientRPC failed, err  %v", err)
	}
	commodityStreamClient = *c
}

func CreateSpuRPC(ctx context.Context, req *commodity.CreateSpuReq, files [][]byte) (id int64, err error) {
	stream, err := commodityStreamClient.CreateSpu(ctx)
	if err != nil {
		logger.Errorf("rpc.CreateSpuRPC CreateSpu failed, err  %v", err)
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	err = stream.Send(req)
	if err != nil {
		logger.Errorf("rpc.CreateSpuRPC SendReq failed, err  %v", err)
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	for _, file := range files {
		err = stream.Send(&commodity.CreateSpuReq{GoodsHeadDrawing: file})
		if err != nil {
			logger.Errorf("rpc.CreateSpuRPC CreateSpu failed, err  %v", err)
			return 0, errno.InternalServiceError.WithMessage(err.Error())
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("rpc.CreateSpuRPC CreateSpu failed, err  %v", err)
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	return resp.SpuID, nil
}

func UpdateSpuRPC(ctx context.Context, req *commodity.UpdateSpuReq, files [][]byte) (err error) {
	stream, err := commodityStreamClient.UpdateSpu(ctx)
	if err != nil {
		logger.Errorf("rpc.UpdateSpuRPC UpdateSpu failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	err = stream.Send(req)
	if err != nil {
		logger.Errorf("rpc.UpdateSpuRPC SendReq failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	for _, file := range files {
		err = stream.Send(&commodity.UpdateSpuReq{GoodsHeadDrawing: file})
		if err != nil {
			logger.Errorf("rpc.UpdateSpuRPC UpdateSpu failed, err  %v", err)
			return errno.InternalServiceError.WithMessage(err.Error())
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("rpc.UpdateSpuRPC UpdateSpu failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	return nil
}

func DeleteSpuRPC(ctx context.Context, req *commodity.DeleteSpuReq) (err error) {
	resp, err := commodityClient.DeleteSpu(ctx, req)
	if err != nil {
		logger.Errorf("rpc.DeleteSpuRPC DeleteSpu failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func CreateSkuRPC(ctx context.Context, req *commodity.CreateSkuReq, files [][]byte) (skuID int64, err error) {
	stream, err := commodityStreamClient.CreateSku(ctx)
	if err != nil {
		logger.Errorf("rpc.CreateSkuRPC: RPC called failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	err = stream.Send(req)
	if err != nil {
		logger.Errorf("rpc.CreateSkuRPC: RPC called failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	for _, file := range files {
		err = stream.Send(&commodity.CreateSkuReq{StyleHeadDrawing: file})
		if err != nil {
			logger.Errorf("rpc.CreateSkuRPC CreateSku failed, err  %v", err)
			return 0, errno.InternalServiceError.WithMessage(err.Error())
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("rpc.CreateSkuRPC: RPC called failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}
	return resp.SkuID, nil
}

func UpdateSkuRPC(ctx context.Context, req *commodity.UpdateSkuReq, files [][]byte) (err error) {
	stream, err := commodityStreamClient.UpdateSku(ctx)
	if err != nil {
		logger.Errorf("rpc.UpdateSkuRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	err = stream.Send(req)
	if err != nil {
		logger.Errorf("rpc.UpdateSkuRPC SendReq failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	for _, file := range files {
		err = stream.Send(&commodity.UpdateSkuReq{StyleHeadDrawing: file})
		if err != nil {
			logger.Errorf("rpc.UpdateSkuRPC UpdateSku failed, err  %v", err)
			return errno.InternalServiceError.WithMessage(err.Error())
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("rpc.UpdateSkuRPC UpdateSku failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	return nil
}

func DeleteSkuRPC(ctx context.Context, req *commodity.DeleteSkuReq) (err error) {
	resp, err := commodityClient.DeleteSku(ctx, req)
	if err != nil {
		logger.Errorf("DeleteSkuRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

func ViewSkuImageRPC(ctx context.Context, req *commodity.ViewSkuImageReq) (images []*model.SkuImage, err error) {
	resp, err := commodityClient.ViewSkuImage(ctx, req)
	if err != nil {
		logger.Errorf("ViewSkuImageRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Images, nil
}

func ViewSkuRPC(ctx context.Context, req *commodity.ViewSkuReq) (sku []*model.Sku, err error) {
	resp, err := commodityClient.ViewSku(ctx, req)
	if err != nil {
		logger.Errorf("ViewSkuRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Skus, nil
}

func UploadSkuAttrRPC(ctx context.Context, req *commodity.UploadSkuAttrReq) (err error) {
	resp, err := commodityClient.UploadSkuAttr(ctx, req)
	if err != nil {
		logger.Errorf("UploadSkuAttrRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

func ListSkuInfo(ctx context.Context, req *commodity.ListSkuInfoReq) (skus []*model.SkuInfo, err error) {
	resp, err := commodityClient.ListSkuInfo(ctx, req)
	if err != nil {
		logger.Errorf("ListSkuInfo: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.SkuInfos, nil
}
