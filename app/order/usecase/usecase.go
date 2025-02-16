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

	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/app/order/domain/service"
	"github.com/west2-online/DomTok/kitex_gen/order"
)

// OrderUseCase 定义在 usecase 层的接口
type OrderUseCase interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderReq) (*order.CreateOrderResp, error)
	ViewOrderList(ctx context.Context, req *order.ViewOrderListReq) (*order.ViewOrderListResp, error)
	ViewOrder(ctx context.Context, req *order.ViewOrderReq) (*order.ViewOrderResp, error)
	CancelOrder(ctx context.Context, req *order.CancelOrderReq) (*order.CancelOrderResp, error)
	ChangeDeliverAddress(ctx context.Context, req *order.ChangeDeliverAddressReq) (*order.ChangeDeliverAddressResp, error)
	DeleteOrder(ctx context.Context, req *order.DeleteOrderReq) (*order.DeleteOrderResp, error)
}

// useCase 实现了 OrderUseCase 接口
type useCase struct {
	db  repository.OrderDB
	svc *service.OrderService
}

func NewOrderCase(db repository.OrderDB, svc *service.OrderService) OrderUseCase {
	return &useCase{
		db:  db,
		svc: svc,
	}
}
