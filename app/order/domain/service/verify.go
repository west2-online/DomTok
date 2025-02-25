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
	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

type OrderVerifyOps func() error

// Verify 验证多个条件
func (svc *OrderService) Verify(ops ...OrderVerifyOps) error {
	for _, opt := range ops {
		if err := opt(); err != nil {
			return err
		}
	}
	return nil
}

// VerifyOrderStatus 验证订单状态
func (svc *OrderService) VerifyOrderStatus(status int32) func() error {
	return func() error {
		if svc.GetOrderStatusMsg(int8(status)) == constants.OrderStatusUnknown {
			return errno.NewErrNo(errno.UnknownOrderStatus, "invalid order status")
		}
		return nil
	}
}

func (svc *OrderService) VerifyAddressID(id int64) func() error {
	return func() error {
		if id <= 0 {
			return errno.NewErrNo(errno.ParamVerifyErrorCode, "address id should lest than 0")
		}
		return nil
	}
}

func (svc *OrderService) VerifyBaseOrderGoods(baseGoods []*model.BaseOrderGoods) func() error {
	return func() error {
		if baseGoods == nil || len(baseGoods) <= 0 {
			return errno.NewErrNo(errno.ParamVerifyErrorCode, "create order needed at least one goods")
		}
		return nil
	}
}
