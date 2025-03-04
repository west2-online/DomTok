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

// TestCommodityService_CalculateWithCoupon 用中文示例:
// 1. 使用表格驱动(table‐driven)测试
// 2. 使用 mockey 对依赖进行打桩
// 3. 如果出错，只校验错误；如果成功，则校验返回值
func TestCommodityService_CalculateWithCoupon(t *testing.T) {
	type TestCase struct {
		Name string
		// 获取登录信息
		MockLoginUID   int64
		MockLoginError error
		// 模拟 GetFullUserCouponsByUId
		MockGetFullUserCoupons      []*model.UserCoupon
		MockGetFullUserCouponsError error
		// 模拟 GetCouponsByUserCoupons
		MockCoupons                      []*model.Coupon
		MockGetCouponsByUserCouponsError error
		// 传入的商品数据
		OrderGoodsList []*model.OrderGoods
		// 预期的返回错误
		ExpectedError error

		// 如果成功(没有错误)时，需要断言的返回值
		// 注意：因为你提到错误时会返回 totalPrice=-1，这里只在无错误时断言
		// 对返回的 goods 进行断言
		ExpectedGoodsResult []*model.OrderGoods
		ExpectedTotalPrice  float64
	}

	testCases := []TestCase{
		{
			Name:           "登录失败",
			MockLoginError: errors.New("登录错误"),
			// 期望最终返回的 error
			ExpectedError: fmt.Errorf("svc.GetCouponByCommoditie get logindata error: %w", errors.New("登录错误")),
		},
		{
			Name:                        "获取用户优惠券失败",
			MockLoginUID:                101,
			MockGetFullUserCouponsError: errors.New("数据库错误"),
			ExpectedError: errno.Errorf(
				errno.InternalDatabaseErrorCode,
				"service: failed to get coupons: %v",
				errors.New("数据库错误"),
			),
		},
		{
			Name:                             "获取优惠券信息失败",
			MockLoginUID:                     101,
			MockGetFullUserCoupons:           []*model.UserCoupon{{CouponId: 1}},
			MockGetCouponsByUserCouponsError: errors.New("GetCouponsByUserCoupons error"),
			ExpectedError: fmt.Errorf(
				"svc.GetCouponByCommodities GetCouponsByUserCoupons error: %w",
				errors.New("GetCouponsByUserCoupons error"),
			),
		},
		{
			Name:         "成功匹配优惠券",
			MockLoginUID: 101,
			MockGetFullUserCoupons: []*model.UserCoupon{
				{CouponId: 1},
			},
			// 包含 2 张券：1 张未过期，1 张已过期
			MockCoupons: []*model.Coupon{
				{
					Id:             1,
					TypeInfo:       constants.CouponTypeSubAmount,
					RangeType:      constants.CouponRangeTypeSPU,
					RangeId:        1001,
					ExpireTime:     time.Now().Add(1 * time.Hour), // 未过期
					ConditionCost:  10,
					DiscountAmount: 5,
				},
				{
					Id:             2,
					TypeInfo:       constants.CouponTypeSubAmount,
					RangeType:      constants.CouponRangeTypeSPU,
					RangeId:        1002,
					ExpireTime:     time.Now().Add(-1 * time.Hour), // 已过期
					ConditionCost:  20,
					DiscountAmount: 10,
				},
			},
			OrderGoodsList: []*model.OrderGoods{
				// 其中第一个 GoodsId 正好和 RangeId=1001 匹配
				{GoodsId: 1001, TotalAmount: 30, FreightAmount: 5, PurchaseQuantity: 1},
				// 第二个 GoodsId=1002 对应的优惠券已过期
				{GoodsId: 1002, TotalAmount: 40, FreightAmount: 5, PurchaseQuantity: 1},
			},
			// 无错误
			ExpectedError: nil,
			// 返回时，假设第一件商品 (GoodsId=1001) 被优惠 5 块，所以优惠后是 25 + 运费5 = 30
			// 第二件商品 (GoodsId=1002) 没有优惠券可用或可用券已过期，所以最终=40 + 5 =45
			// 因此 totalPrice=30+45=75
			ExpectedGoodsResult: []*model.OrderGoods{
				{
					GoodsId:          1002,
					CouponId:         0,
					CouponName:       "",
					TotalAmount:      40,
					FreightAmount:    5,
					DiscountAmount:   45,
					PurchaseQuantity: 1,
					SinglePrice:      45,
				},
				{
					GoodsId:          1001,
					CouponId:         1,
					CouponName:       "", // 下面会 mock assignCouponsAndPrice 设置，如果需要可进一步模拟
					TotalAmount:      30,
					FreightAmount:    5,
					DiscountAmount:   30, // 优惠后+运费
					PurchaseQuantity: 1,
					SinglePrice:      30, // 此处 = DiscountAmount / PurchaseQuantity
				},
			},
			ExpectedTotalPrice: 75,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 1. mock 登录信息
			mockey.Mock(contextLogin.GetLoginData).
				Return(tc.MockLoginUID, tc.MockLoginError).
				Build()

			// 2. mock 数据库操作
			//    这里先用一个假的 gormDB 生成 db
			db := mysql.NewCommodityDB(new(gorm.DB))
			// 模拟 GetFullUserCouponsByUId
			mockey.
				Mock(mockey.GetMethod(db, "GetFullUserCouponsByUId")).
				Return(tc.MockGetFullUserCoupons, tc.MockGetFullUserCouponsError).
				Build()

			// 3. 新建 CommodityService 并替换其 db
			svc := &CommodityService{
				db: db,
			}

			// mock svc.GetCouponsByUserCoupons
			mockey.
				Mock((*CommodityService).GetCouponsByUserCoupons).
				Return(tc.MockCoupons, tc.MockGetCouponsByUserCouponsError).
				Build()

			// 注意：如果 assignCouponsAndPrice 也需要 mock，你可以在这里再进行处理
			// 比如你想完全控制最后返回的 goods 和 totalPrice，你可以 mock:
			// mockey.Mock((*CommodityService).assignCouponsAndPrice).
			//     Return(tc.ExpectedGoodsResult, tc.ExpectedTotalPrice).
			//     Build()
			// 不过上面你也可以直接依赖真实逻辑，以测试真实行为。

			// 调用目标方法
			goodsResult, totalPrice, err := svc.CalculateWithCoupon(context.Background(), tc.OrderGoodsList)

			// 如果出错，就只判断错误，不断言返回值
			if err != nil || tc.ExpectedError != nil {
				if err != nil && tc.ExpectedError != nil {
					convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				} else {
					convey.So(err, convey.ShouldEqual, tc.ExpectedError)
				}
				return
			}

			// 没有错误时，断言结果
			convey.So(goodsResult, convey.ShouldResemble, tc.ExpectedGoodsResult)
			convey.So(totalPrice, convey.ShouldEqual, tc.ExpectedTotalPrice)
		})
	}
}
