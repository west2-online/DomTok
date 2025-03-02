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
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"strings"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/logger"
)

// GetRocketmqAdmin 获取一个 admin 实例以及 brokerAddr
func GetRocketmqAdmin() (admin.Admin, string) {
	nameSrvAddr := []string{config.Rocketmq.NameSrvAddr}
	brokerAddr := config.Rocketmq.BrokerAddr
	adm, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(nameSrvAddr)),
		admin.WithCredentials(primitive.Credentials{
			AccessKey: config.Rocketmq.AccessKey,
			SecretKey: config.Rocketmq.SecretKey,
		}),
	)
	if err != nil {
		logger.Fatalf("create rocketmq admin error: %v", err)
	}

	return adm, brokerAddr
}

// GetRocketmqProducer 获取一个 producer 实例
func GetRocketmqProducer() rocketmq.Producer {
	nameSrvAddr := []string{config.Rocketmq.NameSrvAddr}
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(nameSrvAddr),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: config.Rocketmq.AccessKey,
			SecretKey: config.Rocketmq.SecretKey,
		}),
	)
	if err != nil {
		logger.Fatalf("create rocketmq producer error: %v", err)
	}

	return p
}

// GetRocketmqPushConsumer 获取一个 push consumer 实例
//
// 参数：
// - group: 消费者组名称
// - opts: 消费者的可选配置项
func GetRocketmqPushConsumer(group string, opts ...consumer.Option) rocketmq.PushConsumer {
	nameSrvAddr := []string{config.Rocketmq.NameSrvAddr}
	if opts == nil {
		opts = make([]consumer.Option, 0)
	}
	opts = append(opts,
		consumer.WithGroupName(group),
		consumer.WithNameServer(nameSrvAddr),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: config.Rocketmq.AccessKey,
			SecretKey: config.Rocketmq.SecretKey,
		}),
	)

	p, err := rocketmq.NewPushConsumer(opts...)
	if err != nil {
		logger.Fatalf("create rocketmq push consumer error: %v", err)
	}
	return p
}

func SetRocketMqLoggerLevel(level string) {
	rlog.SetLogLevel(strings.ToLower(level))
}
