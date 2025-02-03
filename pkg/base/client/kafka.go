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

package client

import (
	"fmt"
	"net"
	"time"

	kafukago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

const (
	Timeout = 10 * time.Second // 默认超时时间
)

// GetConn conn不能保证并发安全,仅可作为单线程的长连接使用。
func GetConn() (*kafukago.Conn, error) {
	dialer := getDialer()

	conn, err := dialer.Dial(config.Kafka.Network, config.Kafka.Address)
	if err != nil {
		return nil, fmt.Errorf("failed dial kafka server,error: %w", err)
	}
	return conn, nil
}

// GetNewReader 创建一个reader示例，reader是并发安全的
func GetNewReader(topic string, groupID ...string) *kafukago.Reader {
	id := constants.DefaultReaderGroupID
	if groupID != nil {
		id = groupID[0]
	}

	cfg := kafukago.ReaderConfig{
		Brokers:     []string{config.Kafka.Address}, // 单节点无Leader
		Topic:       topic,
		GroupID:     id,
		MaxBytes:    constants.KafkaReadMaxBytes, // 同上
		MaxAttempts: constants.KafkaRetries,
		Dialer:      getDialer(),
	}
	return kafukago.NewReader(cfg)
}

// GetNewWriter 创建一个writer示例，writer是并发安全的。 errLogger可以传入带有es hook的logger
func GetNewWriter() (*kafukago.Writer, error) {
	addr, err := net.ResolveTCPAddr(config.Kafka.Network, config.Kafka.Address)
	if err != nil {
		return nil, fmt.Errorf("failed create kafka writer,error: %w", err)
	}

	return &kafukago.Writer{
		Addr:                   addr,
		Balancer:               &kafukago.RoundRobin{}, // 轮询写入分区
		MaxAttempts:            constants.KafkaRetries, // 最大尝试次数
		RequiredAcks:           kafukago.RequireOne,    // 每个消息需要一次Act
		Async:                  true,                   // 异步写入
		AllowAutoTopicCreation: false,                  // 不允许自动创建分区
		Transport:              getTransport(),
	}, nil
}

func getDialer() *kafukago.Dialer {
	mechanism := plain.Mechanism{
		Username: config.Kafka.User,
		Password: config.Kafka.Password,
	}
	return &kafukago.Dialer{
		Timeout:       Timeout,
		DualStack:     true,
		SASLMechanism: mechanism,
	}
}

func getTransport() *kafukago.Transport {
	mechanism := plain.Mechanism{
		Username: config.Kafka.User,
		Password: config.Kafka.Password,
	}
	return &kafukago.Transport{
		SASL: mechanism,
	}
}
