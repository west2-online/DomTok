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

package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
)

func (svc *CartService) SendAddGoods(ctx context.Context, uid int64, goods *model.GoodInfo) error {
	msgValue := &model.AddGoodsMsg{
		Uid:   uid,
		Goods: goods,
	}
	v, err := sonic.Marshal(msgValue)
	if err != nil {
		return fmt.Errorf("service.Send: marshal msg failed, err: %w", err)
	}
	msg := []*kafka.Message{
		{
			// 用%来简陋实现一下分区
			K: []byte(strconv.FormatInt(uid%constants.KafkaCartAddGoodsPartitionNum, 10)),
			V: v,
		},
	}

	for i := 0; i < constants.KafkaRetries; i++ {
		errs := svc.MQ.Send(ctx, constants.KafkaCartTopic, msg)
		if len(errs) == 0 {
			return nil
		} else {
			var errMsg string
			for _, e := range errs {
				errMsg = strings.Join([]string{errMsg, e.Error(), ";"}, "")
			}
			err = fmt.Errorf("service.Send: send msg failed, errs: %v", errMsg)
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}

	return err
}
