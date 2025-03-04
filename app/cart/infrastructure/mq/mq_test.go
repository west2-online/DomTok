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
	"time"

	"github.com/bytedance/sonic"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/kafka"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

// 示例初始化方法，请根据你项目实际情况修改。
// 假设返回 *KafkaAdapter 或者你的 MQ 适配器。
func initTestKafkaAdapter(t *testing.T) *KafkaAdapter {
	t.Helper()
	config.Init("test")
	logger.Ignore()

	return NewKafkaAdapter(kafka.NewKafkaInstance())
}

func TestKafkaAdapter_SendAndConsumeAddGoods(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}

	kfk := initTestKafkaAdapter(t)
	ctx := context.Background()

	Convey("测试 Kafka 发送并消费 AddGoods 消息", t, func() {
		// 构建一个测试用的 GoodInfo
		// 这里的字段请根据你的项目结构体改动
		testGood := &model.GoodInfo{
			SkuId:     1001,
			ShopId:    2001,
			VersionId: 3001,
			Count:     2,
		}

		// 设定一个测试用的 UID 用来做分区
		testUID := int64(12345)
		var err error

		Convey("当调用 SendAddGoods 发送消息", func() {
			err = kfk.SendAddGoods(ctx, testUID, testGood)
			So(err, ShouldBeNil)

			Convey("应该能正常从 ConsumeAddGoods 消费到消息", func() {
				// 获取消费通道
				msgCh := kfk.ConsumeAddGoods(ctx)
				So(msgCh, ShouldNotBeNil)

				// 通过 select 和超时来确认是否成功消费到消息
				select {
				case receivedMsg := <-msgCh:
					So(receivedMsg, ShouldNotBeNil)

					// 反序列化消息体
					var addMsg model.AddGoodsMsg
					if unmarshalErr := sonic.Unmarshal(receivedMsg.V, &addMsg); unmarshalErr != nil {
						t.Errorf("反序列化失败: %v", unmarshalErr)
					}

					// 校验反序列化后的内容是否和发送时一致
					So(addMsg.Uid, ShouldEqual, testUID)
					So(addMsg.Goods, ShouldNotBeNil)
					So(addMsg.Goods.SkuId, ShouldEqual, testGood.SkuId)
					So(addMsg.Goods.ShopId, ShouldEqual, testGood.ShopId)
					So(addMsg.Goods.VersionId, ShouldEqual, testGood.VersionId)
					So(addMsg.Goods.Count, ShouldEqual, testGood.Count)

				case <-time.After(5 * time.Second):
					// 如果超过5秒没有收到消息，就认为测试失败
					t.Errorf("等待 AddGoods 消息超时，没有收到消息")
				}
			})
		})
	})
}
