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
	"log"

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

func (db *commodityDB) IsCategoryExistByName(ctx context.Context, name string) (bool, error) {
	var category model.Category
	err := db.client.WithContext(ctx).Where("Name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errno.Errorf(errno.ErrRecordNotFound, "mysql: ErrRecordNotFound record not found: %v", err)
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query category: %v", err)
	}
	return true, nil
}

func (db *commodityDB) IsCategoryExistById(ctx context.Context, id int64) (bool, error) {
	var category model.Category
	err := db.client.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errno.Errorf(errno.ErrRecordNotFound, "mysql: ErrRecordNotFound record not found: %v", err)
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query category: %v", err)
	}
	return true, nil
}

func (db *commodityDB) CreateCategory(ctx context.Context, entity *model.Category) error {
	m := Category{
		Id:        entity.Id,
		Name:      entity.Name,
		CreatorId: entity.CreatorId,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: gorm.DeletedAt{},
	}
	if err := db.client.WithContext(ctx).Create(m).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create category: %v", err)
	}
	return nil
}

func (db *commodityDB) DeleteCategory(ctx context.Context, category *model.Category) error {
	if err := db.client.WithContext(ctx).Delete(Category{Id: category.Id}).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete category: %v", err)
	}
	return nil
}

func (db *commodityDB) GetSpuByIds(ctx context.Context, spuIds []int64) ([]*model.Spu, error) {
	spus := make([]*Spu, 0)
	if err := db.client.WithContext(ctx).Table(constants.SpuTableName).Where("id in (?)", spuIds).Find(&spus).Error; err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "CommodityDB.GetSpuByIds failed: %v", err)
	}
	rets := make([]*model.Spu, 0)
	for _, spu := range spus {
		ret := &model.Spu{
			SpuId:               spu.Id,
			Name:                spu.Name,
			CreatorId:           spu.CreatorId,
			Description:         spu.Description,
			CategoryId:          spu.CategoryId,
			Price:               spu.Price,
			ForSale:             spu.ForSale,
			Shipping:            spu.Shipping,
			CreatedAt:           spu.CreatedAt.Unix(),
			UpdatedAt:           spu.UpdatedAt.Unix(),
			GoodsHeadDrawingUrl: spu.GoodsHeadDrawing,
		}
		rets = append(rets, ret)
	}
	return rets, nil
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

func (db *commodityDB) CreateSku(ctx context.Context, sku *model.Sku) (err error) {
	s := &Sku{
		Id:               sku.SkuID,
		CreatorId:        sku.CreatorID,
		Price:            sku.Price,
		Name:             sku.Name,
		Description:      sku.Description,
		ForSale:          sku.ForSale,
		Stock:            sku.Stock,
		LockStock:        sku.Stock,
		StyleHeadDrawing: sku.StyleHeadDrawingUrl,
		HistoryVersionId: sku.HistoryID,
	}

	skuToSpu := &SpuToSku{
		SkuId: sku.SkuID,
		SpuId: sku.SpuID,
	}
	skuPriceHistory := &SkuPriceHistory{
		Id:          sku.HistoryID,
		SkuId:       sku.SkuID,
		MarkPrice:   sku.Price,
		PrevVersion: 0,
	}

	if err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(s.TableName()).Create(s).Error; err != nil {
			return err
		}
		if err := tx.Table(skuToSpu.TableName()).Create(skuToSpu).Error; err != nil {
			return err
		}
		if err := tx.Table(skuPriceHistory.TableName()).Create(skuPriceHistory).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create sku: %v", err)
	}

	return nil
}

