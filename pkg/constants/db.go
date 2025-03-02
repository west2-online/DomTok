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

package constants

import "time"

const (
	MaxConnections  = 1000             // (DB) 最大连接数
	MaxIdleConns    = 10               // (DB) 最大空闲连接数
	ConnMaxLifetime = 10 * time.Second // (DB) 最大可复用时间
	ConnMaxIdleTime = 5 * time.Minute  // (DB) 最长保持空闲状态时间

	CouponMaxVarCharLen     = 255 // coupon的varchar相关字段最大值
	CouponRangeTypeSPU      = 1
	CouponRangeTypeCategory = 2
	CouponPageSize          = 15
	CouponTypeSubAmount     = 1
	CouponTypeDiscount      = 2
)

const (
	UserTableName       = "users"
	CategoryTableName   = "category"
	OrderTableName      = "orders"
	OrderGoodsTableName = "order_goods"
	SpuTableName        = "spu_info"
	SpuImageTableName   = "spu_image"
	CouponTableName     = "coupon_info"
	UserCouponTableName = "user_coupon"

	SpuSkuTableName        = "spu_to_sku"
	CartTableName          = "cart"
	PaymentTableName       = "payment_orders"
	PaymentRefundTableName = "payment_refunds"
	PaymentLedgerTableName = "payment_ledger"

	SkuTableName             = "sku_info"
	SkuImagesTableName       = "sku_image"
	SkuSaleAttrTableName     = "sku_sale_attr"
	SkuPriceHistoryTableName = "sku_price_history"
)
