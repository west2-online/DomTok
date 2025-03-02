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

package cart

import (
	"github.com/west2-online/DomTok/app/cart/controllers/rpc"
	"github.com/west2-online/DomTok/app/cart/domain/service"
	"github.com/west2-online/DomTok/app/cart/infrastructure/cache"
	"github.com/west2-online/DomTok/app/cart/infrastructure/db"
	"github.com/west2-online/DomTok/app/cart/infrastructure/mq"
	rpccli "github.com/west2-online/DomTok/app/cart/infrastructure/rpc"
	"github.com/west2-online/DomTok/app/cart/usecase"
	"github.com/west2-online/DomTok/kitex_gen/cart"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
	"github.com/west2-online/DomTok/pkg/logger"
)

// InjectCartHandler 注入外部调用
func InjectCartHandler() cart.CartService {
	// 何尝不是一种clientSet
	dbClient, _ := client.InitMySQL()
	dbAdapter := db.NewDBAdapter(dbClient)
	cacheClient, _ := client.NewRedisClient(constants.RedisDBCart)
	cacheAdapter := cache.NewCacheAdapter(cacheClient)
	kafkaAdapter := mq.NewKafkaAdapter(kafka.NewKafkaInstance())
	cClient, err := client.InitCommodityRPC()
	if err != nil {
		logger.Errorf("Failed to init commodity rpc client: %v", err)
	}
	oClient, err := client.InitOrderRPC()
	if err != nil {
		logger.Errorf("Failed to init order rpc client: %v", err)
	}
	rpcImpl := rpccli.NewCartRpcImpl(*cClient, *oClient)
	svc := service.NewCartService(dbAdapter, cacheAdapter, kafkaAdapter, rpcImpl)
	serviceAdapter := usecase.NewCartCase(dbAdapter, cacheAdapter, kafkaAdapter, rpcImpl, svc)
	return rpc.NewCartHandler(serviceAdapter)
}
