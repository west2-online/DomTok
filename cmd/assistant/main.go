package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/west2-online/DomTok/app/assistant/router"
	"github.com/west2-online/DomTok/pkg/constants"
)

// var serviceName = constants.AssistantServiceName

func init() {
	// config.Init(serviceName)
	// logger.Init(serviceName, config.GetLoggerLevel())
}

func main() {
	// get available port from config set
	//// listenAddr, err := utils.GetAvailablePort()
	// if err != nil {
	//	logger.Fatalf("get available port failed, err: %v", err)
	// }

	h := server.New(
		server.WithHostPorts(":8080"),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(constants.ServerMaxRequestBodySize),
	)

	//h.Use(
	//	mw.RecoveryMW(), // recovery
	//	mw.CorsMW(),     // cors
	//	mw.GzipMW(),     // gzip
	//	mw.SentinelMW(), // sentinel
	//)

	router.GeneratedRegister(h)
	h.Spin()
}
