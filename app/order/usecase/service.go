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

package usecase

import (
	"context"
	"fmt"

	"github.com/west2-online/DomTok/app/order/controllers/rpc/pack"
	"github.com/west2-online/DomTok/app/order/domain/model"
	modelrpc "github.com/west2-online/DomTok/kitex_gen/model"
	orderrpc "github.com/west2-online/DomTok/kitex_gen/order"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

// CreateOrder 创建订单
func (uc *useCase) CreateOrder(ctx context.Context, req *orderrpc.CreateOrderReq) (*orderrpc.CreateOrderResp, error) {
	// 1. 构建订单模型
	order := &model.Order{
		AddressID:   req.AddressID,
		AddressInfo: req.AddressInfo,
		Status:      0, // 待支付
	}

	// 2. 构建订单商品模型
	var goods []*model.OrderGoods
	for _, g := range req.BaseOrderGoods {
		goods = append(goods, &model.OrderGoods{
			GoodsID:  g.GoodsID,
			Quantity: int32(g.PurchaseQuantity),
		})
	}

	// 3. 调用服务创建订单
	orderID, err := uc.svc.CreateOrder(ctx, order, goods)
	if err != nil {
		return nil, fmt.Errorf("create order failed: %w", err)
	}

	return &orderrpc.CreateOrderResp{
		Base:    &modelrpc.BaseResp{Code: errno.SuccessCode, Msg: "success"},
		OrderID: orderID,
	}, nil
}

// ViewOrderList 获取订单列表
func (uc *useCase) ViewOrderList(ctx context.Context, req *orderrpc.ViewOrderListReq) (*orderrpc.ViewOrderListResp, error) {
	logger.Infof("[ViewOrderList] 开始处理订单列表查询")
	_, goods, total, err := uc.svc.ViewOrderList(ctx, 0, req.Page, req.Size)
	if err != nil {
		return nil, fmt.Errorf("view order list failed: %w", err)
	}

	// 构建响应
	orderGoods := make([]*modelrpc.OrderGoods, 0, len(goods))
	for _, g := range goods {
		orderGoods = append(orderGoods, pack.BuildOrderGoods(g))
	}

	return &orderrpc.ViewOrderListResp{
		Base:       &modelrpc.BaseResp{Code: errno.SuccessCode, Msg: "success"},
		Total:      total,
		OrderGoods: orderGoods,
	}, nil
}

// ViewOrder 获取订单详情
func (uc *useCase) ViewOrder(ctx context.Context, req *orderrpc.ViewOrderReq) (*orderrpc.ViewOrderResp, error) {
	order, goods, err := uc.svc.ViewOrder(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("view order failed: %w", err)
	}

	// 构建响应
	orderGoods := make([]*modelrpc.OrderGoods, 0, len(goods))
	for _, g := range goods {
		orderGoods = append(orderGoods, pack.BuildOrderGoods(g))
	}

	return &orderrpc.ViewOrderResp{
		Base:        &modelrpc.BaseResp{Code: errno.SuccessCode, Msg: "success"},
		AddressID:   order.AddressID,
		AddressInfo: order.AddressInfo,
		Status:      constants.GetOrderStatusMsg(order.Status),
		OrderGoods:  orderGoods,
	}, nil
}

// CancelOrder 取消订单
func (uc *useCase) CancelOrder(ctx context.Context, req *orderrpc.CancelOrderReq) (*orderrpc.CancelOrderResp, error) {
	if err := uc.svc.CancelOrder(ctx, req.OrderID); err != nil {
		return nil, fmt.Errorf("cancel order failed: %w", err)
	}

	return &orderrpc.CancelOrderResp{
		Base: &modelrpc.BaseResp{Code: errno.SuccessCode, Msg: "success"},
	}, nil
}

// ChangeDeliverAddress 更改配送地址
func (uc *useCase) ChangeDeliverAddress(ctx context.Context, req *orderrpc.ChangeDeliverAddressReq) (*orderrpc.ChangeDeliverAddressResp, error) {
	if err := uc.svc.ChangeDeliverAddress(ctx, req.OrderID, req.AddressID, req.AddressInfo); err != nil {
		return nil, fmt.Errorf("change deliver address failed: %w", err)
	}

	return &orderrpc.ChangeDeliverAddressResp{
		Base: &modelrpc.BaseResp{Code: errno.SuccessCode, Msg: "success"},
	}, nil
}

// DeleteOrder 删除订单
func (uc *useCase) DeleteOrder(ctx context.Context, req *orderrpc.DeleteOrderReq) (*orderrpc.DeleteOrderResp, error) {
	if err := uc.svc.DeleteOrder(ctx, req.OrderID); err != nil {
		return nil, fmt.Errorf("delete order failed: %w", err)
	}

	return &orderrpc.DeleteOrderResp{
		Base: &modelrpc.BaseResp{Code: errno.SuccessCode, Msg: "success"},
	}, nil
}