func (db *commodityDB) UpdateSku(ctx context.Context, sku *model.Sku) error {
	s := &Sku{
		Id:               sku.SkuID,
		CreatorId:        sku.CreatedAt,
		StyleHeadDrawing: sku.StyleHeadDrawingUrl,
		Description:      sku.Description,
		Price:            sku.Price,
		ForSale:          sku.ForSale,
		Stock:            sku.Stock,
		HistoryVersionId: sku.HistoryID,
	}

	skuPriceHistory := &SkuPriceHistory{
		Id:          sku.HistoryID,
		SkuId:       sku.SkuID,
		MarkPrice:   sku.Price,
		PrevVersion: 0,
	}

	if err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(s.TableName()).Where("id = ?", sku.SkuID).Updates(s).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.Errorf(errno.ServiceSkuNotExist, "mysql: sku not found")
			}
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update sku: %v", err)
		}

		var ret Sku
		if err := tx.Table(s.TableName()).Where("id = ?", sku.SkuID).Find(&ret).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get sku: %v", err)
		}

		skuPriceHistory.PrevVersion = ret.HistoryVersionId
		if err := tx.Table(skuPriceHistory.TableName()).Create(skuPriceHistory).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create sku price history: %v", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (db *commodityDB) ViewSku(ctx context.Context, skuIds []*int64, pageNum int, pageSize int) ([]*model.Sku, int64, error) {
	var (
		skus  []Sku
		total int64
	)

	offset := (pageNum - 1) * pageSize
	if err := db.client.WithContext(ctx).Table((&Sku{}).TableName()).Offset(offset).Limit(pageSize).
		Where("id IN (?)", skuIds).Find(&skus).Count(&total).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errno.Errorf(errno.ServiceSkuNotExist, "mysql: sku not found")
		}
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to view sku: %v", err)
	}

	skuIDList := make([]int64, 0, len(skus))
	for _, sku := range skus {
		skuIDList = append(skuIDList, sku.Id)
	}

	var skuToSpuList []SpuToSku
	if err := db.client.WithContext(ctx).Table((&SpuToSku{}).TableName()).Where("sku_id IN (?)", skuIDList).Find(&skuToSpuList).Error; err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get sku to spu: %v", err)
	}

	var skuSaleAttrs []SkuSaleAttr
	if err := db.client.WithContext(ctx).Table((&SkuSaleAttr{}).TableName()).Where("sku_id IN (?)", skuIDList).Find(&skuSaleAttrs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errno.Errorf(errno.ServiceSkuAttrNotExist, "mysql: sku sale attr not found")
		}
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get sku sale attr: %v", err)
	}

	result := make([]*model.Sku, 0, len(skus))
	for _, sku := range skus {
		var spuID int64
		for _, skuToSpu := range skuToSpuList {
			if skuToSpu.SkuId == sku.Id {
				spuID = skuToSpu.SpuId
				break
			}
		}

		var attrValue []*model.AttrValue
		for _, attr := range skuSaleAttrs {
			if attr.SkuId == sku.Id {
				attrValue = append(attrValue, &model.AttrValue{
					SaleAttr:  attr.SaleAttr,
					SaleValue: attr.SaleValue,
				})
			}
		}

		result = append(result, &model.Sku{
			SkuID:               sku.Id,
			CreatorID:           sku.CreatorId,
			Price:               sku.Price,
			Name:                sku.Name,
			Description:         sku.Description,
			ForSale:             sku.ForSale,
			Stock:               sku.Stock,
			StyleHeadDrawingUrl: sku.StyleHeadDrawing,
			CreatedAt:           sku.CreatedAt.Unix(),
			UpdatedAt:           sku.UpdatedAt.Unix(),
			SpuID:               spuID,
			SaleAttr:            attrValue,
			HistoryID:           sku.HistoryVersionId,
			LockStock:           sku.LockStock,
		})
	}

	return result, total, nil
}

func (db *commodityDB) DeleteSku(ctx context.Context, sku *model.Sku) error {
	s := &Sku{
		Id:        sku.SkuID,
		CreatorId: sku.CreatorID,
	}

	if err := db.client.WithContext(ctx).Table(s.TableName()).Where("id = ?", sku.SkuID).Delete(s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.Errorf(errno.ServiceSkuNotExist, "mysql: sku not found")
		}
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete sku: %v", err)
	}

	return nil
}

