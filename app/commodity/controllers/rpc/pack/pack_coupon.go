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
)

func BuildCoupon(coupon *model.Coupon) *modelKitex.Coupon {
	return &modelKitex.Coupon{
		CouponID:       coupon.Id,
		CreatorID:      coupon.Uid,
		DeadlineForGet: coupon.DeadlineForGet.Unix(),
		Name:           coupon.Name,
		TypeInfo:       int32(coupon.TypeInfo),
		ConditionCost:  coupon.ConditionCost,
		DiscountAmount: &coupon.DiscountAmount,
		Discount:       &coupon.Discount,
		RangeType:      int32(coupon.RangeType),
		RangeId:        coupon.RangeId,
		ExpireTime:     coupon.ExpireTime.Unix(),
		Description:    coupon.Description,
	}
}

func BuildCoupons(coupons []*model.Coupon) []*modelKitex.Coupon {
	return base.BuildTypeList(coupons, BuildCoupon)
}

func ConvertOrderGoods(goods *modelKitex.OrderGoods) *model.OrderGoods {
	return &model.OrderGoods{
		MerchantId:         goods.MerchantId,
		GoodsId:            goods.GoodsId,
		GoodsName:          goods.GoodsName,
		StyleId:            goods.StyleId,
		StyleName:          goods.StyleName,
		GoodsVersion:       goods.GoodsVersion,
		StyleHeadDrawing:   goods.StyleHeadDrawing,
		OriginPrice:        goods.OriginPrice,
		SalePrice:          goods.SalePrice,
		SingleFreightPrice: goods.SingleFreightPrice,
		PurchaseQuantity:   goods.PurchaseQuantity,
		TotalAmount:        goods.TotalAmount,
		FreightAmount:      goods.FreightAmount,
	}
}

func ConvertOrderGoodsList(goods []*modelKitex.OrderGoods) []*model.OrderGoods {
	return base.BuildTypeList(goods, ConvertOrderGoods)
}

func BuildOrderGoods(goods *model.OrderGoods) *modelKitex.OrderGoods {
	return &modelKitex.OrderGoods{
		MerchantId:         goods.MerchantId,
		GoodsId:            goods.GoodsId,
		GoodsName:          goods.GoodsName,
		StyleId:            goods.StyleId,
		StyleName:          goods.StyleName,
		GoodsVersion:       goods.GoodsVersion,
		StyleHeadDrawing:   goods.StyleHeadDrawing,
		OriginPrice:        goods.OriginPrice,
		SalePrice:          goods.SalePrice,
		SingleFreightPrice: goods.SingleFreightPrice,
		PurchaseQuantity:   goods.PurchaseQuantity,
		TotalAmount:        goods.TotalAmount,
		FreightAmount:      goods.FreightAmount,
		DiscountAmount:     goods.DiscountAmount,
		PaymentAmount:      goods.PaymentAmount,
		SinglePrice:        goods.SinglePrice,
		CouponId:           goods.CouponId,
		CouponName:         goods.CouponName,
	}
}

func BuildOrderGoodsList(goods []*model.OrderGoods) []*modelKitex.OrderGoods {
	return base.BuildTypeList(goods, BuildOrderGoods)
}
