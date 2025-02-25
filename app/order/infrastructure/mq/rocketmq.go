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
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/samber/lo"

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

// rocketMq 结构体封装了与 RocketMQ 交互所需的管理客户端、Broker 地址、生产者和消费者信息
type rocketMq struct {
	producers map[string]rocketmq.Producer
	consumers []rocketmq.PushConsumer
}

func NewRocketmq() repository.MQ {
	mq := &rocketMq{
		producers: make(map[string]rocketmq.Producer),
		consumers: make([]rocketmq.PushConsumer, 0),
	}
	return mq
}

// SendSyncMsg 同步发送一批消息到指定主题。如果消息列表为空则直接返回 nil，若发送过程中出现错误，将返回自定义错误
//
// 参数：
//   - ctx: 上下文对象，用于控制操作的生命周期
//   - topic: 消息要发送到的主题名称
//   - msgs: 要发送的消息列表，每个消息为 *model.MqMessage 类型
func (mq *rocketMq) SendSyncMsg(ctx context.Context, topic string, msgs ...*model.MqMessage) error {
	if len(msgs) == 0 {
		return nil
	}
	producer, err := mq.getProducer(topic)
	if err != nil {
		return err
	}

	if _, err = producer.SendSync(ctx, convertMsg(topic, msgs...)...); err != nil {
		return errno.NewErrNo(errno.InternalRocketmqErrorCode, fmt.Sprintf("failed to send sync msg, err: %v", err))
	}
	return nil
}

// SubscribeTopic 订阅指定主题的消息，启动一个 goroutine 来执行传入的回调函数 fn。
// 当 fn 返回 false 时，会将一批（或一条）消息重新放回队列等待再次消费。
//
// 参数：
//   - ctx: 上下文对象，用于控制操作的生命周期
//   - topic: 要订阅的主题名称
//   - pullMsgInterval: 拉取消息的时间间隔
//   - fn: 处理消息的回调函数，接收消息体并返回布尔值表示消息是否处理成功
func (mq *rocketMq) SubscribeTopic(ctx context.Context, topic string, pullMsgInterval time.Duration, fn func(ctx context.Context, body []byte) bool) error {
	con, err := mq.getConsumer(topic, pullMsgInterval)
	if err != nil {
		return err
	}

	err = con.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			if !fn(ctx, msg.Body) {
				return consumer.ConsumeRetryLater, nil
			}
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return errno.NewErrNo(errno.InternalRocketmqErrorCode, fmt.Sprintf("failed to subscribe topic, err: %v", err))
	}

	return nil
}

// Shutdown 尝试释放所有资源
func (mq *rocketMq) Shutdown() []error {
	var errs []error
	for _, con := range mq.consumers {
		if err := con.Shutdown(); err != nil {
			errs = append(errs, err)
		}
	}

	for _, producer := range mq.producers {
		if err := producer.Shutdown(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

// getProducer 获取指定主题的生产者实例。如果该主题的生产者已存在则直接返回，否则创建并启动一个新的生产者实例
// 参数：
// - topic: 要获取生产者的主题名称
// 返回值：
// - rocketmq.Producer: 生产者实例
// - error: 若启动生产者失败，返回自定义错误
func (mq *rocketMq) getProducer(topic string) (rocketmq.Producer, error) {
	if _, ok := mq.producers[topic]; ok {
		return mq.producers[topic], nil
	}
	producer := client.GetRocketmqProducer()
	if err := producer.Start(); err != nil {
		return nil, errno.NewErrNo(errno.InternalRocketmqErrorCode, fmt.Sprintf("start rocketmq producer failed, err: %v", err))
	}
	mq.producers[topic] = producer
	return producer, nil
}

// getConsumer 获取指定主题的消费者实例，并按照指定的拉取间隔进行配置，然后启动该消费者实例
// 参数：
// - topic: 要获取消费者的主题名称
// - pullMsgInterval: 拉取消息的时间间隔
// 返回值：
// - rocketmq.PushConsumer: 消费者实例
// - error: 若启动消费者失败，返回自定义错误
func (mq *rocketMq) getConsumer(topic string, pullMsgInterval time.Duration) (rocketmq.PushConsumer, error) {
	con := client.GetRocketmqPushConsumer(fmt.Sprintf(constants.OrderMqConsumerGroupFormat, topic), consumer.WithPullInterval(pullMsgInterval))
	if err := con.Start(); err != nil {
		return nil, errno.NewErrNo(errno.InternalRocketmqErrorCode, fmt.Sprintf("start rocketmq consumer failed, err: %v", err))
	}
	return con, nil
}

func convertMsg(topic string, msgs ...*model.MqMessage) []*primitive.Message {
	return lo.Map(msgs, func(item *model.MqMessage, index int) *primitive.Message {
		m := &primitive.Message{Body: item.Body, Topic: topic}
		if item.IsSetDelayLevel() {
			m.WithDelayTimeLevel(item.GetDelayLevel())
		}
		return m
	})
}
