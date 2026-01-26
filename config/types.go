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

package config

type ServerConfig struct {
	Secret      string `mapstructure:"private-key"`
	PublicKey   string `mapstructure:"public-key"`
	Version     string
	Name        string
	LogLevel    string `mapstructure:"log-level"`
	IntranetUrl string `mapstructure:"intranet-url"`
}

type SnowflakeConfig struct {
	DatacenterID int64 `mapstructure:"datacenter-id"`
}

type ServiceConfig struct {
	Name     string
	AddrList []string
	LB       bool `mapstructure:"load-balance"`
}

type MySQLConfig struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type JaegerConfig struct {
	Addr string
}

type EtcdConfig struct {
	Addr string
}

type RabbitMQConfig struct {
	Addr     string
	Username string
	Password string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type ElasticsearchConfig struct {
	Addr string
	Host string
}

type KafkaConfig struct {
	Address  string
	Network  string
	User     string
	Password string
}

type VolcengineConfig struct {
	ApiKey  string `mapstructure:"api-key"`
	BaseUrl string `mapstructure:"base-url"`
	Region  string `mapstructure:"region"`
	Model   string `mapstructure:"model"`
}

type TosConfig struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
}

type RocketmqConfig struct {
	BrokerAddr  string `mapstructure:"brokerAddr"`
	NameSrvAddr string `mapstructure:"nameSrvAddr"`
	AccessKey   string `mapstructure:"accessKey"`
	SecretKey   string `mapstructure:"secretKey"`
}

type OtelConfig struct {
	CollectorAddr string `mapstructure:"collector-addr"`
}

type AdministratorConfig struct {
	Secret string
}

type Config struct {
	Server        ServerConfig
	Snowflake     SnowflakeConfig
	MySQL         MySQLConfig
	Jaeger        JaegerConfig
	Etcd          EtcdConfig
	RabbitMQ      RabbitMQConfig
	Redis         RedisConfig
	Elasticsearch ElasticsearchConfig
	Kafka         KafkaConfig
	Volcengine    VolcengineConfig
	Tos           TosConfig
	Rocketmq      RocketmqConfig
	Otel          OtelConfig
	Administrator AdministratorConfig
}
