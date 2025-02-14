package mq

import (
	"context"
	"errors"
	"github.com/west2-online/DomTok/pkg/kafka"
)

func (c CommodityMQ) Send(ctx context.Context, topic string, message []*kafka.Message) error {
	errs := c.client.Send(ctx, topic, message)
	var res error
	if len(errs) > 0 {
		for _, err := range errs {
			if err != nil {
				res = errors.Join(res, err)
			}
		}
	}
	return res
}
