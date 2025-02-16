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

package order

import (
	"github.com/west2-online/DomTok/app/order/controllers/rpc"
	"github.com/west2-online/DomTok/app/order/domain/service"
	"github.com/west2-online/DomTok/app/order/infrastructure/mysql"
	"github.com/west2-online/DomTok/app/order/usecase"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/kitex_gen/order"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

// InjectOrderHandler 用于依赖注入
func InjectOrderHandler() order.OrderService {
	// 1. 初始化数据库连接
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}

	// 2. 初始化雪花算法
	sf, err := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfOrderService)
	if err != nil {
		panic(err)
	}

	// 3. 初始化各层依赖
	db := mysql.NewOrderDB(gormDB)
	svc := service.NewOrderService(db, sf)
	uc := usecase.NewOrderCase(db, svc)

	// 4. 返回 handler
	return rpc.NewOrderHandler(uc)
}
