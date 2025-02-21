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

package main

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino"
	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
	"github.com/west2-online/DomTok/app/assistant/cli/server/driver/http"
	"github.com/west2-online/DomTok/app/assistant/router"
	"github.com/west2-online/DomTok/app/assistant/service"
	"github.com/west2-online/DomTok/app/gateway/mw"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

var serviceName = "assistant"

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())

	ai := eino.NewClient()
	httpCli := http.NewClient(&http.ClientConfig{
		BaseUrl: `https://localhost:8888`,
	})
	// TODO: BaseUrl放在哪
	ai.SetServerStrategy(func(functionName string) adapter.ServerCaller {
		return httpCli
	})
	ai.SetBuilder(func(ctx context.Context) (model.ChatModel, error) {
		return ark.NewChatModel(ctx, &ark.ChatModelConfig{
			APIKey:  config.Volcengine.ApiKey,
			BaseURL: config.Volcengine.BaseUrl,
			Region:  config.Volcengine.Region,
			Model:   config.Volcengine.Model,
		})
	})

	service.Use(ai)
}

func main() {
	// get available port from config set
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("get available port failed, err: %v", err)
	}

	h := server.New(
		server.WithHostPorts(listenAddr),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(constants.ServerMaxRequestBodySize),
	)

	h.Use(
		mw.RecoveryMW(), // recovery
		mw.CorsMW(),     // cors
		mw.GzipMW(),     // gzip
		mw.SentinelMW(), // sentinel
	)

	router.GeneratedRegister(h)
	h.Spin()
}
