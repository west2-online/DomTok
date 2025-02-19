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

// CreateOrder 创建订单，包含业务逻辑
func (svc *OrderService) CreateOrder(ctx context.Context, order *model.Order, goods []*model.OrderGoods) (int64, error) {
	orderID := svc.nextID()
	order.Id = orderID

	if err := svc.db.CreateOrder(ctx, order); err != nil {
		return 0, err
	}

	// todo: 未实现

	if err := svc.db.CreateOrderGoods(ctx, goods); err != nil {
		return 0, err
	}

	return orderID, nil
}

// CancelOrder 取消订单，包含状态检查
func (svc *OrderService) CancelOrder(ctx context.Context, orderID int64) error {
	// 1. 检查订单是否存在
	exist, err := svc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 获取订单信息检查状态
	order, err := svc.db.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 3. 只有待支付的订单可以取消
	if order.Status != constants.OrderStatusUnpaidCode {
		return errno.NewErrNo(errno.ServiceError, "order cannot be canceled")
	}

	// 4. 更新订单状态为已取消
	return svc.db.UpdateOrderStatus(ctx, orderID, constants.OrderStatusCancelledCode)
}

// UpdateOrderAddress 更新订单地址
func (svc *OrderService) UpdateOrderAddress(ctx context.Context, orderID int64, addressID int64, addressInfo string) error {
	// 1. 检查订单是否存在
	exist, err := svc.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 获取订单信息检查状态
	order, err := svc.db.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 3. 已完成或已取消的订单不能修改地址
	if order.Status == constants.OrderStatusCompletedCode ||
		order.Status == constants.OrderStatusCancelledCode {
		return errno.NewErrNo(errno.ServiceError, "completed or canceled order cannot change address")
	}

	// 4. 更新地址
	return svc.db.UpdateOrderAddress(ctx, orderID, addressID, addressInfo)
}

func (svc *OrderService) nextID() int64 {
	id, _ := svc.sf.NextVal()
	return id
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

// ViewOrder 获取订单详情
func (svc *OrderService) ViewOrder(ctx context.Context, orderID int64) (*model.Order, []*model.OrderGoods, error) {
	// 1. 检查订单是否存在
	exist, err := svc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		return nil, nil, errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 获取订单信息
	order, err := svc.db.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, nil, err
	}

	// 3. 获取订单商品信息
	goods, err := svc.db.GetOrderGoodsByOrderID(ctx, orderID)
	if err != nil {
		return nil, nil, err
	}

	return order, goods, nil
}

// ChangeDeliverAddress 更改配送地址
func (svc *OrderService) ChangeDeliverAddress(ctx context.Context, orderID, addressID int64, addressInfo string) error {
	// 1. 检查订单是否存在
	exist, err := svc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 获取订单信息检查状态
	order, err := svc.db.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 3. 已完成/取消的订单不能修改地址
	if order.Status >= constants.OrderStatusCompletedCode {
		return errno.NewErrNo(errno.ServiceError, "order cannot change address")
	}

	// 4. 更新地址信息
	return svc.db.UpdateOrderAddress(ctx, orderID, addressID, addressInfo)
}

// DeleteOrder 删除订单
func (svc *OrderService) DeleteOrder(ctx context.Context, orderID int64) error {
	// 1. 检查订单是否存在
	exist, err := svc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 获取订单信息检查状态
	order, err := svc.db.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 3. 已支付的订单不能删除
	if order.Status == constants.OrderStatusPaidCode {
		return errno.NewErrNo(errno.ServiceError, "order cannot be deleted")
	}

	// 4. 删除订单（包含订单商品）
	return svc.db.DeleteOrder(ctx, orderID)
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
