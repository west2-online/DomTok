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

type Order struct {
	Id                    int64
	Status                int8
	Uid                   int64
	TotalAmountOfGoods    decimal.Decimal
	TotalAmountOfFreight  decimal.Decimal
	TotalAmountOfDiscount decimal.Decimal
	PaymentAmount         decimal.Decimal
	PaymentStatus         int8
	PaymentAt             int64
	PaymentStyle          string
	OrderedAt             int64
	DeletedAt             int64
	DeliveryAt            int64
	AddressID             int64
	AddressInfo           string
	CouponId              int64
	CouponName            string
}

type OrderGoods struct {
	OrderID            int64
	MerchantID         int64
	GoodsID            int64
	GoodsName          string
	StyleID            int64
	StyleName          string
	GoodsVersion       int64
	StyleHeadDrawing   string
	OriginPrice        decimal.Decimal
	SalePrice          decimal.Decimal
	SingleFreightPrice decimal.Decimal
	PurchaseQuantity   int64
	TotalAmount        decimal.Decimal
	FreightAmount      decimal.Decimal
	DiscountAmount     decimal.Decimal
	PaymentAmount      decimal.Decimal // 应付金额
	SinglePrice        decimal.Decimal
	CouponId           int64
	CouponName         string
}

type BaseOrderGoods struct {
	MerchantID       int64
	GoodsID          int64
	StyleID          int64
	GoodsVersion     int64
	PurchaseQuantity int64
	CouponID         int64
}

func OG2BOG(o *OrderGoods) *BaseOrderGoods {
	return &BaseOrderGoods{
		MerchantID:       o.MerchantID,
		GoodsID:          o.GoodsID,
		StyleID:          o.StyleID,
		GoodsVersion:     o.GoodsVersion,
		PurchaseQuantity: o.PurchaseQuantity,
		CouponID:         o.CouponId,
	}
}
