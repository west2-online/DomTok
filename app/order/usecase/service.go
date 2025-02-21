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

	"github.com/west2-online/DomTok/app/order/domain/model"
	basecontext "github.com/west2-online/DomTok/pkg/base/context"
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
	if err := uc.svc.CancelOrder(ctx, orderID); err != nil {
		return fmt.Errorf("cancel order failed: %w", err)
	}
	return nil
}

// ChangeDeliverAddress 更改配送地址
func (uc *useCase) ChangeDeliverAddress(ctx context.Context, orderID, addressID int64, addressInfo string) error {
	if err := uc.svc.ChangeDeliverAddress(ctx, orderID, addressID, addressInfo); err != nil {
		return fmt.Errorf("change deliver address failed: %w", err)
	}
	return nil
}

// DeleteOrder 删除订单
func (uc *useCase) DeleteOrder(ctx context.Context, orderID int64) error {
	if err := uc.svc.DeleteOrder(ctx, orderID); err != nil {
		return fmt.Errorf("delete order failed: %w", err)
	}
	return nil
}
