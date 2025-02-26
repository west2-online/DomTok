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

// Order listen`s topic
const (
	SkuStockRollbackTopic                     = "sku-stock-rollback-topic"
	SkuStockRollbackTopicDelayTimeLevel       = RocketMQDelay10M // 10 Minute
	SkuStockRollbackTopicConsumerPullInterval = 1 * time.Second

	OrderPaymentResultTopic                     = "order-payment-result-topic"
	OrderPaymentResultTopicConsumerPullInterval = 1 * time.Second
)

const (
	RocketMQDelay1S = iota + 1
	RocketMQDelay5S
	RocketMQDelay10S
	RocketMQDelay30S
	RocketMQDelay1M
	RocketMQDelay2M
	RocketMQDelay3M
	RocketMQDelay4M
	RocketMQDelay5M
	RocketMQDelay6M
	RocketMQDelay7M
	RocketMQDelay8M
	RocketMQDelay9M
	RocketMQDelay10M
	RocketMQDelay20M
	RocketMQDelay30M
	RocketMQDelay1H
	RocketMQDelay2H
)
