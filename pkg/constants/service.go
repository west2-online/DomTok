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

package constants

import "time"

// Service Name
const (
	GatewayServiceName   = "gateway"
	OrderServiceName     = "order"
	UserServiceName      = "user"
	PaymentServiceName   = "payment"
	CommodityServiceName = "commodity"
	CartServiceName      = "cart"
	AssistantServiceName = "assistant"
)

// UserService
const (
	UserMaximumPasswordLength      = 72 // DO NOT EDIT (ref: bcrypt.GenerateFromPassword)
	UserMinimumPasswordLength      = 5
	UserDefaultEncryptPasswordCost = 10
	UserTestId                     = 1
	UserTestAddr                   = 1
)

// OrderService
const (
	OrderStatusUnpaidCode    = -1
	OrderStatusPaidCode      = 1
	OrderStatusCompletedCode = 2
	OrderStatusCancelledCode = 3

	OrderExpireTime      = 10 * time.Minute
	OrderStatusUnpaid    = "待支付"
	OrderStatusPaid      = "已支付"
	OrderStatusCompleted = "已完成" // 已发货已签收
	OrderStatusCancelled = "已取消"
	OrderStatusUnknown   = "未知状态"

	OrderMqConsumerGroupFormat = "order-%s" // order-topic
)

// GetOrderStatusMsg logic should be equal with /app/order/domain/service/service.go
func GetOrderStatusMsg(code int8) string {
	switch code {
	case OrderStatusUnpaidCode:
		return OrderStatusUnpaid
	case OrderStatusPaidCode:
		return OrderStatusPaid
	case OrderStatusCompletedCode:
		return OrderStatusCompleted
	case OrderStatusCancelledCode:
		return OrderStatusCancelled
	default:
		return OrderStatusUnknown
	}
}

// CommodityService
const (
	CommodityAllowedForSale    = 1
	CommodityNotAllowedForSale = 2

	CommodityDefaultMinCost = 0.0
	CommodityDefaultMaxCost = 1e7
	// CommodityMaxBuyNum 指定了最大商品购买数
	CommodityMaxBuyNum = 1000
)