func (db *commodityDB) ListSkuInfo(ctx context.Context, skuInfo []*model.SkuVersion, pageNum int, pageSize int) ([]*model.Sku, error) {
	var skus []Sku

	skuId := make([]int64, 0, len(skuInfo))
	for _, v := range skuInfo {
		skuId = append(skuId, v.SkuID)
	}

	offset := (pageNum - 1) * pageSize
	if err := db.client.WithContext(ctx).Table((&Sku{}).TableName()).Offset(offset).Limit(pageSize).Where("id IN (?)", skuId).Find(&skus).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to list sku info: %v", err)
	}

	skuIDList := make([]int64, 0, len(skus))
	for _, sku := range skus {
		skuIDList = append(skuIDList, sku.Id)
	}

	var skuToSpuList []SpuToSku
	if err := db.client.WithContext(ctx).Table((&SpuToSku{}).TableName()).Where("sku_id IN (?)", skuIDList).Find(&skuToSpuList).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get sku to spu: %v", err)
	}

	result := make([]*model.Sku, 0, len(skus))

	for _, sku := range skus {
		var spuID int64
		for _, skuToSpu := range skuToSpuList {
			if skuToSpu.SkuId == sku.Id {
				spuID = skuToSpu.SpuId
				break
			}
		}

		result = append(result, &model.Sku{
			SkuID:               sku.Id,
			CreatorID:           sku.CreatorId,
			Price:               sku.Price,
			Name:                sku.Name,
			ForSale:             sku.ForSale,
			LockStock:           sku.LockStock,
			StyleHeadDrawingUrl: sku.StyleHeadDrawing,
			SpuID:               spuID,
			HistoryID:           sku.HistoryVersionId,
		})
	}

	return result, nil
}

func (db *commodityDB) CreateSkuImage(ctx context.Context, skuImage *model.SkuImage) error {
	s := &SkuImages{
		Id:    skuImage.ImageID,
		Url:   skuImage.Url,
		SkuId: skuImage.SkuID,
	}

	if err := db.client.WithContext(ctx).Table(s.TableName()).Create(s).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create sku image: %v", err)
	}

	return nil
}

func (db *commodityDB) UpdateSkuImage(ctx context.Context, skuImage *model.SkuImage) error {
	s := &SkuImages{
		Id:    skuImage.ImageID,
		Url:   skuImage.Url,
		SkuId: skuImage.SkuID,
	}

	if err := db.client.WithContext(ctx).Table(s.TableName()).Where("id = ?", skuImage.ImageID).Updates(s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.Errorf(errno.ServiceSkuImageNotExist, "mysql: sku image not found")
		}
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update sku image: %v", err)
	}

	return nil
}

func (db *commodityDB) ViewSkuImage(ctx context.Context, sku *model.Sku, pageNum int, pageSize int) ([]*model.SkuImage, int64, error) {
	s := &SkuImages{
		SkuId: sku.SkuID,
	}

	var (
		Images []SkuImages
		total  int64
	)

	offset := (pageNum - 1) * pageSize
	if err := db.client.WithContext(ctx).Table(s.TableName()).Offset(offset).Limit(pageSize).
		Where("sku_id = ?", s.SkuId).Find(&Images).Count(&total).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errno.Errorf(errno.ServiceSkuNotExist, "mysql: sku not found")
		}
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to view sku image: %v", err)
	}

	result := make([]*model.SkuImage, 0, len(Images))
	for _, v := range Images {
		result = append(result, &model.SkuImage{
			ImageID:   v.Id,
			SkuID:     v.SkuId,
			Url:       v.Url,
			CreatedAt: v.CreatedAt.Unix(),
		})
	}

	return result, total, nil
}

func (db *commodityDB) DeleteSkuImage(ctx context.Context, imageId int64) error {
	s := &SkuImages{
		Id: imageId,
	}

	if err := db.client.WithContext(ctx).Table(s.TableName()).Where("id = ?", imageId).Delete(s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.Errorf(errno.ServiceSkuImageNotExist, "mysql: sku image not found")
		}
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete sku image: %v", err)
	}

	return nil
}

func (db *commodityDB) IsSpuExist(ctx context.Context, spuId int64) (bool, error) {
	var cnt int64
	if err := db.client.WithContext(ctx).Table(constants.SpuTableName).Where("id = ?", spuId).Count(&cnt).Error; err != nil {
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get count: %v", err)
	}
	return cnt != 0, nil
}

func (db *commodityDB) GetSkuBySkuId(ctx context.Context, skuId int64) (*model.Sku, error) {
	var sku Sku

	if err := db.client.WithContext(ctx).Table(sku.TableName()).Where("id = ?", skuId).First(&sku).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceSkuNotExist, "mysql: sku not found")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to view sku: %v", err)
	}
	result := &model.Sku{
		SkuID:               sku.Id,
		CreatorID:           sku.CreatorId,
		Price:               sku.Price,
		Name:                sku.Name,
		Description:         sku.Description,
		ForSale:             sku.ForSale,
		Stock:               sku.Stock,
		StyleHeadDrawingUrl: sku.StyleHeadDrawing,
		CreatedAt:           sku.CreatedAt.Unix(),
		UpdatedAt:           sku.UpdatedAt.Unix(),
		LockStock:           sku.LockStock,
	}

	return result, nil
}

