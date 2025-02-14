package mq

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
	"strconv"
)

func (c CommodityMQ) SendSaveImage(ctx context.Context, image *model.Image) error {

	data, err := sonic.Marshal(image)
	if err != nil {
		return fmt.Errorf("mq.SendSaveImage: marshal err: %v", err)
	}

	msg := []*kafka.Message{
		{
			K: []byte(strconv.FormatInt(image.Id%constants.KafkaPartitionNum, 10)),
			V: data,
		},
	}
	err = c.Send(ctx, constants.KafkaImageTopic, msg)
	if err != nil {
		return fmt.Errorf("mq.SendSaveImage: send msg err: %v", err)
	}
	return nil
}
