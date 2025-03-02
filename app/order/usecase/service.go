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

func (uc *useCase) CreateOrder(ctx context.Context, addressID int64, baseGoods []*model.BaseOrderGoods) (int64, error) {
	if err := uc.svc.Verify(uc.svc.VerifyAddressID(addressID), uc.svc.VerifyBaseOrderGoods(baseGoods)); err != nil {
		return 0, err
	}

	addressInfo, err := uc.rpc.GetAddressInfo(ctx, addressID)
	if err != nil {
		return 0, err
	}

	goods, err := uc.rpc.QueryGoodsInfo(ctx, baseGoods)
	if err != nil {
		return 0, err
	}

	order, err := uc.svc.MakeOrderByGoods(ctx, addressID, addressInfo, goods)
	if err != nil {
		return 0, err
	}

	if err = uc.svc.CreateOrder(ctx, order, goods); err != nil {
		return 0, err
	}

	if err = uc.svc.WithholdSkuStock(ctx, order.Id, goods); err != nil {
		return 0, err
	}

	return order.Id, nil
}

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

	order, orderGoods, err := uc.db.GetOrderAndGoods(ctx, orderID)
	if err != nil {
		return nil, nil, err
	}

	return order, orderGoods, nil
}

// CancelOrder 取消订单
func (uc *useCase) CancelOrder(ctx context.Context, orderID int64) error {
	// 1. 检查订单是否存在
	exist, _, err := uc.db.IsOrderExist(ctx, orderID)
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
		return errno.NewErrNo(errno.IllegalOperatorCode, "order cannot be canceled")
	}

	// 4. 更新订单状态为已取消
	return uc.db.UpdateOrderStatus(ctx, orderID, constants.OrderStatusCancelledCode)
}

// ChangeDeliverAddress 更改配送地址
func (uc *useCase) ChangeDeliverAddress(ctx context.Context, orderID, addressID int64, addressInfo string) error {
	// 1. 检查订单是否存在
	exist, _, err := uc.db.IsOrderExist(ctx, orderID)
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
		return errno.NewErrNo(errno.IllegalOperatorCode, "order cannot change address")
	}

	// 4. 更新地址信息
	return uc.db.UpdateOrderAddress(ctx, orderID, addressID, addressInfo)
}

// DeleteOrder 删除订单
func (uc *useCase) DeleteOrder(ctx context.Context, orderID int64) error {
	// 1. 检查订单是否存在
	exist, _, err := uc.db.IsOrderExist(ctx, orderID)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
	}

	// 2. 删除订单（包含订单商品）
	return uc.db.DeleteOrder(ctx, orderID)
}

func (uc *useCase) IsOrderExist(ctx context.Context, orderID int64) (bool, int64, error) {
	return uc.db.IsOrderExist(ctx, orderID)
}

func (uc *useCase) OrderPaymentSuccess(ctx context.Context, req *model.PaymentResult) error {
	// 这里不进行 orderID 是否存在的检查，因为这个方法是由 payment 服务调用的，payment 在调用之前已经检查了 orderID 的存在
	status, expired, err := uc.svc.GetPaymentStatusAndOrderExpire(ctx, req.OrderID)
	if err != nil {
		return err
	}

	if uc.svc.IsEqualStatus(status, req.PaymentStatus) {
		return nil
	}

	if err = uc.svc.UpdateOrderAsSuccess(ctx, expired, req); err != nil {
		return err
	}
	return nil
}

func (uc *useCase) OrderPaymentCancel(ctx context.Context, req *model.PaymentResult) error {
	// 这里不进行 orderID 是否存在的检查，因为这个方法是由 payment 服务调用的，payment 在调用之前已经检查了 orderID 的存在
	status, _, err := uc.svc.GetPaymentStatusAndOrderExpire(ctx, req.OrderID)
	if err != nil {
		return err
	}

	if uc.svc.IsEqualStatus(status, req.PaymentStatus) {
		return nil
	}

	if err = uc.svc.CancelOrder(ctx, req); err != nil {
		return err
	}
	return nil
}
