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

package usecase

import (
	"context"
	"fmt"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (uc *useCase) CreateCoupon(ctx context.Context, coupon *model.Coupon) (int64, error) {
	if err := uc.svc.Verify(uc.svc.VerifyCoupon(coupon)); err != nil {
		return -1, err
	}
	err := uc.svc.InitCoupon(ctx, coupon)
	if err != nil {
		return -1, fmt.Errorf("usecase.CreateCoupon error: %w", err)
	}

	couponId, err := uc.db.CreateCoupon(ctx, coupon)
	if err != nil {
		return -1, fmt.Errorf("usecase.CreateCoupon error: %w", err)
	}
	return couponId, nil
}

func (uc *useCase) DeleteCoupon(ctx context.Context, coupon *model.Coupon) (err error) {
	uid, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("usecase.DeleteCoupon get logindata error: %w", err)
	}
	e, coupon, err := uc.db.GetCouponById(ctx, coupon.Id)
	if err != nil {
		return fmt.Errorf("usecase.DeleteCoupon error: %w", err)
	}
	if !e {
		return errno.ParamVerifyError
	}
	if uid != coupon.Uid {
		return errno.AuthInvalid
	}
	err = uc.db.DeleteCouponById(ctx, coupon)
	if err != nil {
		return fmt.Errorf("usecase.DeleteCoupon error: %w", err)
	}
	return
}

func (uc *useCase) GetCreatorCoupons(ctx context.Context, pageNum int64) (coupons []*model.Coupon, err error) {
	if err = uc.svc.Verify(uc.svc.VerifyPageNum(pageNum)); err != nil {
		return nil, err
	}
	uid, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase.CreatorGetCoupons get logindata error: %w", err)
	}
	coupons, err = uc.db.GetCouponsByCreatorId(ctx, uid, pageNum)
	if err != nil {
		return nil, fmt.Errorf("usecase.CreatorGetCoupons get coupons error: %w", err)
	}
	return
}

func (uc *useCase) CreateUserCoupon(ctx context.Context, coupon *model.UserCoupon) (err error) {
	if err = uc.svc.Verify(uc.svc.VerifyRemainUses(coupon.RemainingUses)); err != nil {
		return
	}
	uid, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("usecase.UserGetCoupons get logindata error: %w", err)
	}
	coupon.Uid = uid
	e, _, err := uc.db.GetCouponById(ctx, coupon.CouponId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteCoupon error: %w", err)
	}
	if !e {
		return errno.ParamVerifyError
	}
	err = uc.db.CreateUserCoupon(ctx, coupon)
	if err != nil {
		return fmt.Errorf("usecase.UserGetCoupon error: %w", err)
	}
	return
}

func (uc *useCase) SearchUserCoupons(ctx context.Context, pageNum int64) (coupons []*model.Coupon, err error) {
	if err = uc.svc.Verify(uc.svc.VerifyPageNum(pageNum)); err != nil {
		return nil, err
	}
	uid, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase.UserGetCoupons get logindata error: %w", err)
	}
	userCouponList, err := uc.db.GetUserCouponsByUId(ctx, uid, pageNum)
	if err != nil {
		return nil, fmt.Errorf("usecase.UserGetCoupons error: %w", err)
	}
	coupons, err = uc.svc.GetCouponsByUserCoupons(ctx, userCouponList)
	if err != nil {
		return nil, fmt.Errorf("usecase.UserGetCoupons error: %w", err)
	}
	return
}

func (uc *useCase) GetCouponAndPrice(ctx context.Context, goods []*model.OrderGoods) ([]*model.OrderGoods, float64, error) {
	return uc.svc.CalculateWithCoupon()
}
