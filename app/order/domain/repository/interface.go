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

package repository

import (
	"context"

	"github.com/west2-online/DomTok/app/order/domain/model"
)

// OrderDB 表示订单模块的持久化存储接口
type OrderDB interface {
	IsOrderExist(ctx context.Context, orderID int64) (bool, error)
	CreateOrder(ctx context.Context, order *model.Order) error
	CreateOrderGoods(ctx context.Context, orderGoods []*model.OrderGoods) error
	GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error)
	GetOrderGoodsByOrderID(ctx context.Context, orderID int64) ([]*model.OrderGoods, error)
	GetOrdersByUserID(ctx context.Context, userID int64, page, size int32) ([]*model.Order, int32, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status int32) error
	UpdateOrderAddress(ctx context.Context, orderID int64, addressID int64, addressInfo string) error
	DeleteOrder(ctx context.Context, orderID int64) error
}
