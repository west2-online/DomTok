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

package service

import (
	"context"

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/utils"
)

type OrderService struct {
	db repository.OrderDB
	sf *utils.Snowflake
}

func NewOrderService(db repository.OrderDB, sf *utils.Snowflake) *OrderService {
	return &OrderService{db: db, sf: sf}
}

// IsOrderExist 检查订单是否存在
func (svc *OrderService) IsOrderExist(ctx context.Context, orderID int64) (bool, error) {
	return svc.db.IsOrderExist(ctx, orderID)
}

// OrderExist 检查订单是否存在
func (svc *OrderService) OrderExist(ctx context.Context, orderID int64) error {
	exist, err := svc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not exist")
	}
	return nil
}

// ViewOrderList 获取订单列表
func (svc *OrderService) ViewOrderList(ctx context.Context, userID int64, page, size int32) ([]*model.Order, []*model.OrderGoods, int32, error) {
	// 1. 获取订单列表
	orders, total, err := svc.db.GetOrdersByUserID(ctx, userID, page, size)
	if err != nil {
		return nil, nil, 0, err
	}

	// 2. 如果没有订单，直接返回
	if len(orders) == 0 {
		return nil, nil, total, nil
	}

	// 3. 获取所有订单的商品信息
	var allGoods []*model.OrderGoods
	for _, order := range orders {
		goods, err := svc.db.GetOrderGoodsByOrderID(ctx, order.Id)
		if err != nil {
			return nil, nil, 0, err
		}
		allGoods = append(allGoods, goods...)
	}

	return orders, allGoods, total, nil
}

func (svc *OrderService) GetOrderStatusMsg(code int8) string {
	switch code {
	case constants.OrderStatusUnpaidCode:
		return constants.OrderStatusUnpaid
	case constants.OrderStatusPaidCode:
		return constants.OrderStatusPaid
	case constants.OrderStatusCompletedCode:
		return constants.OrderStatusCompleted
	case constants.OrderStatusCancelledCode:
		return constants.OrderStatusCancelled
	default:
		return constants.OrderStatusUnknown
	}
}
