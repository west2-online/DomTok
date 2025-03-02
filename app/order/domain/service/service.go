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
	"github.com/west2-online/DomTok/pkg/logger"
	"time"
)

type OrderService struct {
	db     repository.OrderDB
	idG    repository.IDGenerator
	rpc    repository.RPC
	mq     repository.MQ
	cache  repository.Cache
	locker repository.Locker
}

func NewOrderService(db repository.OrderDB, sf repository.IDGenerator, rpc repository.RPC,
	mq repository.MQ, cache repository.Cache, locker repository.Locker,
) *OrderService {
	if db == nil || sf == nil || rpc == nil || mq == nil || cache == nil || locker == nil {
		logger.Fatalf("failed get new order service, all arguments should not be nil")
	}
	svc := &OrderService{db: db, idG: sf, rpc: rpc, mq: mq, cache: cache, locker: locker}
	svc.init()
	return svc
}

// IsOrderExist 检查订单是否存在
func (svc *OrderService) IsOrderExist(ctx context.Context, orderID int64) (bool, int64, error) {
	return svc.db.IsOrderExist(ctx, orderID)
}

// OrderExist 检查订单是否存在
func (svc *OrderService) OrderExist(ctx context.Context, orderID int64) error {
	exist, _, err := svc.db.IsOrderExist(ctx, orderID)
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
	return constants.GetOrderStatusMsg(code)
}

func (svc *OrderService) calcOrderExpireTime(createAt int64) int64 {
	return createAt + constants.OrderExpireTime.Milliseconds()
}

func (svc *OrderService) UpdateOrderAsSuccess(ctx context.Context, expired int64, payRel *model.PaymentResult) error {
	// 对状态的一层校验
	if payRel.PaymentStatus != constants.PaymentStatusSuccessCode {
		return errno.NewErrNo(errno.ServiceOrderStatusInvalid, "success orders cannot be updated again")
	}

	if time.Now().UnixMilli() > expired {
		return errno.NewErrNo(errno.ServiceOrderExpired, "order expired")
	}
	// 尝试开始更新
	if err := svc.locker.LockOrder(payRel.OrderID); err != nil {
		return err
	}
	defer svc.logUnLock(payRel.OrderID, svc.locker.UnlockOrder)

	goods, err := svc.db.GetOrderGoodsByOrderID(ctx, payRel.OrderID)
	if err != nil {
		return err
	}

	// 减少库存
	orderStock := model.ConvertOrderGoodsToOrderStock(payRel.OrderID, goods)
	if err = svc.rpc.DescSkuStock(ctx, orderStock); err != nil {
		return err
	}

	if err = svc.db.UpdatePaymentStatus(ctx, payRel); err != nil {
		return err
	}

	if _, err = svc.cache.UpdatePaymentStatus(ctx, &model.CachePaymentStatus{}); err != nil {
		return err
	}

	return nil
}

func (svc *OrderService) CancelOrder(ctx context.Context, payRel *model.PaymentResult) error {
	// 如果订单是失败，那说明回滚过了
	if payRel.PaymentStatus != constants.PaymentStatusFailedCode {
		return errno.NewErrNo(errno.ServiceOrderStatusInvalid, "failed orders cannot be canceled")
	}

	// 尝试开始回滚或者直接取消回滚
	if err := svc.locker.LockOrder(payRel.OrderID); err != nil {
		return err
	}
	defer svc.logUnLock(payRel.OrderID, svc.locker.UnlockOrder)

	goods, err := svc.db.GetOrderGoodsByOrderID(ctx, payRel.OrderID)
	if err != nil {
		return err
	}

	// 释放锁定库存
	orderStock := model.ConvertOrderGoodsToOrderStock(payRel.OrderID, goods)
	if err = svc.rpc.RollbackSkuStock(ctx, orderStock); err != nil {
		return err
	}

	if err = svc.db.UpdatePaymentStatus(ctx, payRel); err != nil {
		return err
	}

	if _, err = svc.cache.UpdatePaymentStatus(ctx, &model.CachePaymentStatus{}); err != nil {
		return err
	}

	return nil
}

func (svc *OrderService) IsEqualStatus(s1, s2 int8) bool {
	return s1 == s2
}

func (svc *OrderService) nextVal() int64 {
	v, _ := svc.idG.NextVal()
	return v
}

func (svc *OrderService) logUnLock(id int64, fn func(id int64) error) {
	if err := fn(id); err != nil {
		logger.Error(err.Error())
	}
}
