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
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
)

func TestCommodityDB_Coupon(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	var uid int64 = 10000

	coupon := &model.Coupon{
		Id:             1,
		Uid:            uid,
		Name:           "Test Coupon",
		ConditionCost:  100,
		DiscountAmount: 10,
		Discount:       0.1,
		RangeId:        0,
		Description:    "Test description",
		ExpireTime:     time.Now(),
		DeadlineForGet: time.Now(),
	}

	Convey("TestCommodityDB_Coupon", t, func() {
		Convey("TestCommodityDB_CreateCoupon", func() {
			id, err := _db.CreateCoupon(ctx, coupon)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
		})

		Convey("TestCommodityDB_GetCouponById", func() {
			exists, ret, err := _db.GetCouponById(ctx, coupon.Id)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
			So(ret, ShouldNotBeNil)
			So(ret.Name, ShouldEqual, coupon.Name)
		})

		Convey("TestCommodityDB_GetCouponsByCreatorId", func() {
			coupons, err := _db.GetCouponsByCreatorId(ctx, uid, 1)
			So(err, ShouldBeNil)
			So(len(coupons), ShouldBeGreaterThanOrEqualTo, 1)
		})

		Convey("TestCommodityDB_DeleteCouponById", func() {
			err := _db.DeleteCouponById(ctx, coupon)
			So(err, ShouldBeNil)

			exists, _, err := _db.GetCouponById(ctx, coupon.Id)
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)
		})
	})
}

func TestCommodityDB_UserCoupon(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	var uid int64 = 10000
	var couponId int64 = 1

	userCoupon := &model.UserCoupon{
		Uid:           uid,
		CouponId:      couponId,
		RemainingUses: 5,
	}

	Convey("TestCommodityDB_UserCoupon", t, func() {
		Convey("TestCommodityDB_CreateUserCoupon", func() {
			err := _db.CreateUserCoupon(ctx, userCoupon)
			So(err, ShouldBeNil)
		})

		Convey("TestCommodityDB_GetUserCouponsByUId", func() {
			coupons, err := _db.GetUserCouponsByUId(ctx, uid, 1)
			So(err, ShouldBeNil)
			So(len(coupons), ShouldBeGreaterThanOrEqualTo, 1)
		})

		Convey("TestCommodityDB_DeleteUserCoupon", func() {
			err := _db.DeleteUserCoupon(ctx, userCoupon)
			So(err, ShouldBeNil)

			fullCoupons, err := _db.GetFullUserCouponsByUId(ctx, uid)
			So(err, ShouldBeNil)
			So(len(fullCoupons), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
