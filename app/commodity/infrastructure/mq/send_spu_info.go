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
	"strconv"

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/kafka"
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
		V: []byte(strconv.FormatInt(id, 10)),
	}

	err := c.Send(ctx, constants.KafkaDeleteSpuTopic, []*kafka.Message{msg})
	if err != nil {
		return errno.Errorf(errno.InternalKafkaErrorCode, "CommodityMQ.SendSpuInfo failed: %v", err)
	}
	return nil
}
