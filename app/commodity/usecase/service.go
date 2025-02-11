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

	"github.com/west2-online/DomTok/app/commodity/entities"
	"github.com/west2-online/DomTok/kitex_gen/model"
)

type PersistencePort interface {
	IsCategoryExist(ctx context.Context, Name string) (bool, error)
	CreateCategory(ctx context.Context, entity *entities.Category) error
	DeleteCategory(ctx context.Context, category *entities.Category) error
	UpdateCategory(ctx context.Context, category *entities.Category) error
	ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*model.CategoryInfo, err error)
}
type CachePort interface{}

type MQPort interface{}

type EsPort interface{}

type UseCase struct {
	DB    PersistencePort
	MQ    MQPort
	Cache CachePort
	Es    EsPort
}

func NewCommodityCase(db PersistencePort, mq MQPort, cache CachePort, es EsPort) *UseCase {
	return &UseCase{
		DB:    db,
		MQ:    mq,
		Cache: cache,
		Es:    es,
	}
}
