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

package pack

import (
	"strconv"

	model2 "github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/kitex_gen/model"
)

func BuildOrder(o *model2.Order) *model.Order {
	return &model.Order{
		Id:                    o.Id,
		Status:                strconv.Itoa(int(o.Status)),
		Uid:                   o.Uid,
		TotalAmountOfGoods:    o.TotalAmountOfGoods,
		TotalAmountOfFreight:  o.TotalAmountOfFreight,
		TotalAmountOfDiscount: o.TotalAmountOfDiscount,
		PaymentAmount:         o.PaymentAmount,
		PaymentStatus:         o.PaymentStatus,
		PaymentAt:             o.PaymentAt,
		PaymentStyle:          o.PaymentStyle,
		OrderedAt:             o.OrderedAt,
		DeletedAt:             o.DeletedAt,
		DeliveryAt:            o.DeliveryAt,
		AddressID:             o.AddressID,
		AddressInfo:           o.AddressInfo,
	}
}

func BuildOrderGoods(g *model2.OrderGoods) *model.OrderGoods {
	return &model.OrderGoods{
		MerchantID:       g.MerchantID,
		GoodsID:          g.GoodsID,
		GoodsName:        g.GoodsName,
		GoodsHeadDrawing: g.GoodsHeadDrawing,
		StyleID:          g.StyleID,
		StyleName:        g.StyleName,
		StyleHeadDrawing: g.StyleHeadDrawing,
		OriginCast:       g.OriginCast,
		SaleCast:         g.SaleCast,
		PurchaseQuantity: g.PurchaseQuantity,
		PaymentAmount:    g.PaymentAmount,
		FreightAmount:    g.FreightAmount,
		SettlementAmount: g.SettlementAmount,
		DiscountAmount:   g.DiscountAmount,
		SingleCast:       g.SingleCast,
		CouponID:         g.CouponID,
	}
}

func BuildOrderWithGoods(o *model2.Order, goods []*model2.OrderGoods) *model.OrderWithGoods {
	idlGoods := make([]*model.OrderGoods, len(goods))
	for i, g := range goods {
		idlGoods[i] = BuildOrderGoods(g)
	}

	return &model.OrderWithGoods{
		Order: &model.Order{
			Id:                    o.Id,
			Status:                strconv.Itoa(int(o.Status)),
			Uid:                   o.Uid,
			TotalAmountOfGoods:    o.TotalAmountOfGoods,
			TotalAmountOfFreight:  o.TotalAmountOfFreight,
			TotalAmountOfDiscount: o.TotalAmountOfDiscount,
			PaymentAmount:         o.PaymentAmount,
			PaymentStatus:         o.PaymentStatus,
			PaymentAt:             o.PaymentAt,
			PaymentStyle:          o.PaymentStyle,
			OrderedAt:             o.OrderedAt,
			DeletedAt:             o.DeletedAt,
			DeliveryAt:            o.DeliveryAt,
			AddressID:             o.AddressID,
			AddressInfo:           o.AddressInfo,
		},
		Goods: idlGoods,
	}
}

func BuildBaseOrder(o *model2.Order) *model.BaseOrder {
	return &model.BaseOrder{
		Id:                 o.Id,
		Status:             strconv.Itoa(int(o.Status)),
		TotalAmountOfGoods: o.TotalAmountOfGoods,
		PaymentAmount:      o.PaymentAmount,
		PaymentStatus:      o.PaymentStatus,
	}
}

func BuildBaseOrderGoods(g *model2.OrderGoods) *model.BaseOrderGoods {
	return &model.BaseOrderGoods{
		MerchantName:     strconv.FormatInt(g.MerchantID, 10),
		GoodsName:        g.GoodsID,
		StyleName:        g.StyleID,
		PurchaseQuantity: g.PurchaseQuantity,
		StyleHeadDrawing: g.StyleHeadDrawing,
		CouponID:         g.CouponID,
	}
}
