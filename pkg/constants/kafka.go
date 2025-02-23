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
	KafkaReadMinBytes      = 512 * B
	KafkaReadMaxBytes      = 1 * MB
	KafkaRetries           = 3
	DefaultReaderGroupID   = "r"
	DefaultTimeRetainHours = 6 // 6小时

	DefaultConsumerChanCap         = 20
	DefaultKafkaProductorSyncWrite = false

	DefaultKafkaNumPartitions     = -1
	DefaultKafkaReplicationFactor = -1
)

// Commodity Service
const (
	KafkaImageTopic            = "ImageTopic"
	KafkaCreateSpuTopic        = "SpuCreateTopic"
	KafkaUpdateSpuTopic        = "SpuUpdateTopic"
	KafkaDeleteSpuTopic        = "SpuDeleteTopic"
	KafkaPartitionNum          = 3
	KafkaImageGroupId          = "ImageGroupId"
	KafkaCreateSpuGroupId      = "CreateSpuGroupId"
	KafkaUpdateSpuGroupId      = "UpdateSpuGroupId"
	KafkaDeleteSpuGroupId      = "DeleteSpuGroupId"
	KafkaCommodityCreateSpuNum = 3
	KafkaCommodityUpdateSpuNum = 3
	KafkaCommodityDeleteSpuNum = 3
)

// CartService
const (
	KafkaCartTopic                = "cart"           // Kafka的话题
	KafkaCartAddGoodsPartitionNum = 10               // Kafka的分区数
	KafkaCartAddGoodsConsumerNum  = 10               // Kafka的并发消费者数
	KafkaCartAddGoodsGroupId      = "cart_add_goods" // Kafka的订阅组id
)
