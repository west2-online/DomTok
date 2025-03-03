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

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/west2-online/DomTok/app/gateway/mw"
	"github.com/west2-online/DomTok/app/gateway/router"
	"github.com/west2-online/DomTok/app/gateway/rpc"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/middleware"
	"github.com/west2-online/DomTok/pkg/utils"
)

var serviceName = constants.GatewayServiceName

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())
	rpc.Init()
}

func main() {
	// get available port from config set
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("get available port failed, err: %v", err)
	}

	p := middleware.TelemetryProvider(serviceName, config.Otel.CollectorAddr)
	defer func() { logger.LogError(p.Shutdown(context.Background())) }()

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
