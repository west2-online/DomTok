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

package rpc

import (
	"context"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/app/cart/domain/repository"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/kitex_gen/commodity/commodityservice"
	kmodel "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/kitex_gen/order"
	"github.com/west2-online/DomTok/kitex_gen/order/orderservice"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/utils"
)

type CartRpcImpl struct {
	commodity commodityservice.Client
	order     orderservice.Client
}

func NewCartRpcImpl(c commodityservice.Client, o orderservice.Client) repository.RpcPort {
	return &CartRpcImpl{
		commodity: c,
		order:     o,
	}
}

func (rpc *CartRpcImpl) GetGoodsInfo(ctx context.Context, cartGoodsIds []*model.CartGoods) ([]*model.CartGoods, error) {
	skuVs := lo.Map(cartGoodsIds, func(item *model.CartGoods, index int) *kmodel.SkuVersion {
		v := &kmodel.SkuVersion{
			SkuID:     item.SkuID,
			VersionID: item.GoodsVersion,
		}
		return v
	})
	skuReq := commodity.ListSkuInfoReq{
		SkuInfos: skuVs,
		PageNum:  1,
		PageSize: int64(len(skuVs)),
	}
	skuInfoResp, err := rpc.commodity.ListSkuInfo(ctx, &skuReq)
	if err = utils.ProcessRpcError("commodity.GetSkuInfo", skuInfoResp, err); err != nil {
		return nil, errno.Errorf(errno.InternalRPCErrorCode, "call commodity.GetSkuInfo failed : %v", err)
	}
	cartGoods := lo.Map(skuInfoResp.SkuInfos, func(item *kmodel.SkuInfo, index int) *model.CartGoods {
		purchaseCount := cartGoodsIds[index].PurchaseQuantity
		return &model.CartGoods{
			MerchantID:       item.CreatorID,
			GoodsID:          item.SpuID,
			GoodsName:        item.Name,
			SkuID:            item.SkuID,
			SkuName:          item.Name,
			GoodsVersion:     item.HistoryID,
			StyleHeadDrawing: item.StyleHeadDrawing,
			PurchaseQuantity: purchaseCount,
			TotalAmount:      decimal.NewFromInt(int64(item.Price) * purchaseCount),
			DiscountAmount:   decimal.NewFromInt(int64(item.Price) * purchaseCount), // 暂时不去调rpc了，时间不是很够
		}
	})
	return cartGoods, nil
}

func (rpc *CartRpcImpl) PurchaseCartGoods(ctx context.Context, cartGoods []*model.CartGoods) (int64, error) {
	baseOrderGoods := lo.Map(cartGoods, func(item *model.CartGoods, index int) *kmodel.BaseOrderGoods {
		return &kmodel.BaseOrderGoods{
			MerchantID:       item.MerchantID,
			GoodsID:          item.GoodsID,
			StyleID:          item.SkuID,
			GoodsVersion:     item.GoodsVersion,
			PurchaseQuantity: item.PurchaseQuantity,
		}
	})
	resp, err := rpc.order.CreateOrder(ctx, &order.CreateOrderReq{
		BaseOrderGoods: baseOrderGoods,
	})
	if err = utils.ProcessRpcError("order.CreateOrder", resp, err); err != nil {
		return -1, errno.Errorf(errno.InternalRPCErrorCode, "call order.CreateOrder failed : %v", err)
	}
	return resp.OrderID, nil
}
