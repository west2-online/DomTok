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
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
)

// SkuLockStockRollback 只负责进行回滚，如果订单被完成那应该由 payment 的 msg 或者直接调用 rpc 来对订单进行更改
func (svc *OrderService) SkuLockStockRollback(ctx context.Context, body []byte) (sucRollback bool) {
	orderStock, err := svc.decodeStocks(body)
	if err != nil {
		logger.Error(err.Error())
		return true
	}

	var payRel model.CachePaymentStatus
	if payRel.PaymentStatus, payRel.OrderExpire, err = svc.GetPaymentStatusAndOrderExpire(ctx, orderStock.OrderID); err != nil {
		logger.Error(err.Error())
		return false
	}

	if payRel.PaymentStatus == constants.PaymentStatusPendingCode {
		if err = svc.rpc.RollbackSkuStock(ctx, orderStock); err != nil {
			logger.Error(err.Error())
			return false
		}
	}

	return true
}

func (svc *OrderService) GetPaymentStatusAndOrderExpire(ctx context.Context, orderID int64) (status int8, expired int64, err error) {
	paymentStatus, exist, err := svc.cache.GetPaymentStatus(ctx, orderID)
	if err != nil {
		logger.Error(err.Error())
	}

	if exist {
		return paymentStatus.PaymentStatus, paymentStatus.OrderExpire, nil
	}

	// 获取 Status
	var orderedAt int64
	paymentStatus.PaymentStatus, orderedAt, err = svc.db.GetOrderStatus(ctx, orderID)
	if err != nil {
		return 0, 0, err
	}
	paymentStatus.OrderExpire = svc.calcOrderExpireTime(orderedAt)

	if e := svc.cache.SetPaymentStatus(ctx, paymentStatus); e != nil {
		logger.Error(e.Error())
	}

	return paymentStatus.PaymentStatus, paymentStatus.OrderExpire, nil
}
