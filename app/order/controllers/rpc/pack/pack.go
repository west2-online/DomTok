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
	model2 "github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/kitex_gen/model"
)

func BuildOrder(o *model2.Order) *model.OrderGoods {
	return &model.OrderGoods{
		MerchantID:       o.UserID,
		GoodsID:          o.ID,
		GoodsName:        "",
		GoodsHeadDrawing: "",
		StyleID:          0,
		StyleName:        "",
		StyleHeadDrawing: "",
		OriginCast:       0,
		SaleCast:         0,
		PaymentAmount:    0,
		FreightAmount:    0,
		SettlementAmount: 0,
		DiscountAmount:   0,
		SingleCast:       0,
	}
}

func BuildOrderGoods(g *model2.OrderGoods) *model.OrderGoods {
	return &model.OrderGoods{
		GoodsID:          g.GoodsID,
		PurchaseQuantity: int64(g.Quantity),
	}
}
