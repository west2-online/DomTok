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

	"github.com/west2-online/DomTok/app/order/domain/model"
	idlmodel "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

func BuildOrder(o *model.Order) *idlmodel.Order {
	return &idlmodel.Order{
		Id:                    o.Id,
		Status:                strconv.Itoa(int(o.Status)),
		Uid:                   o.Uid,
		TotalAmountOfGoods:    utils.DecimalFloat64(&o.TotalAmountOfGoods),
		TotalAmountOfFreight:  utils.DecimalFloat64(&o.TotalAmountOfFreight),
		TotalAmountOfDiscount: utils.DecimalFloat64(&o.TotalAmountOfDiscount),
		PaymentAmount:         utils.DecimalFloat64(&o.PaymentAmount),
		PaymentStatus:         "待支付", // TODO
		PaymentAt:             o.PaymentAt,
		PaymentStyle:          o.PaymentStyle,
		OrderedAt:             o.OrderedAt,
		DeletedAt:             o.DeletedAt,
		DeliveryAt:            o.DeliveryAt,
		AddressID:             o.AddressID,
		AddressInfo:           o.AddressInfo,
		CouponId:              o.CouponId,
		CouponName:            o.CouponName,
	}
}

func BuildOrderGoods(g *model.OrderGoods) *idlmodel.OrderGoods {
	return &idlmodel.OrderGoods{
		OrderId:            g.OrderID,
		MerchantId:         g.MerchantID,
		GoodsId:            g.GoodsID,
		GoodsName:          g.GoodsName,
		StyleId:            g.StyleID,
		StyleName:          g.StyleName,
		GoodsVersion:       g.GoodsVersion,
		StyleHeadDrawing:   g.StyleHeadDrawing,
		OriginPrice:        utils.DecimalFloat64(&g.OriginPrice),
		SalePrice:          utils.DecimalFloat64(&g.SalePrice),
		SingleFreightPrice: utils.DecimalFloat64(&g.SingleFreightPrice),
		PurchaseQuantity:   g.PurchaseQuantity,
		TotalAmount:        utils.DecimalFloat64(&g.TotalAmount),
		FreightAmount:      utils.DecimalFloat64(&g.FreightAmount),
		DiscountAmount:     utils.DecimalFloat64(&g.DiscountAmount),
		PaymentAmount:      utils.DecimalFloat64(&g.PaymentAmount),
		SinglePrice:        utils.DecimalFloat64(&g.SinglePrice),
		CouponId:           g.CouponId,
		CouponName:         g.CouponName,
	}
}

func BuildOrderWithGoods(o *model.Order, goods []*model.OrderGoods) *idlmodel.OrderWithGoods {
	idlGoods := make([]*idlmodel.OrderGoods, len(goods))
	for i, g := range goods {
		idlGoods[i] = BuildOrderGoods(g)
	}

	return &idlmodel.OrderWithGoods{
		Order: BuildOrder(o),
		Goods: idlGoods,
	}
}

func BuildBaseOrder(o *model.Order) *idlmodel.BaseOrder {
	return &idlmodel.BaseOrder{
		Id:                 o.Id,
		Status:             constants.GetOrderStatusMsg(o.Status),
		TotalAmountOfGoods: utils.DecimalFloat64(&o.TotalAmountOfGoods),
		PaymentAmount:      utils.DecimalFloat64(&o.PaymentAmount),
		PaymentStatus:      "", // TODO 根据 o.paymentStatus 转化
	}
}

func BuildBaseOrderGoods(g *model.OrderGoods) *idlmodel.BaseOrderGoods {
	return &idlmodel.BaseOrderGoods{
		MerchantID:       g.MerchantID,
		GoodsID:          g.GoodsID,
		StyleID:          g.StyleID,
		PurchaseQuantity: g.PurchaseQuantity,
		GoodsVersion:     g.GoodsVersion,
	}
}
