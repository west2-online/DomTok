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
	"math/rand/v2"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/west2-online/DomTok/app/order/domain/model"
)

func buildTestModelOrder(t *testing.T) *model.Order {
	t.Helper()
	return &model.Order{
		Id:                    rand.Int64(),
		Status:                1,
		Uid:                   1,
		TotalAmountOfGoods:    decimal.NewFromFloat(100),
		TotalAmountOfFreight:  decimal.NewFromFloat(10),
		TotalAmountOfDiscount: decimal.NewFromFloat(10),
		PaymentAmount:         decimal.NewFromFloat(100),
		PaymentStatus:         -1,
		PaymentAt:             0,
		PaymentStyle:          "支付宝",
		OrderedAt:             time.Now().UnixMilli(),
		DeletedAt:             0,
		DeliveryAt:            0,
		AddressID:             1,
		AddressInfo:           "fake address",
	}
}

func buildTestModelOrderGoods(t *testing.T, id int64) []*model.OrderGoods {
	t.Helper()
	var rel []*model.OrderGoods
	rel = append(rel, &model.OrderGoods{
		OrderID:            id,
		MerchantID:         1,
		GoodsID:            1,
		GoodsVersion:       1,
		GoodsName:          "fake goods",
		StyleID:            1,
		StyleName:          "fake style",
		StyleHeadDrawing:   "fake drawing",
		OriginPrice:        decimal.NewFromFloat(100),
		SalePrice:          decimal.NewFromFloat(100),
		SingleFreightPrice: decimal.NewFromFloat(10),
		PurchaseQuantity:   1,
		TotalAmount:        decimal.NewFromFloat(110),
		FreightAmount:      decimal.NewFromFloat(10),
		DiscountAmount:     decimal.NewFromFloat(0),
		PaymentAmount:      decimal.NewFromFloat(110),
		SinglePrice:        decimal.NewFromFloat(110),
		CouponId:           0,
		CouponName:         "",
	})
	return rel
}
