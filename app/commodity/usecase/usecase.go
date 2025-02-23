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
	"github.com/west2-online/DomTok/kitex_gen/commodity"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/app/commodity/domain/service"
)

type CommodityUseCase interface {
	CreateCategory(ctx context.Context, category *model.Category) (id int64, err error)
	CreateSpu(ctx context.Context, spu *model.Spu) (id int64, err error)
	CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) (int64, error)
	DeleteSpu(ctx context.Context, spuId int64) error
	UpdateSpu(ctx context.Context, spu *model.Spu) error
	UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error
	DeleteSpuImage(ctx context.Context, imageId int64) error
	ViewSpuImages(ctx context.Context, spuId int64, offset, limit int) ([]*model.SpuImage, int64, error)
	ViewSpus(ctx context.Context, req *commodity.ViewSpuReq) ([]*model.Spu, int64, error)
}

type useCase struct {
	db    repository.CommodityDB
	svc   *service.CommodityService
	cache repository.CommodityCache
	mq    repository.CommodityMQ
	es    repository.CommodityElastic
}

func NewCommodityCase(db repository.CommodityDB, svc *service.CommodityService, cache repository.CommodityCache, mq repository.CommodityMQ, es repository.CommodityElastic) *useCase {
	return &useCase{
		db:    db,
		svc:   svc,
		cache: cache,
		mq:    mq,
		es:    es,
	}
}
