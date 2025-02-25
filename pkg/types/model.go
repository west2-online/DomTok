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

package types

// PaymentResultMessage 是 payment 模块与 order 模块之间的 msg 格式约定, 请勿修改
// DO NOT EDIT
type PaymentResultMessage struct {
	OrderID       int64  // 订单 id
	PaymentStatus int    // 状态
	PaymentAt     int64  // 支付时间, 毫秒级时间戳
	PaymentStyle  string // 支付类型
}
