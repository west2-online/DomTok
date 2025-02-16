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
	GatewayServiceName = "gateway"
	OrderServiceName   = "order"
	UserServiceName    = "user"
)

// UserService
const (
	UserMaximumPasswordLength      = 72 // DO NOT EDIT (ref: bcrypt.GenerateFromPassword)
	UserMinimumPasswordLength      = 5
	UserDefaultEncryptPasswordCost = 10
)

// OrderService
const (
	OrderStatusUnpaidCode    int32 = 0
	OrderStatusPaidCode      int32 = 1
	OrderStatusCompletedCode int32 = 2
	OrderStatusCancelledCode int32 = 3
)

// OrderService Status Messages
const (
	OrderStatusUnpaid    = "待支付"
	OrderStatusPaid      = "已支付"
	OrderStatusCompleted = "已完成"
	OrderStatusCancelled = "已取消"
	OrderStatusUnknown   = "未知状态"
)

func GetOrderStatusMsg(code int32) string {
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
