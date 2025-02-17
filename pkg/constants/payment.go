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

const (
	OrderNotExist     = false
	PaymentExist      = true
	PaymentNotExist   = false
	UserNotExist      = false
	PaymentSecretKey  = "west2online"
	RedisStoreSuccess = true  // 成功
	RedisStoreFailed  = false // Redis 存储失败
	// TODO 这两个常量要变
	ExpirationTime = 15
	PingTime       = 2
)
const (
	PaymentStatusPending    = iota // 待支付
	PaymentStatusProcessing        // 处理中
	PaymentStatusSuccess           // 成功支付
	PaymentStatusFailed            // 支付失败
)
