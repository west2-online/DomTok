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

	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
)

// SkuLockStockRollback 只负责进行回滚，如果订单被完成那应该由 payment 的 msg 或者直接调用 rpc 来对订单进行更改
func (svc *OrderService) SkuLockStockRollback(ctx context.Context, body []byte) (sucRollback bool) {
	orderStock, err := svc.decodeStocks(body)
	if err != nil {
		logger.Errorf(err.Error())
		return true
	}

	paymentStatus, exist, err := svc.cache.GetPaymentStatus(ctx, orderStock.OrderID)
	if err != nil {
		logger.Errorf(err.Error())
	}

	if !exist {
		var orderedAt int64
		paymentStatus.PaymentStatus, orderedAt, err = svc.db.GetOrderStatus(ctx, orderStock.OrderID)
		if err != nil {
			logger.Errorf(err.Error())
			return true
		}
		paymentStatus.OrderExpire = svc.GetOrderExpireTime(orderedAt)
	}

	if paymentStatus.PaymentStatus == constants.PaymentStatusPendingCode {
		if err = svc.rpc.IncrSkuLockStock(ctx, orderStock); err != nil {
			logger.Errorf(err.Error())

			if e := svc.cache.SetPaymentStatus(ctx, paymentStatus); e != nil {
				logger.Errorf(e.Error())
			}
			return false
		}
	}

	return true
}
