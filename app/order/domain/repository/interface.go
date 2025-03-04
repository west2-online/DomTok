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
	"time"

	"github.com/west2-online/DomTok/app/order/domain/model"
)

// OrderDB 表示订单模块的持久化存储接口
type OrderDB interface {
	CreateOrder(ctx context.Context, order *model.Order, goods []*model.OrderGoods) error
	CreateOrderGoods(ctx context.Context, orderGoods []*model.OrderGoods) error

	GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error)
	GetOrderGoodsByOrderID(ctx context.Context, orderID int64) ([]*model.OrderGoods, error)
	GetOrdersByUserID(ctx context.Context, userID int64, page, size int32) ([]*model.Order, int32, error)
	GetOrderAndGoods(ctx context.Context, orderID int64) (*model.Order, []*model.OrderGoods, error)
	GetOrderStatus(ctx context.Context, id int64) (int8, int64, error) // GetOrderStatus Return paymentStatus orderedAt error

	UpdateOrderStatus(ctx context.Context, orderID int64, status int32) error
	UpdateOrderAddress(ctx context.Context, orderID int64, addressID int64, addressInfo string) error
	UpdatePaymentStatus(ctx context.Context, message *model.PaymentResult) error

	DeleteOrder(ctx context.Context, orderID int64) error

	IsOrderExist(ctx context.Context, orderID int64) (bool, int64, error)
}

type MQ interface {
	SendSyncMsg(ctx context.Context, topic string, msgs ...*model.MqMessage) error
	SubscribeTopic(ctx context.Context, topic string, pullMsgInterval time.Duration, fn func(ctx context.Context, body []byte) bool) error
	Shutdown() []error
}

type RPC interface {
	GetAddressInfo(ctx context.Context, addressId int64) (string, error)
	QueryGoodsInfo(ctx context.Context, goods []*model.BaseOrderGoods) ([]*model.OrderGoods, error)
	WithholdSkuStock(ctx context.Context, stocks *model.OrderStock) error
	RollbackSkuStock(ctx context.Context, stocks *model.OrderStock) error
	DescSkuStock(ctx context.Context, stocks *model.OrderStock) error
}

type Cache interface {
	SetPaymentStatus(ctx context.Context, s *model.CachePaymentStatus) error
	GetPaymentStatus(ctx context.Context, orderID int64) (*model.CachePaymentStatus, bool, error)
	UpdatePaymentStatus(ctx context.Context, s *model.CachePaymentStatus) (exist bool, err error)
	DeletePaymentStatus(ctx context.Context, orderID int64) error
}

type Locker interface {
	LockOrder(orderID int64) error
	UnlockOrder(orderID int64) error
}

type IDGenerator interface {
	NextVal() (int64, error)
}
