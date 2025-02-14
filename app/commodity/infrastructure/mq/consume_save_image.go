package mq

import (
	"context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
)

func (c CommodityMQ) ConsumeSaveImage(ctx context.Context) <-chan *kafka.Message {
	msgChan := c.client.Consume(ctx,
		constants.KafkaImageTopic,
		constants.KafkaPartitionNum,
		constants.KafkaImageGroupId,
		constants.DefaultConsumerChanCap,
	)
	return msgChan
}
