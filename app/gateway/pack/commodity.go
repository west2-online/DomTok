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
	"io/ioutil"
	"mime/multipart"

	"github.com/west2-online/DomTok/app/gateway/model/model"
	modelKitex "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/base"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/upyun"
)

func BuildFileDataBytes(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, errno.OSOperationError.WithError(err)
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, errno.IOOperationError.WithError(err)
	}
	return data, err
}

func BuildSpuImage(img *modelKitex.SpuImage) *model.SpuImage {
	return &model.SpuImage{
		ImageID:   img.ImageID,
		SpuID:     img.SpuID,
		URL:       upyun.GetImageUrl(img.Url),
		CreatedAt: img.CreatedAt,
		UpdatedAt: img.UpdatedAt,
	}
}

func BuildSpuImages(imgs []*modelKitex.SpuImage) []*model.SpuImage {
	return base.BuildTypeList(imgs, BuildSpuImage)
}

func BuildSpu(spu *modelKitex.Spu) *model.Spu {
	return &model.Spu{
		SpuID:            spu.SpuID,
		Name:             spu.Name,
		CreatorID:        spu.CreatorID,
		CreatedAt:        spu.CreatedAt,
		UpdatedAt:        spu.UpdatedAt,
		CategoryID:       spu.CategoryID,
		Description:      spu.Description,
		GoodsHeadDrawing: spu.GoodsHeadDrawing,
		Price:            spu.Price,
		ForSale:          spu.ForSale,
		Shipping:         spu.Shipping,
	}
}

func BuildSpus(spus []*modelKitex.Spu) []*model.Spu {
	return base.BuildTypeList(spus, BuildSpu)
}

func BuildCoupon(coupon *modelKitex.Coupon) *model.Coupon {
	return &model.Coupon{
		CouponID:       coupon.CouponID,
		CreatorID:      coupon.CreatorID,
		DeadlineForGet: coupon.DeadlineForGet,
		Name:           coupon.Name,
		TypeInfo:       coupon.TypeInfo,
		ConditionCost:  coupon.ConditionCost,
		DiscountAmount: coupon.DiscountAmount,
		Discount:       coupon.Discount,
		RangeType:      coupon.RangeType,
		RangeId:        coupon.RangeId,
		ExpireTime:     coupon.ExpireTime,
		Description:    coupon.Description,
	}
}

func BuildCoupons(coupons []*modelKitex.Coupon) []*model.Coupon {
	return base.BuildTypeList(coupons, BuildCoupon)
}

func BuildCategory(category *modelKitex.CategoryInfo) *model.CategoryInfo {
	return &model.CategoryInfo{
		CategoryID: category.CategoryID,
		Name:       category.Name,
	}
}

func BuildCategorys(categorys []*modelKitex.CategoryInfo) []*model.CategoryInfo {
	return base.BuildTypeList(categorys,BuildCategory)
}
