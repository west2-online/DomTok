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
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
)

// send 对内部库的send添加了一层重试
func (c *KafkaAdapter) send(ctx context.Context, msg []*kafka.Message) (err error) {
	// 这里参数没有动的必要，直接设为固定，实际也可以改为调用时传入
	for i := 0; i < constants.KafkaRetries; i++ {
		errs := c.mq.Send(ctx, constants.KafkaCartTopic, msg)
		if len(errs) == 0 {
			return nil
		} else {
			var errMsg string
			for _, e := range errs {
				errMsg = strings.Join([]string{errMsg, e.Error(), ";"}, "")
			}
			err = fmt.Errorf("mq.Send: send msg failed, errs: %v", errMsg)
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return err
}
