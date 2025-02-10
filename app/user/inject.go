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

package user

import (
	"github.com/west2-online/DomTok/app/user/controllers/rpc"
	"github.com/west2-online/DomTok/app/user/repository/mysql"
	"github.com/west2-online/DomTok/app/user/service"
	"github.com/west2-online/DomTok/app/user/usecase"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/kitex_gen/user"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

// InjectUserHandler 用于依赖注入
// 从这个文件的位置就可以看出来极其特殊, 独立于架构之外, 服务于业务
func InjectUserHandler() user.UserService {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	sf, err := utils.NewSnowflake(config.GetDataCenterID(), constants.WorkerOfUserService)
	if err != nil {
		panic(err)
	}

	db := mysql.NewUserDB(gormDB)
	svc := service.NewUserService(db, sf)
	uc := usecase.NewUserCase(db, svc)

	return rpc.NewUserHandler(uc)
}
