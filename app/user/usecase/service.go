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

package usecase

import (
	"context"

	"github.com/west2-online/DomTok/app/user/entities"
)

// PersistencePort 表示持久化存储接口 (或者也可以叫做 DBPort)
type PersistencePort interface {
	IsUserExist(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, entity *entities.User) error
}

// CachePort 表示缓存接口
type CachePort interface{}

// TemplateRPCPort 表示模板服务的 RPC 接口
// 比如 template 服务有非常多的 RPC 接口，但是目前只需要一个接口，也就是 User 不依赖于它所不需要的接口
type TemplateRPCPort interface {
	GetTemplateInfo(ctx context.Context) (*entities.TemplateModel, error)
}

type UseCase struct {
	DB             PersistencePort
	templateClient TemplateRPCPort
}

func NewUserCase(db PersistencePort, tempate TemplateRPCPort) *UseCase {
	return &UseCase{
		DB:             db,
		templateClient: tempate,
	}
}
