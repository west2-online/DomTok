package mq

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/kafka"
	"strconv"
)

func (c *CommodityMQ) SendCreateSpuInfo(ctx context.Context, spu *model.Spu) error {
	v, err := sonic.Marshal(spu)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "CommodityMQ.SendSpuInfo failed: %v", err)
	}

	msg := &kafka.Message{
		K: []byte(strconv.FormatInt(spu.SpuId%constants.KafkaCommodityCreateSpuNum, 10)),
		V: v,
	}
	err = c.Send(ctx, constants.KafkaCreateSpuTopic, []*kafka.Message{msg})
	if err != nil {
		return errno.Errorf(errno.InternalKafkaErrorCode, "CommodityMQ.SendSpuInfo failed: %v", err)
	}
	return nil
}

func (c *CommodityMQ) SendUpdateSpuInfo(ctx context.Context, spu *model.Spu) error {
	v, err := sonic.Marshal(spu)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "CommodityMQ.SendSpuInfo failed: %v", err)
	}

	msg := &kafka.Message{
		K: []byte(strconv.FormatInt(spu.SpuId%constants.KafkaCommodityUpdateSpuNum, 10)),
		V: v,
	}
	err = c.Send(ctx, constants.KafkaUpdateSpuTopic, []*kafka.Message{msg})
	if err != nil {
		return errno.Errorf(errno.InternalKafkaErrorCode, "CommodityMQ.SendSpuInfo failed: %v", err)
	}
	return nil
}

func (c *CommodityMQ) SendDeleteSpuInfo(ctx context.Context, id int64) error {

	msg := &kafka.Message{
		K: []byte(strconv.FormatInt(id%constants.KafkaCommodityDeleteSpuNum, 10)),
		V: []byte(strconv.FormatInt(int64(id), 10)),
	}

	err := c.Send(ctx, constants.KafkaDeleteSpuTopic, []*kafka.Message{msg})
	if err != nil {
		return errno.Errorf(errno.InternalKafkaErrorCode, "CommodityMQ.SendSpuInfo failed: %v", err)
	}
	return nil
}
