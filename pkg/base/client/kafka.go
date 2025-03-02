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

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

const (
	Timeout = 3 * time.Second // 默认超时时间
)

// GetConn conn不能保证并发安全,仅可作为单线程的长连接使用。
func GetConn() (*kafukago.Conn, error) {
	dialer := getDialer()

	conn, err := dialer.Dial(config.Kafka.Network, config.Kafka.Address)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalKafkaErrorCode, fmt.Sprintf("failed dial kafka server,error: %v", err))
	}
	return conn, nil
}

// GetNewReader 创建一个reader示例，reader是并发安全的
func GetNewReader(topic string, groupID string) *kafukago.Reader {
	if groupID == "" {
		groupID = constants.DefaultReaderGroupID
	}
	if err := createIfNotExist(topic); err != nil {
		logger.Fatalf("Failed to create topic %s, %v", topic, err)
	}

	cfg := kafukago.ReaderConfig{
		Brokers:     []string{config.Kafka.Address}, // 单节点无Leader
		Topic:       topic,
		GroupID:     groupID,
		MaxBytes:    constants.KafkaReadMaxBytes, // 同上
		MaxAttempts: constants.KafkaRetries,
		ErrorLogger: logger.GetKafkaErrorLogger(),
		Dialer:      getDialer(),
	}
	return kafukago.NewReader(cfg)
}

// GetNewWriter 创建一个writer示例，writer是并发安全的。
func GetNewWriter(topic string, async bool) (*kafukago.Writer, error) {
	if err := createIfNotExist(topic); err != nil {
		return nil, err
	}

	addr, err := net.ResolveTCPAddr(config.Kafka.Network, config.Kafka.Address)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalNetworkErrorCode, fmt.Sprintf("failed create kafka writer,error: %v", err))
	}

	return &kafukago.Writer{
		Addr:                   addr,
		Topic:                  topic,
		Balancer:               &kafukago.RoundRobin{}, // 轮询写入分区
		MaxAttempts:            constants.KafkaRetries, // 最大尝试次数
		RequiredAcks:           kafukago.RequireOne,    // 每个消息需要一次Act
		Async:                  async,                  // 异步写入
		AllowAutoTopicCreation: true,
		ErrorLogger:            logger.GetKafkaErrorLogger(),
		Transport:              getTransport(),
	}, nil
}

func createIfNotExist(topic string) error {
	conn, err := GetConn()
	if err != nil {
		return err
	}

	err = conn.CreateTopics(kafukago.TopicConfig{
		Topic:             topic,
		NumPartitions:     constants.DefaultKafkaNumPartitions,
		ReplicationFactor: constants.DefaultKafkaReplicationFactor,
	})
	if err != nil {
		return errno.NewErrNo(errno.InternalKafkaErrorCode, fmt.Sprintf("failed to create topic, err: %v", err))
	}
	return nil
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
