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
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
)

type Order struct {
	Id                    int64
	Status                int8
	Uid                   int64
	TotalAmountOfGoods    decimal.Decimal `gorm:"type:decimal(15,4)"`
	TotalAmountOfFreight  decimal.Decimal `gorm:"type:decimal(15,4)"`
	TotalAmountOfDiscount decimal.Decimal `gorm:"type:decimal(15,4)"`
	PaymentAmount         decimal.Decimal `gorm:"type:decimal(15,4)"`
	PaymentStatus         int8
	PaymentAt             int64
	PaymentStyle          string
	OrderedAt             int64
	DeletedAt             gorm.DeletedAt `gorm:"index"`
	DeliveryAt            int64
	AddressID             int64
	AddressInfo           string
	CouponId              int64
	CouponName            string
}

func (Order) TableName() string {
	return constants.OrderTableName
}

type OrderGoods struct {
	OrderID            int64
	MerchantID         int64
	GoodsID            int64
	GoodsVersion       int64
	GoodsName          string
	StyleID            int64
	StyleName          string
	StyleHeadDrawing   string
	OriginPrice        decimal.Decimal `gorm:"type:decimal(11,4)"`
	SalePrice          decimal.Decimal `gorm:"type:decimal(11,4)"`
	SingleFreightPrice decimal.Decimal `gorm:"type:decimal(11,4)"` // 单个运费金额
	PurchaseQuantity   int64
	TotalAmount        decimal.Decimal `gorm:"type:decimal(15,4)"` // 优惠前总金额 = (salePrice * count)+ freight
	FreightAmount      decimal.Decimal `gorm:"type:decimal(15,4)"`
	DiscountAmount     decimal.Decimal `gorm:"type:decimal(15,4)"`
	PaymentAmount      decimal.Decimal `gorm:"type:decimal(15,4)"` // 应付金额
	SinglePrice        decimal.Decimal `gorm:"type:decimal(11,4)"`
	CouponId           int64
	CouponName         string
}

func (OrderGoods) TableName() string {
	return constants.OrderGoodsTableName
}

func (db *orderDB) order2Model(order *Order) *model.Order {
	return &model.Order{
		Id:                    order.Id,
		Status:                order.Status,
		Uid:                   order.Uid,
		TotalAmountOfGoods:    order.TotalAmountOfGoods,
		TotalAmountOfFreight:  order.TotalAmountOfFreight,
		TotalAmountOfDiscount: order.TotalAmountOfDiscount,
		PaymentAmount:         order.PaymentAmount,
		PaymentStatus:         order.PaymentStatus,
		PaymentAt:             order.PaymentAt,
		PaymentStyle:          order.PaymentStyle,
		OrderedAt:             order.OrderedAt,
		//DeletedAt:             order.DeletedAt,
		DeliveryAt:  order.DeliveryAt,
		AddressID:   order.AddressID,
		AddressInfo: order.AddressInfo,
		CouponId:    order.CouponId,
		CouponName:  order.CouponName,
	}
}

func (db *orderDB) goods2Model(goods *OrderGoods) *model.OrderGoods {
	return &model.OrderGoods{
		OrderID:            goods.OrderID,
		MerchantID:         goods.MerchantID,
		GoodsID:            goods.GoodsID,
		GoodsName:          goods.GoodsName,
		StyleID:            goods.StyleID,
		StyleName:          goods.StyleName,
		GoodsVersion:       goods.GoodsVersion,
		StyleHeadDrawing:   goods.StyleHeadDrawing,
		OriginPrice:        goods.OriginPrice,
		SalePrice:          goods.SalePrice,
		SingleFreightPrice: goods.SingleFreightPrice,
		PurchaseQuantity:   goods.PurchaseQuantity,
		TotalAmount:        goods.TotalAmount,
		FreightAmount:      goods.FreightAmount,
		DiscountAmount:     goods.DiscountAmount,
		PaymentAmount:      goods.PaymentAmount,
		SinglePrice:        goods.SinglePrice,
		CouponId:           goods.CouponId,
		CouponName:         goods.CouponName,
	}
}

func (db *orderDB) model2Order(order *model.Order) *Order {
	return &Order{
		Id:                    order.Id,
		Status:                order.Status,
		Uid:                   order.Uid,
		TotalAmountOfGoods:    order.TotalAmountOfGoods,
		TotalAmountOfFreight:  order.TotalAmountOfFreight,
		TotalAmountOfDiscount: order.TotalAmountOfDiscount,
		PaymentAmount:         order.PaymentAmount,
		PaymentStatus:         order.PaymentStatus,
		PaymentAt:             order.PaymentAt,
		PaymentStyle:          order.PaymentStyle,
		OrderedAt:             order.OrderedAt,
		//DeletedAt:             gorm.DeletedAt{},
		DeliveryAt:  order.DeliveryAt,
		AddressID:   order.AddressID,
		AddressInfo: order.AddressInfo,
		CouponId:    order.CouponId,
		CouponName:  order.CouponName,
	}
}

func (db *orderDB) model2Goods(goods *model.OrderGoods) *OrderGoods {
	return &OrderGoods{
		OrderID:            goods.OrderID,
		MerchantID:         goods.MerchantID,
		GoodsID:            goods.GoodsID,
		GoodsName:          goods.GoodsName,
		StyleID:            goods.StyleID,
		StyleName:          goods.StyleName,
		GoodsVersion:       goods.GoodsVersion,
		StyleHeadDrawing:   goods.StyleHeadDrawing,
		OriginPrice:        goods.OriginPrice,
		SalePrice:          goods.SalePrice,
		SingleFreightPrice: goods.SingleFreightPrice,
		PurchaseQuantity:   goods.PurchaseQuantity,
		TotalAmount:        goods.TotalAmount,
		FreightAmount:      goods.FreightAmount,
		DiscountAmount:     goods.DiscountAmount,
		PaymentAmount:      goods.PaymentAmount,
		SinglePrice:        goods.SinglePrice,
		CouponId:           goods.CouponId,
		CouponName:         goods.CouponName,
	}
}
