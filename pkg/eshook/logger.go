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

package eshook

import (
	elastic "github.com/elastic/go-elasticsearch"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// InitLoggerWithHook 初始化带有EsHook的logger
// index: 索引的名字
func InitLoggerWithHook(index string, esclient *elastic.Client) {
	if config.Elasticsearch == nil {
		return
	}

	if !client.IsESConnected(esclient) {
		logger.Warn("es not worked!")
		return
	}

	hook := NewElasticHook(esclient, config.Elasticsearch.Host, index)
	logger.AddLoggerHook(hook.Fire)
}
