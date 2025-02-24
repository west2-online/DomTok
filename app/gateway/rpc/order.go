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

	"github.com/west2-online/DomTok/kitex_gen/order"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func InitOrderRPC() {
	c, err := client.InitOrderRPC()
	if err != nil {
		logger.Fatalf("api.rpc.order InitOrderRPC failed, err is %v", err)
	}
	orderClient = *c
}

// CreateOrderRPC 创建订单RPC调用
func CreateOrderRPC(ctx context.Context, req *order.CreateOrderReq) (orderID int64, err error) {
	resp, err := orderClient.CreateOrder(ctx, req)
	if err != nil {
		logger.Errorf("CreateOrderRPC: RPC called failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.OrderID, nil
}

// ViewOrderListRPC 查看订单列表RPC调用
func ViewOrderListRPC(ctx context.Context, req *order.ViewOrderListReq) (*order.ViewOrderListResp, error) {
	resp, err := orderClient.ViewOrderList(ctx, req)
	if err != nil {
		logger.Errorf("ViewOrderListRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp, nil
}

// ViewOrderRPC 查看订单详情RPC调用
func ViewOrderRPC(ctx context.Context, req *order.ViewOrderReq) (*order.ViewOrderResp, error) {
	resp, err := orderClient.ViewOrder(ctx, req)
	if err != nil {
		logger.Errorf("ViewOrderRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp, nil
}

// CancelOrderRPC 取消订单RPC调用
func CancelOrderRPC(ctx context.Context, req *order.CancelOrderReq) error {
	resp, err := orderClient.CancelOrder(ctx, req)
	if err != nil {
		logger.Errorf("CancelOrderRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

// ChangeDeliverAddressRPC 修改配送地址RPC调用
func ChangeDeliverAddressRPC(ctx context.Context, req *order.ChangeDeliverAddressReq) error {
	resp, err := orderClient.ChangeDeliverAddress(ctx, req)
	if err != nil {
		logger.Errorf("ChangeDeliverAddressRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

// DeleteOrderRPC 删除订单RPC调用
func DeleteOrderRPC(ctx context.Context, req *order.DeleteOrderReq) error {
	resp, err := orderClient.DeleteOrder(ctx, req)
	if err != nil {
		logger.Errorf("DeleteOrderRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

func IsOrderExistRPC(ctx context.Context, req *order.IsOrderExistReq) (bool, error) {
	resp, err := orderClient.IsOrderExist(ctx, req)
	if err != nil {
		logger.Errorf("IsOrderExistRPC: RPC called failed: %v", err.Error())
		return false, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return false, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Exist, nil
}
