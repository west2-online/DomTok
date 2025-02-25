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
	"github.com/west2-online/DomTok/app/commodity/controllers/rpc"
	"github.com/west2-online/DomTok/app/commodity/domain/service"
	"github.com/west2-online/DomTok/app/commodity/infrastructure/es"
	"github.com/west2-online/DomTok/app/commodity/infrastructure/mq"
	"github.com/west2-online/DomTok/app/commodity/infrastructure/mysql"
	"github.com/west2-online/DomTok/app/commodity/infrastructure/redis"
	"github.com/west2-online/DomTok/app/commodity/usecase"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/kafka"
	"github.com/west2-online/DomTok/pkg/utils"
)

func InjectCommodityHandlerr() commodity.CommodityService {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	sf, err := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
	if err != nil {
		panic(err)
	}

	redisCache, err := client.NewRedisClient(constants.RedisDBCommodity)
	if err != nil {
		panic(err)
	}

	elastic, err := client.NewEsCommodityClient()
	if err != nil {
		panic(err)
	}

	kafMQ := kafka.NewKafkaInstance()

	db := mysql.NewCommodityDB(gormDB)
	re := redis.NewCommodityCache(redisCache)
	kaf := mq.NewCommodityMQ(kafMQ)
	e := es.NewCommodityElastic(elastic)
	svc := service.NewCommodityService(db, sf, re, kaf, e)
	uc := usecase.NewCommodityCase(db, svc, re, kaf, e)

	return rpc.NewCommodityHandler(uc)
}
