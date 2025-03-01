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

package model

import (
	"time"

	"github.com/west2-online/DomTok/pkg/constants"
)

type Coupon struct {
	Id             int64
	Uid            int64
	Name           string
	TypeInfo       int64
	ConditionCost  float64
	DiscountAmount float64
	Discount       float64
	RangeType      int64
	RangeId        int64
	Description    string
	ExpireTime     time.Time
	DeadlineForGet time.Time
}

type UserCoupon struct {
	Uid           int64
	CouponId      int64
	RemainingUses int64
}

type OrderGoods struct {
	GoodsId            int64
	OriginPrice        float64
	SalePrice          float64
	SingleFreightPrice float64
	TotalAmount        float64
	FreightAmount      float64
	DiscountAmount     float64
	SinglePrice        float64
	CouponId           int64
	CouponName         string
}

func (c *Coupon) CalculateDiscountPrice(originalPrice float64) float64 {
	switch c.TypeInfo {
	case constants.CouponRangeTypeSPU:
		// 减价格
		return originalPrice - c.DiscountAmount
	case constants.CouponRangeTypeCategory:
		// 打折
		return originalPrice * c.Discount
	default:
		return originalPrice
	}
}
