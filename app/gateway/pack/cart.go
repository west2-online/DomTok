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
	"github.com/west2-online/DomTok/app/gateway/model/model"
	kmodel "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/base"
)

func BuildCartGoods(goods *kmodel.CartGoods) *model.CartGoods {
	return &model.CartGoods{
		MerchantId:       goods.MerchantId,
		GoodsId:          goods.GoodsId,
		GoodsName:        goods.GoodsName,
		SkuId:            goods.SkuId,
		SkuName:          goods.SkuName,
		GoodsVersion:     goods.GoodsVersion,
		StyleHeadDrawing: goods.StyleHeadDrawing,
		PurchaseQuantity: goods.PurchaseQuantity,
		TotalAmount:      goods.TotalAmount,
		DiscountAmount:   goods.DiscountAmount,
	}
}

func BuildCartGoodsList(goods []*kmodel.CartGoods) []*model.CartGoods {
	return base.BuildTypeList(goods, BuildCartGoods)
}