func (db *commodityDB) GetSkuImageByImageId(ctx context.Context, imageId int64) (*model.SkuImage, error) {
	s := &SkuImages{
		Id: imageId,
	}

	if err := db.client.WithContext(ctx).Table(s.TableName()).Where("id = ?", s.Id).Find(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServiceSkuImageNotExist, "mysql: sku image not found")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get sku image: %v", err)
	}

	result := &model.SkuImage{
		ImageID:   s.Id,
		SkuID:     s.SkuId,
		Url:       s.Url,
		CreatedAt: s.CreatedAt.Unix(),
	}

	return result, nil
}

func (db *commodityDB) GetSkuIdBySpuID(ctx context.Context, spuId int64, pageNum int, pageSize int) ([]*int64, error) {
	i := &SpuToSku{
		SpuId: spuId,
	}

	var skuIds []SpuToSku

	offset := (pageNum - 1) * pageSize
	if err := db.client.WithContext(ctx).Table(i.TableName()).Offset(offset).Limit(pageSize).Where("spu_id = ?", i.SpuId).Find(&skuIds).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get sku by spu id: %v", err)
	}

	skuIdsList := make([]*int64, 0, len(skuIds))
	for _, v := range skuIds {
		skuIdsList = append(skuIdsList, &v.SkuId)
	}

	return skuIdsList, nil
}

func (db *commodityDB) UploadSkuAttr(ctx context.Context, sku *model.Sku, attr *model.AttrValue, id int64) error {
	s := &SkuSaleAttr{
		Id:               id,
		SkuId:            sku.SkuID,
		HistoryVersionId: sku.HistoryID,
		SaleAttr:         attr.SaleAttr,
		SaleValue:        attr.SaleValue,
	}

	if err := db.client.WithContext(ctx).Table(s.TableName()).Create(s).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to upload sku attr: %v", err)
	}

	return nil
}

func (db *commodityDB) UpdateCategory(ctx context.Context, category *model.Category) error {
	if err := db.client.WithContext(ctx).Model(&model.Category{}).Where("id = ?", category.Id).Updates(category).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update category: %v", err)
	}
	return nil
}

func (db *commodityDB) ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*model.CategoryInfo, err error) {
	offset := (pageNum - 1) * pageSize
	cs := make([]*Category, 0)
	if err := db.client.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&cs).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to list categories: %v", err)
	}
	resp = make([]*model.CategoryInfo, 0)
	for _, c := range cs {
		resp = append(resp, &model.CategoryInfo{
			Name:       c.Name,
			CategoryID: c.Id,
		})
	}
	return resp, nil
}

func (db *commodityDB) IncrLockStock(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, info := range infos {
			var lockStock int

			if err := tx.Raw("select lock_stock from "+constants.SkuTableName+
				" where id = ? for update", info.SkuID).Scan(&lockStock).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update lock stock: %v", err)
			}

			if lockStock < 0 {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: failed to increase, invalid stock num:%d ", lockStock)
			}

			if err := tx.Table(constants.SkuTableName).Where("id = ?", info.SkuID).
				UpdateColumn("lock_stock", gorm.Expr("lock_stock + ?", info.Count)).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update lock stock: %v", err)
			}
		}

		return nil
	})
	return err
}

func (db *commodityDB) DecrLockStock(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, info := range infos {
			var lockStock int

			if err := tx.Raw("SELECT lock_stock FROM "+constants.SkuTableName+
				" WHERE id = ? FOR UPDATE", info.SkuID).Scan(&lockStock).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to lock stock row: %v", err)
			}

			if lockStock <= 0 || lockStock < int(info.Count) {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: not enough locked stock to decrease (available: %d, requested: %d)", lockStock, info.Count)
			}

			if err := tx.Table(constants.SkuTableName).Where("id = ?", info.SkuID).
				UpdateColumn("lock_stock", gorm.Expr("lock_stock - ?", info.Count)).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to decrease lock stock: %v", err)
			}
		}
		return nil
	})

	return err
}

