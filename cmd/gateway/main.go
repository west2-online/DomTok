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

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/opensergo/sentinel/adapter"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/gateway/router"
	"github.com/west2-online/DomTok/gateway/rpc"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
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

	h := server.New(
		server.WithHostPorts(listenAddr),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(constants.ServerMaxRequestBodySize),
	)

	h.Use(recovery.Recovery(recovery.WithRecoveryHandler(recoveryHandler)))

	// Cors
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		MaxAge:           constants.CorsMaxAge,
		ExposeHeaders:    []string{"Content-Length"},
	}))

	// gzip
	h.Use(gzip.Gzip(gzip.BestSpeed))

	// Sentinel
	initSentinel()
	h.Use(adapter.SentinelServerMiddleware(
		adapter.WithServerResourceExtractor(func(c context.Context, ctx *app.RequestContext) string {
			return "api"
		}),
		adapter.WithServerBlockFallback(func(ctx context.Context, c *app.RequestContext) {
			logger.Errorf("frequent requests have been rejected by the gateway. clientIP: %v\n", c.ClientIP())
			c.AbortWithStatusJSON(consts.StatusOK, map[string]interface{}{
				"code":    errno.InternalServiceErrorCode,
				"message": "服务器当前处于请求高峰，请稍后再试",
			})
		}),
	))

	router.GeneratedRegister(h)
	h.Spin()
}

func recoveryHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
	logger.Errorf("[Recovery] InternalServiceError err=%v\n stack=%s\n", err, stack)
	c.JSON(consts.StatusInternalServerError, map[string]interface{}{
		"code":    errno.InternalServiceErrorCode,
		"message": "内部服务错误，请稍后再试",
	})
}

func initSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		logger.Fatalf("Unexpected error: %+v", err)
	}

	// limit QPS to 100
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "api",
			Threshold:              constants.SentinelThreshold,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       constants.SentinelStatIntervalInMs,
		},
	})
	if err != nil {
		logger.Fatalf("Unexpected error: %+v", err)
		return
	}
}
