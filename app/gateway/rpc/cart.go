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

	"github.com/west2-online/DomTok/kitex_gen/cart"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func InitCartRPC() {
	c, err := client.InitCartRPC()
	if err != nil {
		logger.Fatalf("api.rpc.cart InitCartRPC failed, err is %v", err)
	}
	cartClient = *c
}

func AddGoodsIntoCartRPC(ctx context.Context, req *cart.AddGoodsIntoCartRequest) (err error) {
	resp, err := cartClient.AddGoodsIntoCart(ctx, req)
	if err != nil {
		logger.Errorf("AddGoodsIntoCartRPC RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		// 对外暴露的信息，实际错误的log已经在rpc server通过mw打印了
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

func ShowCartGoodsRPC(ctx context.Context, req *cart.ShowCartGoodsListRequest) (err error) {
	resp, err := cartClient.ShowCartGoodsList(ctx, req)
	if err != nil {
		logger.Errorf("ShowCartGoodsRPC RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		// 对外暴露的信息，实际错误的log已经在rpc server通过mw打印了
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}
