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

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

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
