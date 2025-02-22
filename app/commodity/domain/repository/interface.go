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

package repository

import (
	"context"

	kmodel "github.com/west2-online/DomTok/kitex_gen/model"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/kafka"
)

type CommodityDB interface {
	IsCategoryExist(ctx context.Context, name string) (bool, error)
	CategoryCreatorId(ctx context.Context, Id int64) (int64,error)
	CreateCategory(ctx context.Context, entity *model.Category) error
	DeleteCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*kmodel.CategoryInfo, err error)

	CreateSpu(ctx context.Context, spu *model.Spu) error
	CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) error
	DeleteSpu(ctx context.Context, spuId int64) error
	IsExistSku(ctx context.Context, spuId int64) (bool, error)
	GetSpuBySpuId(ctx context.Context, spuId int64) (*model.Spu, error)
	GetSpuImage(ctx context.Context, spuImageId int64) (*model.SpuImage, error)
	UpdateSpu(ctx context.Context, spu *model.Spu) error
	UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error
	DeleteSpuImage(ctx context.Context, spuImageId int64) error
	DeleteSpuImagesBySpuId(ctx context.Context, spuId int64) (ids []int64, url []string, err error)
	GetImagesBySpuId(ctx context.Context, spuId int64, offset, limit int) ([]*model.SpuImage, int64, error)
}

type CommodityCache interface {
	IsExist(ctx context.Context, key string) bool
	GetSpuImages(ctx context.Context, key string) (*model.SpuImages, error)
	SetSpuImages(ctx context.Context, key string, images *model.SpuImages)
}

type CommodityMQ interface {
	Send(ctx context.Context, topic string, message []*kafka.Message) error
}
