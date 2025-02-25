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

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/app/order/domain/service"
)

// OrderUseCase 定义在 usecase 层的接口
type OrderUseCase interface {
	// CreateOrder 返回 orderID 和 error
	CreateOrder(ctx context.Context, addressID int64, goods []*model.BaseOrderGoods) (int64, error)
	ViewOrderList(ctx context.Context, page, size int32) ([]*model.Order, []*model.OrderGoods, int32, error)
	ViewOrder(ctx context.Context, orderID int64) (*model.Order, []*model.OrderGoods, error)
	CancelOrder(ctx context.Context, orderID int64) error
	ChangeDeliverAddress(ctx context.Context, orderID, addressID int64, addressInfo string) error
	DeleteOrder(ctx context.Context, orderID int64) error
	IsOrderExist(ctx context.Context, orderID int64) (bool, error)
}

// useCase 实现了 OrderUseCase 接口
type useCase struct {
	db  repository.OrderDB
	svc *service.OrderService
	rpc repository.RPC
}

func NewOrderCase(db repository.OrderDB, svc *service.OrderService, rpc repository.RPC) OrderUseCase {
	return &useCase{
		db:  db,
		svc: svc,
		rpc: rpc,
	}
}
