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

package commodity

import (
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/DomTok/app/commodity/controllers/rpc"
	"github.com/west2-online/DomTok/app/commodity/repository/cache"
	"github.com/west2-online/DomTok/app/commodity/repository/db"
	"github.com/west2-online/DomTok/app/commodity/repository/es"
	"github.com/west2-online/DomTok/app/commodity/repository/mq"
	"github.com/west2-online/DomTok/app/commodity/usecase"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/kitex_gen/commodity/commodityservice"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
	"net"

	"github.com/west2-online/DomTok/pkg/logger"
)

// constants部分看了其他的pr有写了，防止冲突先找个数代替，到时候合并完再改掉
var (
	serviceName    = "constants.CommodityServiceName" //TODO
	serviceAdapter *usecase.UseCase
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())

	dbClient, _ := client.InitMySQL()
	dbAdapter := db.NewDBAdapter(dbClient)

	cacheClient, _ := client.NewRedisClient(0) // TODO
	cacheAdapter := cache.NewCacheAdapter(cacheClient)

	mqClient, _ := client.GetConn()
	mqAdapter := mq.NewMQAdapter(mqClient)

	//TODO: ES
	esClient := es.EsAdapter{}

	serviceAdapter = usecase.NewCommodityCase(dbAdapter, mqAdapter, cacheAdapter, esClient)
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("Commodity: new etcd registry failed, err: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("Commodity: get available port failed, err: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Commodity: resolve tcp addr failed, err: %v", err)
	}

	svr := commodityservice.NewServer(
		rpc.NewCommodityHandler(serviceAdapter),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: serviceName,
		}),
		server.WithMuxTransport(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{
			MaxConnections: constants.MaxConnections,
			MaxQPS:         constants.MaxQPS,
		}),
	)
	if err = svr.Run(); err != nil {
		logger.Fatalf("Commodity: run server failed, err: %v", err)
	}
}
