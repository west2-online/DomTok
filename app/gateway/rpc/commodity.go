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

	if !utils.IsSuccess(resp.Base) {
		return 0, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
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

	resp, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("rpc.UpdateSpuRPC UpdateSpu failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
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

func DeleteSpuImageRPC(ctx context.Context, req *commodity.DeleteSpuImageReq) (err error) {
	resp, err := commodityClient.DeleteSpuImage(ctx, req)
	if err != nil {
		logger.Errorf("rpc.DeleteSpuImage DeleteSpuImageRPC failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func CreateSpuImageRPC(ctx context.Context, req *commodity.CreateSpuImageReq, files [][]byte) (id int64, err error) {
	stream, err := commodityStreamClient.CreateSpuImage(ctx)
	if err != nil {
		logger.Errorf("rpc.CreateSpuImageRPC CreateSpuImage failed, err  %v", err)
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	err = stream.Send(req)
	if err != nil {
		logger.Errorf("rpc.CreateSpuImageRPC SendReq failed, err  %v", err)
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	for _, file := range files {
		err = stream.Send(&commodity.CreateSpuImageReq{
			Data: file,
		})
		if err != nil {
			logger.Errorf("rpc.CreateSpuImageRPC CreateSpuImage failed, err  %v", err)
			return 0, errno.InternalServiceError.WithMessage(err.Error())
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("rpc.CreateSpuImageRPC CreateSpuImage failed, err  %v", err)
		return 0, errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return 0, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}

	return resp.ImageID, nil
}

func UpdateSpuImageRPC(ctx context.Context, req *commodity.UpdateSpuImageReq, files [][]byte) (err error) {
	stream, err := commodityStreamClient.UpdateSpuImage(ctx)
	if err != nil {
		logger.Errorf("rpc.UpdateSpuImageRPC UpdateSpuImage failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	err = stream.Send(req)
	if err != nil {
		logger.Errorf("rpc.UpdateSpuImageRPC SendReq failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	for _, file := range files {
		err = stream.Send(&commodity.UpdateSpuImageReq{
			Data: file,
		})
		if err != nil {
			logger.Errorf("rpc.UpdateSpuImageRPC UpdateSpuImage failed, err  %v", err)
			return errno.InternalServiceError.WithMessage(err.Error())
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("rpc.UpdateSpuImageRPC UpdateSpuImage failed, err  %v", err)
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}

	return nil
}
