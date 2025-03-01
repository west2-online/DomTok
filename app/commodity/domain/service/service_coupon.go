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

package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (svc *CommodityService) InitCoupon(ctx context.Context, coupon *model.Coupon) error {
	uid, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("service.CreateCoupon get logindata error: %w", err)
	}
	coupon.Uid = uid
	coupon.Id = svc.nextID()
	return nil
}

func (svc *CommodityService) GetCouponsByUserCoupons(ctx context.Context, userCouponList []*model.UserCoupon) (couponList []*model.Coupon, err error) {
	couponIDs := make([]int64, 0, len(userCouponList))

	for _, userCoupon := range userCouponList {
		couponIDs = append(couponIDs, userCoupon.CouponId)
	}

	couponList, err = svc.db.GetCouponsByIDs(ctx, couponIDs)
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "service: failed to get coupons: %v", err)
	}

	return couponList, nil
}

func (svc *CommodityService) CalculateWithCoupon(ctx context.Context, goods []*model.OrderGoods) ([]*model.OrderGoods, float64, error) {
	uid, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return nil, -1, fmt.Errorf("svc.GetCouponByCommoditie get logindata error: %w", err)
	}
	userCoupons, err := svc.db.GetFullUserCouponsByUId(ctx, uid)
	if err != nil {
		return nil, -1, errno.Errorf(errno.InternalDatabaseErrorCode, "service: failed to get coupons: %v", err)
	}
	couponList, err := svc.GetCouponsByUserCoupons(ctx, userCoupons)
	if err != nil {
		return nil, -1, fmt.Errorf("svc.GetCouponByCommodities GetCouponsByUserCoupons error: %w", err)
	}
	// 直接在原切片上通过双指针修改，减少内存开销
	validPointer := 0
	for i := 0; i < len(couponList); i++ {
		// todo: 时间可能要统一一下
		if time.Now().Before(couponList[i].ExpireTime) {
			couponList[validPointer] = couponList[i]
			validPointer++
		}
	}
	couponList = couponList[:validPointer]

	/*
		思路：使用排序+双指针来处理查找过程，缺点是空间复杂度较大
		时间复杂度：O(NlogN+MlogM+N+M)
		空间复杂度：O(N+M)
	*/
	// 按 RangeType 分组
	var couponsForSpu, couponsForCategory []*model.Coupon
	for _, c := range couponList {
		switch c.RangeType {
		case constants.CouponRangeTypeSPU: // 按 SpuId 匹配
			couponsForSpu = append(couponsForSpu, c)
		case constants.CouponRangeTypeCategory: // 按 CategoryId 匹配
			couponsForCategory = append(couponsForCategory, c)
		default:
			continue
		}
	}

	// 升序排序
	sort.Slice(couponsForSpu, func(i, j int) bool {
		return couponsForSpu[i].RangeId < couponsForSpu[j].RangeId
	})
	sort.Slice(couponsForCategory, func(i, j int) bool {
		return couponsForCategory[i].RangeId < couponsForCategory[j].RangeId
	})

	matchMap := make(map[int64][]*model.Coupon)

	// 将 spuList 按 SpuId 做升序排序
	sort.Slice(goods, func(i, j int) bool {
		return goods[i].GoodsId < goods[j].GoodsId
	})
	// 构造双指针
	i, j := 0, 0
	for i < len(goods) && j < len(couponsForSpu) {
		spuId := goods[i].GoodsId
		couponRangeId := couponsForSpu[j].RangeId
		switch {
		case spuId == couponRangeId:
			matchMap[spuId] = append(matchMap[spuId], couponsForSpu[j])
			j++ // 优惠券移动
		case spuId < couponRangeId:
			i++ // spuId移动
		default:
			j++
		}
	}

	/*
		// 同理处理CategoryId
		sort.Slice(spuList, func(i, j int) bool {
			return spuList[i].CategoryId < spuList[j].CategoryId
		})
		i, j = 0, 0
		for i < len(spuList) && j < len(couponsForCategory) {
			categoryId := spuList[i].CategoryId
			couponRangeId := couponsForCategory[j].RangeId

			switch {
			case categoryId == couponRangeId:
				spuId := spuList[i].SpuId
				matchMap[spuId] = append(matchMap[spuId], couponsForCategory[j])
				j++
			case categoryId < couponRangeId:
				i++
			default:
				j++
			}
		}

	*/
	goodsResult, totalPrice := svc.assignCoupons(goods, matchMap)
	return goodsResult, totalPrice, nil
}

// assignCouponsByPrice 以商品价格降序为优先级，从 matchMap 中给每个 SPU 匹配优惠券
func (svc *CommodityService) assignCoupons(goodsList []*model.OrderGoods,
	matchMap map[int64][]*model.Coupon,
) ([]*model.OrderGoods, float64) {
	// 按 spu.Price 进行降序排序，让价格最高的商品优先匹配
	sort.Slice(goodsList, func(i, j int) bool {
		return goodsList[i].TotalAmount > goodsList[j].TotalAmount
	})

	// 用于标记本次交易里某张优惠券是否已被占用
	usedCoupons := make(map[int64]bool)

	var totalPrice float64

	for _, spu := range goodsList {
		bestPrice := spu.TotalAmount
		var bestCoupon *model.Coupon

		couponCandidates, ok := matchMap[spu.GoodsId]
		// 没有可用券
		if !ok || len(couponCandidates) == 0 {
			totalPrice += bestPrice + spu.SingleFreightPrice
			spu.DiscountAmount = bestPrice + spu.SingleFreightPrice
			continue
		}

		// 遍历所有匹配到的优惠券，找到最优
		for _, c := range couponCandidates {
			// 如果已经被其他 SPU 占用，则跳过
			if usedCoupons[c.Id] {
				continue
			}

			// 判断是否可用
			canUse := false
			if spu.TotalAmount >= c.ConditionCost {
				canUse = true
			}

			// 如果当前券可用，计算折后价
			if canUse {
				discountedPrice := c.CalculateDiscountPrice(spu.TotalAmount)
				// 更优解，替换
				if discountedPrice < bestPrice {
					bestPrice = discountedPrice
					bestCoupon = c
				}
			}
		}

		if bestCoupon != nil {
			spu.CouponId = bestCoupon.Id
			spu.CouponName = bestCoupon.Name
			usedCoupons[bestCoupon.Id] = true
		}

		totalPrice += bestPrice + spu.SingleFreightPrice
		spu.DiscountAmount = bestPrice + spu.SingleFreightPrice
	}

	return goodsList, totalPrice
}
