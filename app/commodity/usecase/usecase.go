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

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/app/commodity/domain/service"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
)

type CommodityUseCase interface {
	CreateCategory(ctx context.Context, category *model.Category) (int64, error)
	DeleteCategory(ctx context.Context, category *model.Category) (err error)
	UpdateCategory(ctx context.Context, category *model.Category) (err error)
	ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*model.CategoryInfo, err error)
	CreateSpu(ctx context.Context, spu *model.Spu) (id int64, err error)
	CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) (int64, error)
	DeleteSpu(ctx context.Context, spuId int64) error
	UpdateSpu(ctx context.Context, spu *model.Spu) error
	UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error
	DeleteSpuImage(ctx context.Context, imageId int64) error
	ViewSpuImages(ctx context.Context, spuId int64, offset, limit int) ([]*model.SpuImage, int64, error)
	ViewSpus(ctx context.Context, req *commodity.ViewSpuReq) ([]*model.Spu, int64, error)
	ListSpuInfo(ctx context.Context, ids []int64) ([]*model.Spu, error)

	IncrLockStock(ctx context.Context, infos []*model.SkuBuyInfo) error
	DecrLockStock(ctx context.Context, infos []*model.SkuBuyInfo) error
	DecrStock(ctx context.Context, infos []*model.SkuBuyInfo) error

	CreateCoupon(ctx context.Context, coupon *model.Coupon) (int64, error)
	DeleteCoupon(ctx context.Context, coupon *model.Coupon) (err error)
	GetCreatorCoupons(ctx context.Context, pageNum int64) (coupons []*model.Coupon, err error)
	CreateUserCoupon(ctx context.Context, coupon *model.UserCoupon) (err error)
	SearchUserCoupons(ctx context.Context, pageNum int64) (coupons []*model.Coupon, err error)
	GetCouponAndPrice(ctx context.Context, goods []*model.OrderGoods) ([]*model.OrderGoods, float64, error)

	CreateSku(ctx context.Context, sku *model.Sku, ext string) (s *model.Sku, err error)
	UpdateSku(ctx context.Context, sku *model.Sku, ext string) (err error)
	DeleteSku(ctx context.Context, sku *model.Sku) (err error)
	ViewSku(ctx context.Context, sku *model.Sku, pageNum *int64, pageSize *int64, isSpuId bool) (Skus []*model.Sku, total int64, err error)
	UploadSkuAttr(ctx context.Context, attr *model.AttrValue, Sku *model.Sku) (err error)
	ListSkuInfo(ctx context.Context, skuInfos []*model.SkuVersion, pageNum int64, pageSize int64) (SkuInfos []*model.Sku, total int64, err error)
	ViewSkuPriceHistory(ctx context.Context, skuPrice *model.SkuPriceHistory, pageNum int64, pageSize int64) ([]*model.SkuPriceHistory, error)
	CreateSkuImage(ctx context.Context, skuImage *model.SkuImage, data []byte) (int64, error)
	UpdateSkuImage(ctx context.Context, skuImage *model.SkuImage, data []byte) (err error)
	ViewSkuImages(ctx context.Context, sku *model.Sku, pageNum *int64, pageSize *int64) (Images []*model.SkuImage, total int64, err error)
	DeleteSkuImage(ctx context.Context, imageId int64) (err error)
}

type useCase struct {
	db    repository.CommodityDB
	svc   *service.CommodityService
	cache repository.CommodityCache
	mq    repository.CommodityMQ
	es    repository.CommodityElastic
}

func NewCommodityCase(db repository.CommodityDB, svc *service.CommodityService, cache repository.CommodityCache,
	mq repository.CommodityMQ, es repository.CommodityElastic,
) *useCase {
	return &useCase{
		db:    db,
		svc:   svc,
		cache: cache,
		mq:    mq,
		es:    es,
	}
}
