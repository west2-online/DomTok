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

package constants

const (
	// LogFilePath 对应 ${pwd}/{LogFilePath}/log.log 相对于当前运行路径而言
	LogFilePath = "log"

	// wd/log/{ServiceName}/data/*.log
	LogFilePathTemplate      = "%s/%s/%s/%s/std.log"
	ErrorLogFilePathTemplate = "%s/%s/%s/%s/stderr.log"

	// DefaultLogLevel 是默认的日志等级. Supported Level: debug info warn error fatal
	DefaultLogLevel = "INFO"

	StackTraceKey = "stacktrace"
	ServiceKey    = "service"
	SourceKey     = "source"
	ErrorCodeKey  = "error_code"

	KlogSource  = "klog"
	MysqlSource = "mysql"
	RedisSource = "redis"
	KafkaSource = "kafka"
)
