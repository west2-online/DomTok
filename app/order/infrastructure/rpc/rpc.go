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

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/kitex_gen/commodity/commodityservice"
	kmodel "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/kitex_gen/user/userservice"
	"github.com/west2-online/DomTok/pkg/utils"
)

type orderRpcImpl struct {
	user      userservice.Client
	commodity commodityservice.Client
}

func NewOrderRpcImpl(u userservice.Client, c commodityservice.Client) repository.RPC {
	return &orderRpcImpl{u, c}
}

// TODO 等address 接口
func (rpc *orderRpcImpl) GetAddressInfo(ctx context.Context, addressId int64) (string, error) {
	return "", nil
}

func (rpc *orderRpcImpl) QueryGoodsInfo(ctx context.Context, goods []*model.BaseOrderGoods) ([]*model.OrderGoods, error) {
	skuVs := lo.Map(goods, func(item *model.BaseOrderGoods, index int) *kmodel.SkuVersion {
		v := &kmodel.SkuVersion{
			SkuID:     item.StyleID,
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
	if err = utils.ProcessRpcError("commodity.ListSkuInfo", skuInfoResp, err); err != nil {
		return nil, err
	}

	orderGoods := lo.Map(skuInfoResp.SkuInfos, func(item *kmodel.SkuInfo, index int) *model.OrderGoods {
		purchaseCount := goods[index].PurchaseQuantity
		couponId := goods[index].CouponID
		return &model.OrderGoods{
			// OrderID:            0, 后续由 service 模块赋值
			MerchantID:         item.CreatorID,
			GoodsID:            item.SpuID,
			GoodsName:          item.Name,
			StyleID:            item.SkuID,
			StyleName:          item.Name,
			GoodsVersion:       item.HistoryID,
			StyleHeadDrawing:   item.StyleHeadDrawing,
			OriginPrice:        decimal.NewFromInt(int64(item.Price)),
			SalePrice:          decimal.NewFromInt(int64(item.Price)),
			SingleFreightPrice: decimal.NewFromInt(0),
			PurchaseQuantity:   purchaseCount,
			TotalAmount:        decimal.NewFromInt(int64(item.Price) * purchaseCount),
			FreightAmount:      decimal.NewFromInt(0),
			// DiscountAmount:     decimal.Decimal{}, 优惠券计算
			//  : decimal.Decimal{}, 优惠券计算
			// SinglePrice: decimal.Decimal{}, 最终更新
			CouponId: couponId,
			// CouponName: "", 后续更新
		}
	})

	return orderGoods, nil
}

// DescSkuLockStock 预扣除商品数量
func (rpc *orderRpcImpl) DescSkuLockStock(ctx context.Context, stock *model.OrderStock) error {
	infos := lo.Map(stock.Stocks, func(item *model.Stock, index int) *kmodel.SkuBuyInfo {
		return &kmodel.SkuBuyInfo{
			SkuID: item.SkuID,
			Count: item.Count,
		}
	})

	resp, err := rpc.commodity.DescSkuLockStock(ctx, &commodity.DescSkuLockStockReq{Infos: infos})
	if err = utils.ProcessRpcError("commodity.DescSkuLockStock", resp, err); err != nil {
		return err
	}

	return nil
}

// IncrSkuLockStock 增加商品库存, 用于预扣接口的回滚
func (rpc *orderRpcImpl) IncrSkuLockStock(ctx context.Context, stock *model.OrderStock) error {
	infos := lo.Map(stock.Stocks, func(item *model.Stock, index int) *kmodel.SkuBuyInfo {
		return &kmodel.SkuBuyInfo{
			SkuID: item.SkuID,
			Count: item.Count,
		}
	})

	resp, err := rpc.commodity.IncrSkuLockStock(ctx, &commodity.IncrSkuLockStockReq{Infos: infos})
	if err = utils.ProcessRpcError("commodity.IncrSkuLockStock", resp, err); err != nil {
		return err
	}

	return nil
}

// DescSkuStock 确认商品数量扣除
func (rpc *orderRpcImpl) DescSkuStock(ctx context.Context, stock *model.OrderStock) error {
	infos := lo.Map(stock.Stocks, func(item *model.Stock, index int) *kmodel.SkuBuyInfo {
		return &kmodel.SkuBuyInfo{
			SkuID: item.SkuID,
			Count: item.Count,
		}
	})

	resp, err := rpc.commodity.DescSkuStock(ctx, &commodity.DescSkuStockReq{Infos: infos})
	if err = utils.ProcessRpcError("commodity.IncrSkuLockStock", resp, err); err != nil {
		return err
	}

	return nil
}
