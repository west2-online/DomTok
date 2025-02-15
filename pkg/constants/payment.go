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
	PaymentOrderNotExist               = -1
	PaymentStatusPending               = 0 // 待支付
	PaymentStatusProcessing            = 1 // 处理中
	PaymentStatusSuccess               = 2 // 成功支付
	PaymentStatusFailed                = 3 // 支付失败
	PaymentOrderNotExistToken          = ""
	PaymentOrderNotExistExpirationTime = 0
	UserNotExist                       = -1
	UserNotExistToken                  = ""
	UserNotExistExpirationTime         = 0
	HavePaidToken                      = ""
	HavePaidExpirationTime             = 0
	ErrorToken                         = ""
	ErrorExpirationTime                = 0
	PaymentSecretKey                   = "west2online"
	RedisStoreSuccess                  = 0  // 成功
	RedisStoreFailed                   = -1 // Redis 存储失败
)
