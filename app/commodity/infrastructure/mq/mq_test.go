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

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/kafka"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func initTest(t *testing.T) repository.CommodityMQ {
	t.Helper()
	config.Init("test")
	logger.Ignore()

	return NewCommodityMQ(kafka.NewKafkaInstance())
}

var testDeleteId int64 = 1

func TestCommodityMQ_SendAndConsume(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	initTest(t)

	kfk := NewCommodityMQ(kafka.NewKafkaInstance())
	ctx := context.Background()

	Convey("Test kafka send and consume message", t, func() {
		var err error
		var msg <-chan *kafka.Message

		err = kfk.SendCreateSpuInfo(ctx, &model.Spu{})
		So(err, ShouldBeNil)
		msg = kfk.ConsumeCreateSpuInfo(ctx)
		So(msg, ShouldNotBeNil)

		err = kfk.SendUpdateSpuInfo(ctx, &model.Spu{})
		So(err, ShouldBeNil)
		msg = kfk.ConsumeUpdateSpuInfo(ctx)
		So(msg, ShouldNotBeNil)

		err = kfk.SendDeleteSpuInfo(ctx, testDeleteId)
		So(err, ShouldBeNil)
		msg = kfk.ConsumeDeleteSpuInfo(ctx)
		So(msg, ShouldNotBeNil)
	})
}
