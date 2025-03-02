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
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/pkg/kafka"
)

type CommodityDB interface {
	IsCategoryExistByName(ctx context.Context, name string) (bool, error)
	GetCategoryById(ctx context.Context, id int64) (*model.Category, error)
	CreateCategory(ctx context.Context, entity *model.Category) error
	DeleteCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*model.CategoryInfo, err error)

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
	GetSpuByIds(ctx context.Context, spuIds []int64) ([]*model.Spu, error)

	CreateCoupon(ctx context.Context, coupon *model.Coupon) (int64, error)
	GetCouponById(ctx context.Context, id int64) (bool, *model.Coupon, error)
	GetCouponsByCreatorId(ctx context.Context, uid int64, pageNum int64) ([]*model.Coupon, error)
	DeleteCouponById(ctx context.Context, coupon *model.Coupon) error

	GetCouponsByIDs(ctx context.Context, couponIDs []int64) ([]*model.Coupon, error)
	CreateUserCoupon(ctx context.Context, coupon *model.UserCoupon) error
	GetUserCouponsByUId(ctx context.Context, uid int64, pageNum int64) ([]*model.UserCoupon, error)
	GetFullUserCouponsByUId(ctx context.Context, uid int64) ([]*model.UserCoupon, error)
	DeleteUserCoupon(ctx context.Context, coupon *model.UserCoupon) error

	IncrLockStock(ctx context.Context, infos []*model.SkuBuyInfo) error
	DecrLockStock(ctx context.Context, infos []*model.SkuBuyInfo) error
	IncrStock(ctx context.Context, infos []*model.SkuBuyInfo) error
	DecrStock(ctx context.Context, infos []*model.SkuBuyInfo) error
	GetSkuById(ctx context.Context, id int64) (*model.Sku, error)
	DecrStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error
	DecrLockStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error
	IncrLockStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error

	CreateSku(ctx context.Context, sku *model.Sku) error
	UpdateSku(ctx context.Context, sku *model.Sku) error
	ViewSku(ctx context.Context, skuIds []*int64, PageNum int, PageSize int) ([]*model.Sku, int64, error)
	DeleteSku(ctx context.Context, sku *model.Sku) error
	CreateSkuImage(ctx context.Context, skuImage *model.SkuImage) error
	UpdateSkuImage(ctx context.Context, skuImage *model.SkuImage) error
	ViewSkuImage(ctx context.Context, sku *model.Sku, PageNum int, PageSize int) ([]*model.SkuImage, int64, error)
	DeleteSkuImage(ctx context.Context, imageId int64) error
	IsSpuExist(ctx context.Context, spuId int64) (bool, error)
	GetSkuBySkuId(ctx context.Context, skuId int64) (*model.Sku, error)
	GetSkuImageByImageId(ctx context.Context, imageId int64) (*model.SkuImage, error)
	GetSkuIdBySpuID(ctx context.Context, spuId int64, PageNum int, PageSize int) ([]*int64, error)
	UploadSkuAttr(ctx context.Context, sku *model.Sku, attr *model.AttrValue, id int64) error
	ListSkuInfo(ctx context.Context, skuInfo []*model.SkuVersion, PageNum int, PageSize int) ([]*model.Sku, error)
}

type CommodityCache interface {
	IsExist(ctx context.Context, key string) bool
	GetSpuImages(ctx context.Context, key string) (*model.SpuImages, error)
	SetSpuImages(ctx context.Context, key string, images *model.SpuImages)

	SetSkuImages(ctx context.Context, key string, skuImages []*model.SkuImage)
	GetSkuImages(ctx context.Context, key string) ([]*model.SkuImage, error)
	DeleteSkuImages(ctx context.Context, key string) error

	GetLockStockNum(ctx context.Context, key string) (int64, error)
	SetLockStockNum(ctx context.Context, key string, num int64)
	IncrLockStockNum(ctx context.Context, infos []*model.SkuBuyInfo) error
	DecrLockStockNum(ctx context.Context, infos []*model.SkuBuyInfo) error
	GetLockStockKey(id int64) string
	GetStockKey(id int64) string
	DecrStockNum(ctx context.Context, infos []*model.SkuBuyInfo) error
	IsHealthy(ctx context.Context) error
	Lock(ctx context.Context, keys []string, ttl time.Duration) error
	UnLock(ctx context.Context, keys []string) error
	GetSkuKey(id int64) string
}

type CommodityMQ interface {
	Send(ctx context.Context, topic string, message []*kafka.Message) error
	SendCreateSpuInfo(ctx context.Context, spu *model.Spu) error
	SendUpdateSpuInfo(ctx context.Context, spu *model.Spu) error
	SendDeleteSpuInfo(ctx context.Context, id int64) error
	ConsumeCreateSpuInfo(ctx context.Context) <-chan *kafka.Message
	ConsumeUpdateSpuInfo(ctx context.Context) <-chan *kafka.Message
	ConsumeDeleteSpuInfo(ctx context.Context) <-chan *kafka.Message
}

type CommodityElastic interface {
	IsExist(ctx context.Context, indexName string) bool
	CreateIndex(ctx context.Context, indexName string) error
	AddItem(ctx context.Context, indexName string, spu *model.Spu) error
	RemoveItem(ctx context.Context, indexName string, id int64) error
	UpdateItem(ctx context.Context, indexName string, spu *model.Spu) error
	SearchItems(ctx context.Context, indexName string, query *commodity.ViewSpuReq) ([]int64, int64, error)
	BuildQuery(req *commodity.ViewSpuReq) *elastic.BoolQuery
}
