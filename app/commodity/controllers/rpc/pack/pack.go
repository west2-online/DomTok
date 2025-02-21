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
	model2 "github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/kitex_gen/model"
)

func BuildImages(i []*model2.SkuImage) []*model.SkuImage {
	result := make([]*model.SkuImage, 0, len(i)) // 预分配容量
	for _, v := range i {
		result = append(result, &model.SkuImage{
			ImageID:   v.ImageID,
			SkuID:     v.SkuID,
			Url:       v.Url,
			CreatedAt: v.CreatedAt,
			DeletedAt: &v.DeletedAt,
		})
	}
	return result
}

func BuildSkus(i []*model2.Sku) []*model.Sku {
	result := make([]*model.Sku, 0, len(i)) // 预分配容量
	for _, v := range i {
		attr := make([]*model.AttrValue, 0, len(v.SaleAttr))
		for _, value := range v.SaleAttr {
			attr = append(attr, &model.AttrValue{
				SaleAttr:  value.SaleAttr,
				SaleValue: value.SaleValue,
			})
		}

		result = append(result, &model.Sku{
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

func BuildSkuInfos(i []*model2.Sku) []*model.SkuInfo {
	result := make([]*model.SkuInfo, 0, len(i)) // 预分配容量
	for _, v := range i {
		result = append(result, &model.SkuInfo{
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
