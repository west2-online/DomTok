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
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/pkg/constants"
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
	s := Spu{
		Id:               spu.SpuId,
		Name:             spu.Name,
		CreatorId:        spu.CreatorId,
		Description:      spu.Description,
		CategoryId:       spu.CategoryId,
		GoodsHeadDrawing: spu.GoodsHeadDrawingUrl,
		Price:            spu.Price,
		ForSale:          spu.ForSale,
		Shipping:         spu.Shipping,
	}

	if err := db.client.WithContext(ctx).Table(s.TableName()).Create(&s).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to spu: %v", err)
	}
	return nil
}

func (db *commodityDB) CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) error {
	s := SpuImage{
		Id:    spuImage.ImageID,
		SpuId: spuImage.SpuID,
		Url:   spuImage.Url,
	}
	if err := db.client.WithContext(ctx).Table(s.TableName()).Create(&s).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to spu image: %v", err)
	}
	return nil
}

func (db *commodityDB) DeleteSpu(ctx context.Context, spuId int64) error {
	s := Spu{}
	if err := db.client.WithContext(ctx).Table(s.TableName()).Where("id = ?", spuId).Delete(&s).Error; err != nil {
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

func (db *commodityDB) GetImagesBySpuId(ctx context.Context, spuId int64, offset, limit int) ([]*model.SpuImage, int64, error) {
	imgs := make([]*SpuImage, 0)
	var cnt int64
	if err := db.client.WithContext(ctx).Table(constants.SpuImageTableName).Where("spu_id = ?", spuId).
		Order("created_at").Limit(limit).Offset(offset).Find(&imgs).Count(&cnt).Error; err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get images: %v", err)
	}
	ret := make([]*model.SpuImage, 0)
	for _, img := range imgs {
		ret = append(ret, &model.SpuImage{
			ImageID:   img.Id,
			SpuID:     img.SpuId,
			Url:       img.Url,
			CreatedAt: img.CreatedAt.Unix(),
			UpdatedAt: img.UpdatedAt.Unix(),
		})
	}
	return ret, cnt, nil
}

func (db *commodityDB) GetSpuBySpuId(ctx context.Context, spuId int64) (*model.Spu, error) {
	s := Spu{}
	if err := db.client.WithContext(ctx).Table(constants.SpuTableName).Where("id = ?", spuId).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceSpuNotExist, "spu not exist")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get spu: %v", err)
	}
	ret := &model.Spu{
		SpuId:               s.Id,
		Name:                s.Name,
		CreatorId:           s.CreatorId,
		CategoryId:          s.CategoryId,
		Description:         s.Description,
		GoodsHeadDrawingUrl: s.GoodsHeadDrawing,
		Price:               s.Price,
		ForSale:             s.ForSale,
		Shipping:            s.Shipping,
		CreatedAt:           s.CreatedAt.Unix(),
		UpdatedAt:           s.UpdatedAt.Unix(),
	}

	return ret, nil
}

func (db *commodityDB) UpdateSpu(ctx context.Context, spu *model.Spu) error {
	s := Spu{
		Id:               spu.SpuId,
		Name:             spu.Name,
		CreatorId:        spu.CreatorId,
		Description:      spu.Description,
		CategoryId:       spu.CategoryId,
		GoodsHeadDrawing: spu.GoodsHeadDrawingUrl,
		Price:            spu.Price,
		ForSale:          spu.ForSale,
		Shipping:         spu.Shipping,
	}
	if err := db.client.WithContext(ctx).Table(constants.SpuTableName).Where("id=?", spu.SpuId).Updates(&s).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update spu: %v", err)
	}
	return nil
}

func (db *commodityDB) UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error {
	img := SpuImage{
		Id:    spuImage.ImageID,
		SpuId: spuImage.SpuID,
		Url:   spuImage.Url,
	}
	if err := db.client.WithContext(ctx).Table(constants.SpuImageTableName).Updates(&img).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update spu image: %v", err)
	}
	return nil
}

func (db *commodityDB) DeleteSpuImage(ctx context.Context, spuImageId int64) error {
	s := SpuImage{}
	if err := db.client.WithContext(ctx).Table(s.TableName()).Where("id = ?", spuImageId).Delete(&s).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete spu image: %v", err)
	}
	return nil
}

func (db *commodityDB) GetSpuImage(ctx context.Context, spuImageId int64) (*model.SpuImage, error) {
	img := SpuImage{}

	if err := db.client.WithContext(ctx).Table(img.TableName()).Where("id=?", spuImageId).First(&img).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceImgNotExist, "spu image not exist")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get spu image: %v", err)
	}
	ret := &model.SpuImage{
		ImageID: img.Id,
		SpuID:   img.SpuId,
		Url:     img.Url,
	}
	return ret, nil
}

func (db *commodityDB) DeleteSpuImagesBySpuId(ctx context.Context, spuId int64) (ids []int64, url []string, err error) {
	ids = make([]int64, 0)
	url = make([]string, 0)
	imgs := make([]*SpuImage, 0)

	if err = db.client.WithContext(ctx).Table(constants.SpuImageTableName).Where("spu_id = ?", spuId).Delete(imgs).Error; err != nil {
		return nil, nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete images: %v", err)
	}
	for _, img := range imgs {
		ids = append(ids, img.SpuId)
		url = append(url, img.Url)
	}
	return ids, url, nil
}
