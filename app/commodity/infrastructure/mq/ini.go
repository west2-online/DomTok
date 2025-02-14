package mq

import (
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/pkg/kafka"
)

type CommodityMQ struct {
	client *kafka.Kafka
}

func NewCommodityMQ(client *kafka.Kafka) repository.CommodityMQ {
	return &CommodityMQ{client: client}
}
