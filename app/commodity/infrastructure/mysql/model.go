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
	"time"

	"gorm.io/gorm"

	"github.com/west2-online/DomTok/pkg/constants"
)

type Category struct {
	Id        int64 `gorm:"primary_key"`
	Name      string
	CreatorId int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// gorm.Model
}

type Spu struct {
	Id               int64 `gorm:"primary_key"`
	Name             string
	CreatorId        int64
	Description      string
	CategoryId       int64
	GoodsHeadDrawing string
	Price            float64
	ForSale          int
	Shipping         float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type SpuImage struct {
	Id        int64 `gorm:"primary_key"`
	Url       string
	SpuId     int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//	gorm.Model
}

type SpuToSku struct {
	SkuId     int64
	SpuId     int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Coupon struct {
	Id             int64 `gorm:"primary_key"`
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
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type UserCoupon struct {
	Uid           int64
	CouponId      int64
	RemainingUses int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Sku struct {
	Id               int64 `gorm:"primary_key"`
	CreatorId        int64
	Price            float64
	Name             string
	Description      string
	ForSale          int
	HistoryStock     int64
	StyleHeadDrawing string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Stock            int64
	LockStock        int64
}

// 对应表名

func (spu *Spu) TableName() string {
	return constants.SpuTableName
}

func (spu *SpuImage) TableName() string {
	return constants.SpuImageTableName
}

func (Category) TableName() string {
	return constants.CategoryTableName
}

func (s *SpuToSku) TableName() string {
	return constants.SpuSkuTableName
}

func (Coupon) TableName() string {
	return constants.CouponTableName
}

func (UserCoupon) TableName() string {
	return constants.UserCouponTableName
}
