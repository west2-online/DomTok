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

	"github.com/west2-online/DomTok/app/order/usecase"
	"github.com/west2-online/DomTok/kitex_gen/order"
)

type OrderHandler struct {
	useCase usecase.OrderUseCase
}

func NewOrderHandler(useCase usecase.OrderUseCase) order.OrderService {
	return &OrderHandler{useCase: useCase}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderReq) (resp *order.CreateOrderResp, err error) {
	resp, err = h.useCase.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *OrderHandler) ViewOrderList(ctx context.Context, req *order.ViewOrderListReq) (resp *order.ViewOrderListResp, err error) {
	resp, err = h.useCase.ViewOrderList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *OrderHandler) ViewOrder(ctx context.Context, req *order.ViewOrderReq) (resp *order.ViewOrderResp, err error) {
	resp, err = h.useCase.ViewOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *order.CancelOrderReq) (resp *order.CancelOrderResp, err error) {
	resp, err = h.useCase.CancelOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *OrderHandler) ChangeDeliverAddress(ctx context.Context, req *order.ChangeDeliverAddressReq) (resp *order.ChangeDeliverAddressResp, err error) {
	resp, err = h.useCase.ChangeDeliverAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, req *order.DeleteOrderReq) (resp *order.DeleteOrderResp, err error) {
	resp, err = h.useCase.DeleteOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
