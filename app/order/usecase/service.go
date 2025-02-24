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
	basecontext "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

// ViewOrderList 获取订单列表
func (uc *useCase) ViewOrderList(ctx context.Context, page, size int32) ([]*model.Order, []*model.OrderGoods, int32, error) {
	// 从 RPC 上下文中获取用户ID
	userID, err := basecontext.GetLoginData(ctx)
	if err != nil {
		return nil, nil, 0, errno.NewErrNo(errno.AuthInvalidCode, "invalid user id")
	}

	return uc.svc.ViewOrderList(ctx, userID, page, size)
}

// ViewOrder 获取订单详情
func (uc *useCase) ViewOrder(ctx context.Context, orderID int64) (*model.Order, []*model.OrderGoods, error) {
	if err := uc.svc.OrderExist(ctx, orderID); err != nil {
		return nil, nil, err
	}

	order, orderGoods, err := uc.db.GetOrderWithGoods(ctx, orderID)
	if err != nil {
		return nil, nil, err
	}

	return order, orderGoods, nil
}

// CancelOrder 取消订单
func (uc *useCase) CancelOrder(ctx context.Context, orderID int64) error {
	// 1. 检查订单是否存在
	exist, err := uc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 获取订单信息检查状态
	order, err := uc.db.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 3. 只有待支付的订单可以取消
	if order.Status != constants.OrderStatusUnpaidCode {
		return errno.NewErrNo(errno.ServiceError, "order cannot be canceled")
	}

	// 4. 更新订单状态为已取消
	return uc.db.UpdateOrderStatus(ctx, orderID, constants.OrderStatusCancelledCode)
}

// ChangeDeliverAddress 更改配送地址
func (uc *useCase) ChangeDeliverAddress(ctx context.Context, orderID, addressID int64, addressInfo string) error {
	// 1. 检查订单是否存在
	exist, err := uc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 获取订单信息检查状态
	order, err := uc.db.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 3. 已完成/取消的订单不能修改地址
	if order.Status >= constants.OrderStatusCompletedCode {
		return errno.NewErrNo(errno.ServiceError, "order cannot change address")
	}

	// 4. 更新地址信息
	return uc.db.UpdateOrderAddress(ctx, orderID, addressID, addressInfo)
}

// DeleteOrder 删除订单
func (uc *useCase) DeleteOrder(ctx context.Context, orderID int64) error {
	// 1. 检查订单是否存在
	exist, err := uc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 删除订单（包含订单商品）
	return uc.db.DeleteOrder(ctx, orderID)
}

func (uc *useCase) IsOrderExist(ctx context.Context, orderID int64) (bool, error) {
	return uc.db.IsOrderExist(ctx, orderID)
}
