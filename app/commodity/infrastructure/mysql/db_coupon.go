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
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (db *commodityDB) CreateCoupon(ctx context.Context, coupon *model.Coupon) (int64, error) {
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
		return -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create coupon: %v", err)
	}
	return dbModel.Id, nil
}

func (db *commodityDB) GetCouponById(ctx context.Context, id int64) (bool, *model.Coupon, error) {
	dbModel := &Coupon{
		Id: id,
	}
	if err := db.client.WithContext(ctx).First(dbModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get spu: %v", err)
	}
	return true, &model.Coupon{
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

func (db *commodityDB) GetCouponsByCreatorId(ctx context.Context, uid int64, pageNum int64) ([]*model.Coupon, error) {
	dbModel := make([]*Coupon, 0)
	offset := (pageNum - 1) * constants.CouponPageSize
	if err := db.client.WithContext(ctx).Where("uid = ?", uid).Offset(int(offset)).Limit(constants.CouponPageSize).Find(&dbModel).Error; err != nil {
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

func (db *commodityDB) DeleteCouponById(ctx context.Context, coupon *model.Coupon) error {
	dbModel := &Coupon{
		Id: coupon.Id,
	}
	if err := db.client.WithContext(ctx).Where("id = ? AND uid = ?", coupon.Id, coupon.Uid).Delete(dbModel).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete coupon: %v", err)
	}
	return nil
}

func (db *commodityDB) GetCouponsByIDs(ctx context.Context, couponIDs []int64) ([]*model.Coupon, error) {
	dbModels := make([]*Coupon, 0)

	if err := db.client.WithContext(ctx).
		Where("id IN ?", couponIDs).
		Find(&dbModels).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to find coupons: %v", err)
	}
	couponList := make([]*model.Coupon, 0, len(dbModels))
	for _, dbCoupon := range dbModels {
		couponList = append(couponList, &model.Coupon{
			Id:             dbCoupon.Id,
			Uid:            dbCoupon.Uid,
			Name:           dbCoupon.Name,
			TypeInfo:       dbCoupon.TypeInfo,
			ConditionCost:  dbCoupon.ConditionCost,
			DiscountAmount: dbCoupon.DiscountAmount,
			Discount:       dbCoupon.Discount,
			RangeType:      dbCoupon.RangeType,
			RangeId:        dbCoupon.RangeId,
			Description:    dbCoupon.Description,
			ExpireTime:     dbCoupon.ExpireTime,
			DeadlineForGet: dbCoupon.DeadlineForGet,
		})
	}

	return couponList, nil
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

func (db *commodityDB) GetUserCouponsByUId(ctx context.Context, uid int64, pageNum int64) ([]*model.UserCoupon, error) {
	dbModel := make([]*UserCoupon, 0)
	offset := (pageNum - 1) * constants.CouponPageSize
	if err := db.client.WithContext(ctx).Where("uid = ?", uid).Offset(int(offset)).Limit(constants.CouponPageSize).Find(&dbModel).Error; err != nil {
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

func (db *commodityDB) GetFullUserCouponsByUId(ctx context.Context, uid int64) ([]*model.UserCoupon, error) {
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
	dbModel := &UserCoupon{}
	err := db.client.WithContext(ctx).
		Where("uid = ? AND coupon_id = ?", coupon.Uid, coupon.CouponId).
		First(dbModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: user coupon not found")
		}
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to find user coupon: %v", err)
	}

	if dbModel.RemainingUses > 1 {
		dbModel.RemainingUses--
		if err := db.client.WithContext(ctx).Where("uid = ? AND coupon_id = ?", coupon.Uid, coupon.CouponId).Save(dbModel).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to save user coupon: %v", err)
		}
	} else {
		if err := db.client.WithContext(ctx).
			Where("uid = ? AND coupon_id = ?", coupon.Uid, coupon.CouponId).
			Delete(&UserCoupon{}).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete user coupon: %v", err)
		}
	}
	return nil
}