func (db *commodityDB) IncrStock(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, info := range infos {
			var stock int
			if err := tx.Raw("select stock from "+constants.SkuTableName+
				" where id = ? for update", info.SkuID).Scan(&stock).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update stock: %v", err)
			}
			if stock < 0 {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: failed to increase, invalid stock num:%d ", stock)
			}
			if err := tx.Table(constants.SkuTableName).Where("id = ?", info.SkuID).
				UpdateColumn("stock", gorm.Expr("stock + ?", info.Count)).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update stock: %v", err)
			}
		}

		return nil
	})
	return err
}

func (db *commodityDB) DecrStock(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, info := range infos {
			var s model.SkuStock
			if err := tx.Raw("SELECT stock, lock_stock FROM "+constants.SkuTableName+" WHERE id = ? FOR UPDATE", info.SkuID).Scan(&s).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to fetch sku stock info: %v", err)
			}

			stock := s.Stock
			lockStock := s.LockStock

			if stock < info.Count || stock <= 0 {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: not enough  stock to decrease (available: %d, requested: %d)", stock, info.Count)
			}

			if lockStock > stock || lockStock <= 0 || lockStock-info.Count < 0 {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: not enough  stock to decrease, lockStock = %d, stock = %d", lockStock, stock)
			}

			updates := map[string]interface{}{
				"stock":      gorm.Expr("stock - ?", info.Count),
				"lock_stock": gorm.Expr("lock_stock - ?", info.Count),
			}

			if err := tx.Table(constants.SkuTableName).Where("id = ?", info.SkuID).Updates(updates).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to decrease stock and lock stock: %v", err)
			}
		}
		return nil
	})

	return err
}

func (c *commodityDB) GetSkuById(ctx context.Context, id int64) (*model.Sku, error) {
	var s Sku
	if err := c.client.WithContext(ctx).Table(constants.SkuTableName).Where("id = ?", id).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ErrRecordNotFound, "mysql: sku not exist")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get sku: %v", err)
	}
	return &model.Sku{
		SkuID:     id,
		Stock:     s.Stock,
		LockStock: s.LockStock,
	}, nil
}

func (c *commodityDB) DecrStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := c.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, info := range infos {
			var s model.SkuStock

			if err := tx.Raw("SELECT stock, lock_stock FROM "+constants.SkuTableName+" WHERE id = ?", info.SkuID).Scan(&s).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to fetch sku stock: %v", err)
			}

			stock := s.Stock
			lockStock := s.LockStock

			if stock < info.Count || stock <= 0 {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: not enough  stock to decrease (available: %d, requested: %d)", stock, info.Count)
			}
			log.Println(s)
			if lockStock > stock || lockStock <= 0 || lockStock-info.Count < 0 {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: not enough  stock to decrease, lockStock = %d, stock = %d", lockStock, stock)
			}

			updates := map[string]interface{}{
				"stock":      gorm.Expr("stock - ?", info.Count),
				"lock_stock": gorm.Expr("lock_stock - ?", info.Count),
			}

			if err := tx.Table(constants.SkuTableName).Where("id = ?", info.SkuID).Updates(updates).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to decrease stock and lock stock: %v", err)
			}
		}
		return nil
	})
	return err
}

func (db *commodityDB) DecrLockStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, info := range infos {
			var lockStock int64

			if err := tx.Raw("SELECT lock_stock FROM "+constants.SkuTableName+" WHERE id = ?", info.SkuID).Scan(&lockStock).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to fetch sku lock stock: %v", err)
			}

			if lockStock < info.Count || lockStock <= 0 {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: not enough stock")
			}

			if err := tx.Table(constants.SkuTableName).Where("id = ?", info.SkuID).
				UpdateColumn("lock_stock", gorm.Expr("lock_stock - ?", info.Count)).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update stock: %v", err)
			}
		}
		return nil
	})
	return err
}

func (db *commodityDB) IncrLockStockInNX(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, info := range infos {
			var lockStock int64

			if err := tx.Raw("SELECT lock_stock FROM "+constants.SkuTableName+" WHERE id = ?", info.SkuID).Scan(&lockStock).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to fetch sku lock stock: %v", err)
			}

			if lockStock < 0 {
				return errno.Errorf(errno.InsufficientStockErrorCode, "mysql: not enough stock")
			}

			if err := tx.Table(constants.SkuTableName).Where("id = ?", info.SkuID).
				UpdateColumn("lock_stock", gorm.Expr("lock_stock + ?", info.Count)).Error; err != nil {
				return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update stock: %v", err)
			}
		}
		return nil
	})
	return err
}
