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
	"fmt"
	"strings"

	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
)

func (c *KafkaAdapter) send(ctx context.Context, msg []*kafka.Message) (err error) {
	if !c.done.Load() {
		err = c.mq.SetWriter(constants.KafkaCartTopic, true)
		if err != nil {
			return err
		}
		c.done.Swap(true)
	}
	errs := c.mq.Send(ctx, constants.KafkaCartTopic, msg)
	if len(errs) != 0 {
		var errMsg string
		for _, e := range errs {
			errMsg = strings.Join([]string{errMsg, e.Error(), ";"}, "")
		}
		err = fmt.Errorf("mq.Send: send msg failed, errs: %v", errMsg)
		return err
	}
	return nil
}
