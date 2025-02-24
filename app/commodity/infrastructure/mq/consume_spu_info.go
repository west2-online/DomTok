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

package mq

import (
	"context"

	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
)

func (c *CommodityMQ) ConsumeCreateSpuInfo(ctx context.Context) <-chan *kafka.Message {
	return c.client.Consume(ctx, constants.KafkaCreateSpuTopic, constants.KafkaCommodityCreateSpuNum,
		constants.KafkaCreateSpuGroupId, constants.KafkaESConsumerChanCap)
}

func (c *CommodityMQ) ConsumeUpdateSpuInfo(ctx context.Context) <-chan *kafka.Message {
	return c.client.Consume(ctx, constants.KafkaUpdateSpuTopic, constants.KafkaCommodityUpdateSpuNum,
		constants.KafkaUpdateSpuGroupId, constants.KafkaESConsumerChanCap)
}

func (c *CommodityMQ) ConsumeDeleteSpuInfo(ctx context.Context) <-chan *kafka.Message {
	return c.client.Consume(ctx, constants.KafkaDeleteSpuTopic, constants.KafkaCommodityDeleteSpuNum,
		constants.KafkaDeleteSpuGroupId, constants.KafkaESConsumerChanCap)
}
