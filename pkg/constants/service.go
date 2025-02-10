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
	OrderStatusUnpaidCode = 0
	OrderStatusPaidCode   = 1
	OrderStatusFailCode   = 2

	OrderStatusUnpaid  = "未支付"
	OrderStatusPaid    = "待支付"
	OrderStatusFail    = "支付失败"
	OrderStatusUnknown = "未知状态"
)
