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

// Service Name
const (
	GatewayServiceName   = "gateway"
	OrderServiceName     = "order"
	UserServiceName      = "user"
	PaymentServiceName   = "payment"
	CommodityServiceName = "commodity"
	CartServiceName      = "cart"
)

// UserService
const (
	UserMaximumPasswordLength      = 72 // DO NOT EDIT (ref: bcrypt.GenerateFromPassword)
	UserMinimumPasswordLength      = 5
	UserDefaultEncryptPasswordCost = 10
	UserTestId                     = 1
)

// OrderService
const (
	OrderStatusUnpaidCode    = -1
	OrderStatusPaidCode      = 1
	OrderStatusCompletedCode = 2
	OrderStatusCancelledCode = 3
)

// OrderService Status Messages
const (
	OrderStatusUnpaid    = "待支付"
	OrderStatusPaid      = "已支付"
	OrderStatusCompleted = "已完成"
	OrderStatusCancelled = "已取消"
	OrderStatusUnknown   = "未知状态"
)

// CommodityService
const (
	CommodityAllowedForSale    = 1
	CommodityNotAllowedForSale = 2
)
