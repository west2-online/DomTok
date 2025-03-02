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

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

type CommodityVerifyOps func() error

func (svc *CommodityService) Verify(opts ...CommodityVerifyOps) error {
	for _, opt := range opts {
		if err := opt(); err != nil {
			return err
		}
	}
	return nil
}

func (svc *CommodityService) VerifyForSaleStatus(status int) CommodityVerifyOps {
	return func() error {
		if status > 0 && status != constants.CommodityAllowedForSale && status != constants.CommodityNotAllowedForSale {
			return errno.ParamVerifyError
		}
		return nil
	}
}

func (svc *CommodityService) VerifyCoupon(coupon *model.Coupon) CommodityVerifyOps {
	return func() error {
		if coupon.DiscountAmount > coupon.ConditionCost {
			return errno.ParamVerifyError
		}
		if coupon.Discount > 1 || coupon.Discount <= 0 {
			return errno.ParamVerifyError
		}
		if coupon.ExpireTime.After(coupon.DeadlineForGet) {
			return errno.ParamVerifyError
		}
		if len(coupon.Name) >= constants.CouponMaxVarCharLen || len(coupon.Description) >= constants.CouponMaxVarCharLen {
			return errno.ParamVerifyError
		}
		switch coupon.RangeType {
		case constants.CouponRangeTypeSPU:
			_, err := svc.db.GetSpuBySpuId(context.Background(), coupon.RangeId)
			if err != nil {
				return fmt.Errorf("check spu exist failed or non-exist: %w", err)
			}
		case constants.CouponRangeTypeCategory:
		//	e, err := svc.db.IsCategoryExistById(context.Background(), coupon.RangeId)
		//	if err != nil {
		//		return fmt.Errorf("check sku exist failed: %w", err)
		//	}
		//	if !e {
		//		return errno.ParamVerifyError
		//	}
		default:
			return errno.ParamVerifyError
		}
		return nil
	}
}

func (svc *CommodityService) VerifyPageNum(pageNum int64) CommodityVerifyOps {
	return func() error {
		if pageNum < 1 {
			return errno.ParamVerifyError
		}
		return nil
	}
}

func (svc *CommodityService) VerifyRemainUses(times int64) CommodityVerifyOps {
	return func() error {
		if times < 1 {
			return errno.ParamVerifyError
		}
		return nil
	}
}

func (svc *CommodityService) VerifyCategoryId(ctx context.Context, categoryId int64) CommodityVerifyOps {
	return func() error {
		_, err := svc.db.IsCategoryExistById(ctx, categoryId)
		if err != nil {
			return fmt.Errorf("CommodityService.VerifyCategoryId failed :%w", err)
		}
		return nil
	}
}
