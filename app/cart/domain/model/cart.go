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

package model

import "github.com/shopspring/decimal"

type Cart struct {
	UserId  int64
	SkuJson string
}

type GoodInfo struct {
	SkuId     int64
	ShopId    int64
	VersionId int64
	Count     int64
}

type CartGoods struct {
	MerchantID       int64
	GoodsID          int64
	GoodsName        string
	SkuID            int64
	SkuName          string
	GoodsVersion     int64
	StyleHeadDrawing string
	PurchaseQuantity int64
	TotalAmount      decimal.Decimal
	DiscountAmount   decimal.Decimal
}
