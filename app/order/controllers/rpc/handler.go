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

	"github.com/west2-online/DomTok/app/order/controllers/rpc/pack"
	"github.com/west2-online/DomTok/app/order/usecase"
	idlmodel "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/kitex_gen/order"
)

type OrderHandler struct {
	useCase usecase.OrderUseCase
}

func NewOrderHandler(useCase usecase.OrderUseCase) order.OrderService {
	return &OrderHandler{useCase: useCase}
}

// TODO：需要实现
func (h *OrderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderReq) (resp *order.CreateOrderResp, err error) {
	resp = new(order.CreateOrderResp)
	return resp, nil
}

func (h *OrderHandler) ViewOrderList(ctx context.Context, req *order.ViewOrderListReq) (resp *order.ViewOrderListResp, err error) {
	resp = new(order.ViewOrderListResp)
	orders, goods, total, err := h.useCase.ViewOrderList(ctx, req.GetPage(), req.GetSize())
	if err != nil {
		return resp, err
	}

	idlOrders := make([]*idlmodel.BaseOrderWithGoods, 0, len(orders))
	for i, o := range orders {
		idlOrders = append(idlOrders, &idlmodel.BaseOrderWithGoods{
			Order: pack.BuildBaseOrder(o),
			Goods: []*idlmodel.BaseOrderGoods{pack.BuildBaseOrderGoods(goods[i])},
		})
	}

	resp.Total = total
	resp.OrderList = idlOrders
	return resp, nil
}

func (h *OrderHandler) ViewOrder(ctx context.Context, req *order.ViewOrderReq) (resp *order.ViewOrderResp, err error) {
	resp = new(order.ViewOrderResp)
	o, goods, err := h.useCase.ViewOrder(ctx, req.GetOrderID())
	if err != nil {
		return resp, err
	}
	resp.Data = pack.BuildOrderWithGoods(o, goods)
	return resp, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *order.CancelOrderReq) (resp *order.CancelOrderResp, err error) {
	resp = new(order.CancelOrderResp)
	return resp, h.useCase.CancelOrder(ctx, req.GetOrderID())
}

func (h *OrderHandler) ChangeDeliverAddress(ctx context.Context, req *order.ChangeDeliverAddressReq) (resp *order.ChangeDeliverAddressResp, err error) {
	resp = new(order.ChangeDeliverAddressResp)
	return resp, h.useCase.ChangeDeliverAddress(ctx, req.GetOrderID(), req.GetAddressID(), req.GetAddressInfo())
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, req *order.DeleteOrderReq) (resp *order.DeleteOrderResp, err error) {
	resp = new(order.DeleteOrderResp)
	return resp, h.useCase.DeleteOrder(ctx, req.GetOrderID())
}

func (h *OrderHandler) IsOrderExist(ctx context.Context, req *order.IsOrderExistReq) (resp *order.IsOrderExistResp, err error) {
	resp = new(order.IsOrderExistResp)
	resp.Exist, err = h.useCase.IsOrderExist(ctx, req.GetOrderID())
	return resp, err
}
