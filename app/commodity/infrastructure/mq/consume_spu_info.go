package mq

import (
	"context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
)

func (c *CommodityMQ) ConsumeCreateSpuInfo(ctx context.Context) <-chan *kafka.Message {
	return c.client.Consume(ctx, constants.KafkaCreateSpuTopic, constants.KafkaCommodityCreateSpuNum,
		constants.KafkaCreateSpuGroupId, constants.DefaultConsumerChanCap)
}

func (c *CommodityMQ) ConsumeUpdateSpuInfo(ctx context.Context) <-chan *kafka.Message {
	return c.client.Consume(ctx, constants.KafkaUpdateSpuTopic, constants.KafkaCommodityUpdateSpuNum,
		constants.KafkaUpdateSpuGroupId, constants.DefaultConsumerChanCap)
}

func (c *CommodityMQ) ConsumeDeleteSpuInfo(ctx context.Context) <-chan *kafka.Message {
	return c.client.Consume(ctx, constants.KafkaDeleteSpuTopic, constants.KafkaCommodityDeleteSpuNum,
		constants.KafkaDeleteSpuGroupId, constants.DefaultConsumerChanCap)
}
