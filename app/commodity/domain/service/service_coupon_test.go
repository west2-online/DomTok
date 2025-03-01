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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/infrastructure/mysql"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

// TestCommodityService_CalculateWithCoupon 为 CalculateWithCoupon 接口编写测试
func TestCommodityService_CalculateWithCoupon(t *testing.T) {
	// 定义测试用例结构体
	type TestCase struct {
		Name string
		// 模拟登录数据
		MockLoginUID   int64
		MockLoginError error
		// 模拟 db.GetFullUserCouponsByUId 返回
		MockUserCoupons             []*model.UserCoupon
		MockGetFullUserCouponsError error
		// 模拟 svc.GetCouponsByUserCoupons 返回（这里返回的优惠券列表中可能包含过期优惠券）
		MockCouponList                   []*model.Coupon
		MockGetCouponsByUserCouponsError error
		// 输入的 spu 列表
		SpuList []*model.Spu
		// 模拟 assignCoupons 返回结果
		MockAssignedMap map[int64]*model.Coupon
		MockPriceMap    map[int64]float64
		MockTotalPrice  float64
		// 模拟 ConvertMapsToAssignedCoupon 返回结果
		ExpectedAssignedCoupons []*model.AssignedCoupon
		// 预期的错误
		ExpectedError error
		// 预期总价格（仅在无错误时校验）
		ExpectedTotalPrice float64
	}

	// 构造测试用例
	testCases := []TestCase{
		{
			Name:           "登录失败",
			MockLoginError: errors.New("登录错误"),
			ExpectedError:  fmt.Errorf("svc.GetCouponByCommoditie get logindata error: %w", errors.New("登录错误")),
		},
		{
			Name:                        "获取用户优惠券失败",
			MockLoginUID:                101,
			MockUserCoupons:             nil,
			MockGetFullUserCouponsError: errors.New("数据库错误"),
			ExpectedError:               errno.Errorf(errno.InternalDatabaseErrorCode, "service: failed to get coupons: %v", errors.New("数据库错误")),
		},
		{
			Name:                             "获取优惠券列表失败",
			MockLoginUID:                     101,
			MockUserCoupons:                  []*model.UserCoupon{{CouponId: 1}},
			MockGetFullUserCouponsError:      nil,
			MockGetCouponsByUserCouponsError: errors.New("服务错误"),
			ExpectedError:                    fmt.Errorf("svc.GetCouponByCommodities GetCouponsByUserCoupons error: %w", errors.New("服务错误")),
		},
		{
			Name:            "成功匹配优惠券",
			MockLoginUID:    101,
			MockLoginError:  nil,
			MockUserCoupons: []*model.UserCoupon{{CouponId: 1}},
			// 返回的优惠券列表中包含两个优惠券，其中一个未过期，一个已过期（已过期的会被过滤掉）
			MockCouponList: []*model.Coupon{
				{
					Id:             1,
					RangeType:      constants.CouponRangeTypeSPU, // 假设此常量在项目中已定义
					RangeId:        1001,
					ExpireTime:     time.Now().Add(1 * time.Hour),
					ConditionCost:  50,
					DiscountAmount: 10,
				},
				{
					Id:             2,
					RangeType:      constants.CouponRangeTypeCategory,
					RangeId:        2001,
					ExpireTime:     time.Now().Add(-1 * time.Hour),
					ConditionCost:  30,
					DiscountAmount: 5,
				},
			},
			MockGetCouponsByUserCouponsError: nil,
			// 输入的商品列表：一个 SPUId 为 1001、CategoryId 随意；另一个 SPUId 为 1002
			SpuList: []*model.Spu{
				{SpuId: 1001, CategoryId: 3001, Price: 60},
				{SpuId: 1002, CategoryId: 2001, Price: 40},
			},
			// 模拟 assignCoupons 返回，假设只有第一个商品匹配到优惠券，计算折后价格为 55，而第二个商品无优惠
			MockAssignedMap: map[int64]*model.Coupon{
				1001: {Id: 1, RangeType: constants.CouponRangeTypeSPU, RangeId: 1001, ConditionCost: 50, DiscountAmount: 10},
			},
			MockPriceMap: map[int64]float64{
				1001: 55,
				1002: 40,
			},
			MockTotalPrice: 95,
			ExpectedAssignedCoupons: []*model.AssignedCoupon{
				{SpuId: 1001, Coupon: &model.Coupon{
					Id: 1, RangeType: constants.CouponRangeTypeSPU,
					RangeId: 1001, ConditionCost: 50, DiscountAmount: 10,
				}, DiscountedPrice: 55},
			},
			ExpectedError:      nil,
			ExpectedTotalPrice: 95,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 模拟获取登录数据
			mockey.Mock(contextLogin.GetLoginData).Return(tc.MockLoginUID, tc.MockLoginError).Build()

			// 构造一个 db 对象并模拟 GetFullUserCouponsByUId
			db := mysql.NewCommodityDB(new(gorm.DB))
			mockey.Mock(mockey.GetMethod(db, "GetFullUserCouponsByUId")).Return(tc.MockUserCoupons, tc.MockGetFullUserCouponsError).Build()

			// 创建一个 CommodityService 对象，并将 db 注入
			svc := new(CommodityService)
			svc.db = db

			// 模拟 svc.GetCouponsByUserCoupons 方法返回优惠券列表
			mockey.Mock((*CommodityService).GetCouponsByUserCoupons).Return(tc.MockCouponList, tc.MockGetCouponsByUserCouponsError).Build()

			// 模拟 assignCoupons 方法，返回预设的 assignedMap、priceMap 和 totalPrice
			mockey.Mock((*CommodityService).assignCoupons).Return(tc.MockAssignedMap, tc.MockPriceMap, tc.MockTotalPrice).Build()

			// 模拟 ConvertMapsToAssignedCoupon 函数，返回预期的 AssignedCoupon 列表
			mockey.Mock(model.ConvertMapsToAssignedCoupon).Return(tc.ExpectedAssignedCoupons).Build()

			// 调用 CalculateWithCoupon 方法
			assignedCoupons, totalPrice, err := svc.CalculateWithCoupon(context.Background(), tc.SpuList)
			if err != nil && tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}

			// 如果没有错误，再对返回的优惠券列表和总价格进行校验
			if err == nil {
				convey.So(assignedCoupons, convey.ShouldResemble, tc.ExpectedAssignedCoupons)
				convey.So(totalPrice, convey.ShouldEqual, tc.ExpectedTotalPrice)
			}
		})
	}
}
