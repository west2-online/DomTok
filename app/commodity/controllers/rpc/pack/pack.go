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

package pack

import (
	"github.com/west2-online/DomTok/app/commodity/domain/model"
	modelKitex "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/base"
	"github.com/west2-online/DomTok/pkg/upyun"
)

func BuildImage(img *model.SpuImage) *modelKitex.SpuImage {
	return &modelKitex.SpuImage{
		ImageID:   img.ImageID,
		SpuID:     img.SpuID,
		Url:       img.Url,
		CreatedAt: img.CreatedAt,
		UpdatedAt: img.UpdatedAt,
	}
}

func BuildImages(imgs []*model.SpuImage) []*modelKitex.SpuImage {
	return base.BuildTypeList(imgs, BuildImage)
}

func BuildSpu(spu *model.Spu) *modelKitex.Spu {
	return &modelKitex.Spu{
		SpuID:            spu.SpuId,
		Name:             spu.Name,
		CreatorID:        spu.CreatorId,
		CategoryID:       spu.CategoryId,
		Description:      spu.Description,
		GoodsHeadDrawing: upyun.GetImageUrl(spu.GoodsHeadDrawingUrl),
		Price:            spu.Price,
		ForSale:          int32(spu.ForSale),
		Shipping:         spu.Shipping,
		CreatedAt:        spu.CreatedAt,
		UpdatedAt:        spu.UpdatedAt,
	}
}

func BuildSpus(spus []*model.Spu) []*modelKitex.Spu {
	return base.BuildTypeList(spus, BuildSpu)
}

func BuildSkuImages(i []*model.SkuImage) []*modelKitex.SkuImage {
	result := make([]*modelKitex.SkuImage, 0, len(i))
	for _, v := range i {
		result = append(result, &modelKitex.SkuImage{
			ImageID:   v.ImageID,
			SkuID:     v.SkuID,
			Url:       v.Url,
			CreatedAt: v.CreatedAt,
			DeletedAt: &v.DeletedAt,
		})
	}
	return result
}

func BuildSkus(i []*model.Sku) []*modelKitex.Sku {
	result := make([]*modelKitex.Sku, 0, len(i))
	for _, v := range i {
		attr := make([]*modelKitex.AttrValue, 0, len(v.SaleAttr))
		for _, value := range v.SaleAttr {
			attr = append(attr, &modelKitex.AttrValue{
				SaleAttr:  value.SaleAttr,
				SaleValue: value.SaleValue,
			})
		}

		result = append(result, &modelKitex.Sku{
			SkuID:            v.SkuID,
			CreatorID:        v.CreatorID,
			Price:            v.Price,
			Name:             v.Name,
			Description:      v.Description,
			ForSale:          int32(v.ForSale),
			Stock:            v.Stock,
			StyleHeadDrawing: v.StyleHeadDrawingUrl,
			CreatedAt:        v.CreatedAt,
			UpdatedAt:        v.UpdatedAt,
			DeletedAt:        &v.DeletedAt,
			SpuID:            v.SpuID,
			SaleAttr:         attr,
			HistoryID:        v.HistoryID,
			LockStock:        v.LockStock,
		})
	}
	return result
}

func BuildSkuInfo(s *model.Sku) *modelKitex.SkuInfo {
	var sku model.Sku
	if s == nil {
		sku.SkuID = -1
		sku.CreatorID = -1
		sku.HistoryID = -1
	} else {
		sku.SkuID = s.SkuID
		sku.CreatorID = s.CreatorID
		sku.HistoryID = s.HistoryID
	}
	result := &modelKitex.SkuInfo{
		SkuID:     sku.SkuID,
		CreatorID: sku.CreatorID,
		HistoryID: sku.HistoryID,
	}
	return result
}

func BuildSkuInfos(i []*model.Sku) []*modelKitex.SkuInfo {
	result := make([]*modelKitex.SkuInfo, 0, len(i)) // 预分配容量
	for _, v := range i {
		result = append(result, &modelKitex.SkuInfo{
			SkuID:            v.SkuID,
			CreatorID:        v.CreatorID,
			Price:            v.Price,
			Name:             v.Name,
			ForSale:          int32(v.ForSale),
			LockStock:        v.LockStock,
			StyleHeadDrawing: v.StyleHeadDrawingUrl,
			SpuID:            v.SpuID,
			HistoryID:        v.HistoryID,
		})
	}
	return result
}

func BuildSkuPriceHistory(i []*model.SkuPriceHistory) []*modelKitex.PriceHistory {
	result := make([]*modelKitex.PriceHistory, 0, len(i))
	for _, v := range i {
		result = append(result, &modelKitex.PriceHistory{
			HistoryID: v.Id,
			SkuID:     v.SkuId,
			Price:     int64(v.MarkPrice),
			CreatedAt: v.CreatedAt,
		})
	}
	return result
}

func BuildCategory(category *model.CategoryInfo) *modelKitex.CategoryInfo {
	return &modelKitex.CategoryInfo{
		CategoryID: category.CategoryID,
		Name:       category.Name,
	}
}

func BuildCategorys(categorys []*model.CategoryInfo) []*modelKitex.CategoryInfo {
	return base.BuildTypeList(categorys, BuildCategory)
}
