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

package mysql

import (
	"context"
	"github.com/west2-online/DomTok/pkg/constants"
	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/pkg/errno"
)

// commodityDB impl domain.CommodityDB defined domain
type commodityDB struct {
	client *gorm.DB
}

func NewCommodityDB(client *gorm.DB) repository.CommodityDB {
	return &commodityDB{client: client}
}

func (db *commodityDB) CreateCategory(ctx context.Context, name string) error {
	return nil
}

func (db *commodityDB) CreateSpu(ctx context.Context, spu *model.Spu) error {
	s := &Spu{
		Id:               spu.SpuId,
		Name:             spu.Name,
		CreatorId:        spu.CreatorId,
		Description:      spu.Description,
		CategoryId:       spu.CategoryId,
		GoodsHeadDrawing: spu.GoodsHeadDrawingName,
		Price:            spu.Price,
		ForSale:          spu.ForSale,
		Shipping:         spu.Shipping,
	}

	if err := db.client.WithContext(ctx).Table(s.TableName()).Create(spu).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to spu: %v", err)
	}
	return nil
}

func (db *commodityDB) CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) error {
	s := &SpuImage{
		Id:    spuImage.ImageID,
		SpuId: spuImage.SpuID,
		Url:   spuImage.Url,
	}
	if err := db.client.WithContext(ctx).Table(s.TableName()).Create(s).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to spu image: %v", err)
	}
	return nil
}

func (db *commodityDB) DeleteSpu(ctx context.Context, spuId int64) error {
	s := &Spu{}
	if err := db.client.WithContext(ctx).Table(s.TableName()).Where("spu_id = ?", spuId).Delete(s); err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete spu: %v", err)
	}
	return nil
}

func (db *commodityDB) IsExistSku(ctx context.Context, spuId int64) (bool, error) {
	var cnt int64
	if err := db.client.WithContext(ctx).Table(constants.SpuSkuTableName).Where("spu_id = ?", spuId).Count(&cnt).Error; err != nil {
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get count: %v", err)
	}
	return cnt != 0, nil
}

func (db *commodityDB) GetSpuBySpuId(ctx context.Context, spuId int64) (*model.Spu, error) {
	var s Spu
	if err := db.client.WithContext(ctx).Table(constants.SpuTableName).Where("spu_id = ?", spuId).Find(&s).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get spu: %v", err)
	}
	ret := &model.Spu{
		SpuId:                s.Id,
		Name:                 s.Name,
		CreatorId:            s.CreatorId,
		CategoryId:           s.CategoryId,
		Description:          s.Description,
		GoodsHeadDrawingName: s.GoodsHeadDrawing,
		Price:                s.Price,
		ForSale:              s.ForSale,
		Shipping:             s.Shipping,
		CreatedAt:            s.CreatedAt.Unix(),
		UpdatedAt:            s.UpdatedAt.Unix(),
		DeletedAt:            s.DeletedAt.Unix(),
	}

	return ret, nil
}

func (db *commodityDB) UpdateSpu(ctx context.Context, spu *model.Spu) error {
	if err := db.client.WithContext(ctx).Table(constants.SpuTableName).Where("spu_id = ?", spu.SpuId).Updates(spu); err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update spu: %v", err)
	}
	return nil
}

func (db *commodityDB) UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error {
	if err := db.client.WithContext(ctx).Table(constants.SpuImageTableName).Where("spu_id = ?", spuImage.ImageID).Updates(spuImage); err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update spu image: %v", err)
	}
	return nil
}

func (db *commodityDB) DeleteSpuImage(ctx context.Context, spuImageId int64) error {
	s := &SpuImage{}
	if err := db.client.WithContext(ctx).Table(s.TableName()).Where("spu_id = ?", spuImageId).Delete(s); err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete spu image: %v", err)
	}
	return nil
}

func (db *commodityDB) DeleteSpuImageToSpu(ctx context.Context, spuImageId int64, spuId int64) error {
	return nil
}

func (db *commodityDB) GetSpuImage(ctx context.Context, spuImageId int64) (*model.SpuImage, error) {
	//TODO implement me
	panic("implement me")
}
