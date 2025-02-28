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
	"testing"

	"github.com/apache/rocketmq-client-go/v2/rlog"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func initTest(t *testing.T) repository.MQ {
	t.Helper()
	config.Init("test")
	logger.Ignore()
	rlog.SetLogLevel("fatal")
	return NewRocketmq()
}

var (
	testTopic   = "test111"
	testMsgBody = "test"
)

func TestOrder_RocketMqSendAndConsume(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	initTest(t)
	mq := NewRocketmq()
	ctx := context.Background()
	Convey("Test rocketmq send and consume message", t, func() {
		wait := make(chan struct{})
		err := mq.SubscribeTopic(ctx, testTopic, 0, func(ctx context.Context, body []byte) bool {
			close(wait)
			return true
		})
		So(err, ShouldBeNil)
		err = mq.SendSyncMsg(ctx, testTopic, &model.MqMessage{Body: []byte(testMsgBody)})
		So(err, ShouldBeNil)
		<-wait
	})
}
