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

package kafka

import (
	"context"
	"errors"
	"io"

	"github.com/samber/lo"
	kafkago "github.com/segmentio/kafka-go"

	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
)

type Kafka struct {
	readers      []*kafkago.Reader
	writers      map[string]*kafkago.Writer
	consumeChans map[string]chan *Message
}

type Message struct {
	K, V []byte
}

// NewKafkaInstance 返回一个新的 kafka 实例
func NewKafkaInstance() *Kafka {
	return &Kafka{
		readers:      make([]*kafkago.Reader, 0),
		writers:      make(map[string]*kafkago.Writer, 0),
		consumeChans: make(map[string]chan *Message),
	}
}

// Consume 根据 consumerNum开启指定数量的协程, 并将消息通过 channel 传递
//
// 注意: 不要手动关闭返回的 channel
func (k *Kafka) Consume(ctx context.Context, topic string, consumerNum int, groupID string, chanCap ...int) <-chan *Message {
	if k.consumeChans[topic] != nil {
		return k.consumeChans[topic]
	}

	chCap := constants.DefaultConsumerChanCap
	if chanCap != nil {
		chCap = chanCap[0]
	}
	ch := make(chan *Message, chCap)
	k.consumeChans[topic] = ch

	for i := range consumerNum {
		k.readers = append(k.readers, client.GetNewReader(topic, groupID))
		go k.consume(ctx, topic, k.readers[i])
	}
	return ch
}

func (k *Kafka) consume(ctx context.Context, topic string, r *kafkago.Reader) {
	ch := k.consumeChans[topic]
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Errorf("read message from kafka reader failed,err: %v", err.Error())
			return
		}

		ch <- &Message{K: msg.Key, V: msg.Value}
	}
}

// Send 发送消息到指定的 topic
func (k *Kafka) Send(ctx context.Context, topic string, messages []*Message) []error {
	if k.writers[topic] == nil {
		if err := k.SetWriter(topic); err != nil {
			return []error{err}
		}
	}

	return k.send(ctx, topic, messages)
}

// SetWriter 针对特定的 topic 生成一个并发安全的 writer,
// SetWriter 会在 topic 不存在的时候创建他
func (k *Kafka) SetWriter(topic string, asyncWrite ...bool) error {
	async := constants.DefaultKafkaProductorSyncWrite
	if asyncWrite != nil {
		async = asyncWrite[0]
	}

	w, err := client.GetNewWriter(topic, async)
	if err != nil {
		return err
	}

	k.writers[topic] = w
	return nil
}

func (k *Kafka) send(ctx context.Context, topic string, messages []*Message) []error {
	msgs := lo.Map(messages, func(item *Message, index int) kafkago.Message {
		return kafkago.Message{
			Key:   item.K,
			Value: item.V,
		}
	})

	err := k.writers[topic].WriteMessages(ctx, msgs...)
	switch e := err.(type) { //nolint
	case nil:
		return nil
	case kafkago.WriteErrors:
		return e
	default:
		return []error{err}
	}
}

func (k *Kafka) Close() {
	for _, reader := range k.readers {
		if err := reader.Close(); err != nil {
			logger.Errorf("close kafka reader failed, err: %v", err)
		}
	}

	for _, writer := range k.writers {
		if err := writer.Close(); err != nil {
			logger.Errorf("close kafka writer failed, err: %v", err)
		}
	}

	for _, ch := range k.consumeChans {
		close(ch)
	}
}
