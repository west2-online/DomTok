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

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (db *commodityDB) CreateCoupon(ctx context.Context, coupon *model.Coupon) error {
	dbModel := &Coupon{
		Id:             coupon.Id,
		Uid:            coupon.Uid,
		Name:           coupon.Name,
		TypeInfo:       coupon.TypeInfo,
		ConditionCost:  coupon.ConditionCost,
		DiscountAmount: coupon.DiscountAmount,
		Discount:       coupon.Discount,
		RangeType:      coupon.RangeType,
		RangeId:        coupon.RangeId,
		Description:    coupon.Description,
		ExpireTime:     coupon.ExpireTime,
		DeadlineForGet: coupon.DeadlineForGet,
	}
	if err := db.client.WithContext(ctx).Create(dbModel).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create coupon: %v", err)
	}
	return nil
}

func (db *commodityDB) GetCouponById(ctx context.Context, id int64) (*model.Coupon, error) {
	dbModel := &Coupon{
		Id: id,
	}
	if err := db.client.WithContext(ctx).First(dbModel).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to find coupon: %v", err)
	}
	return &model.Coupon{
		Id:             dbModel.Id,
		Uid:            dbModel.Uid,
		Name:           dbModel.Name,
		TypeInfo:       dbModel.TypeInfo,
		ConditionCost:  dbModel.ConditionCost,
		DiscountAmount: dbModel.DiscountAmount,
		Discount:       dbModel.Discount,
		RangeType:      dbModel.RangeType,
		RangeId:        dbModel.RangeId,
		Description:    dbModel.Description,
		ExpireTime:     dbModel.ExpireTime,
		DeadlineForGet: dbModel.DeadlineForGet,
	}, nil
}

func (db *commodityDB) CheckExistCouponByCreatorId(ctx context.Context, uid int64) ([]*model.Coupon, error) {
	dbModel := make([]*Coupon, 0)
	if err := db.client.WithContext(ctx).Where("uid = ?", uid).Find(&dbModel).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to find coupon: %v", err)
	}
	result := make([]*model.Coupon, 0)
	for _, coupon := range dbModel {
		result = append(result, &model.Coupon{
			Id:             coupon.Id,
			Uid:            coupon.Uid,
			Name:           coupon.Name,
			TypeInfo:       coupon.TypeInfo,
			ConditionCost:  coupon.ConditionCost,
			DiscountAmount: coupon.DiscountAmount,
			Discount:       coupon.Discount,
			RangeType:      coupon.RangeType,
			RangeId:        coupon.RangeId,
			Description:    coupon.Description,
			ExpireTime:     coupon.ExpireTime,
			DeadlineForGet: coupon.DeadlineForGet,
		})
	}
	return result, nil
}

func (db *commodityDB) DeleteCouponById(ctx context.Context, id int64) error {
	dbModel := &Coupon{
		Id: id,
	}
	if err := db.client.WithContext(ctx).Delete(dbModel).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete coupon: %v", err)
	}
	return nil
}

func (db *commodityDB) CreateUserCoupon(ctx context.Context, coupon *model.UserCoupon) error {
	dbModel := &UserCoupon{
		Uid:           coupon.Uid,
		CouponId:      coupon.CouponId,
		RemainingUses: coupon.RemainingUses,
	}
	if err := db.client.WithContext(ctx).Create(dbModel).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create user coupon: %v", err)
	}
	return nil
}

func (db *commodityDB) GetUserCouponByUId(ctx context.Context, uid int64) ([]*model.UserCoupon, error) {
	dbModel := make([]*UserCoupon, 0)
	if err := db.client.WithContext(ctx).Where("uid = ?", uid).Find(&dbModel).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to find user coupon: %v", err)
	}
	result := make([]*model.UserCoupon, 0)
	for _, coupon := range dbModel {
		result = append(result, &model.UserCoupon{
			Uid:           coupon.Uid,
			CouponId:      coupon.CouponId,
			RemainingUses: coupon.RemainingUses,
		})
	}
	return result, nil
}

func (db *commodityDB) DeleteUserCoupon(ctx context.Context, coupon *model.UserCoupon) error {
	dbModel := &UserCoupon{
		Uid:      coupon.Uid,
		CouponId: coupon.CouponId,
	}
	if err := db.client.WithContext(ctx).Delete(dbModel).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete user coupon: %v", err)
	}
	return nil
}
