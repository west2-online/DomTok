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

	entities "github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/app/commodity/domain/service"
	Model "github.com/west2-online/DomTok/kitex_gen/model"
)

type CommodityUseCase interface {
	// 增删改查
	CreateCategory(ctx context.Context, entity *entities.Category) (int64, error)
	DeleteCategory(ctx context.Context, category *entities.Category) error
	UpdateCategory(ctx context.Context, category *entities.Category) error
	ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*Model.CategoryInfo, err error)
}

type useCase struct {
	db    repository.CommodityDB
	svc   *service.CommodityService
	cache repository.CommodityCache
}

func NewCommodityCase(db repository.CommodityDB, svc *service.CommodityService, cache repository.CommodityCache) *useCase {
	return &useCase{
		db:    db,
		svc:   svc,
		cache: cache,
	}
}
