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
	"math"

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
		PageSize: math.MaxInt64,
	}
	skuInfoResp, err := rpc.commodity.ListSkuInfo(ctx, &skuReq)
	if err = utils.ProcessRpcError("commodity.ListSkuInfo", skuInfoResp, err); err != nil {
		return nil, err
	}

	orderGoods := lo.Map(skuInfoResp.SkuInfos, func(item *kmodel.SkuInfo, index int) *model.OrderGoods {
		purchaseCount := goods[index].PurchaseQuantity
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
			// SinglePrice: decimal.Decimal{}, 最终更新
			// CouponId: , // 优惠券模块负责
			// CouponName: "", 后续更新
		}
	})

	return orderGoods, nil
}

// WithholdSkuStock 预扣除商品数量
func (rpc *orderRpcImpl) WithholdSkuStock(ctx context.Context, stocks *model.OrderStock) error {
	infos := stockToSkuBuyInfo(stocks)

	resp, err := rpc.commodity.IncrSkuLockStock(ctx, &commodity.IncrSkuLockStockReq{Infos: infos})
	if err = utils.ProcessRpcError("commodity.IncrSkuLockStock", resp, err); err != nil {
		return err
	}

	return nil
}

// RollbackSkuStock 增加商品库存, 用于预扣接口的回滚
func (rpc *orderRpcImpl) RollbackSkuStock(ctx context.Context, stock *model.OrderStock) error {
	infos := stockToSkuBuyInfo(stock)

	resp, err := rpc.commodity.DescSkuLockStock(ctx, &commodity.DescSkuLockStockReq{Infos: infos})
	if err = utils.ProcessRpcError("commodity.IncrSkuLockStock", resp, err); err != nil {
		return err
	}

	return nil
}

// DescSkuStock 确认商品数量扣除
func (rpc *orderRpcImpl) DescSkuStock(ctx context.Context, stock *model.OrderStock) error {
	infos := stockToSkuBuyInfo(stock)

	resp, err := rpc.commodity.DescSkuStock(ctx, &commodity.DescSkuStockReq{Infos: infos})
	if err = utils.ProcessRpcError("commodity.IncrSkuLockStock", resp, err); err != nil {
		return err
	}

	return nil
}

// CalcOrderGoodsPrice 通过 coupon 的接口计算订单商品的最终价格
func (rpc *orderRpcImpl) CalcOrderGoodsPrice(ctx context.Context, goods []*model.OrderGoods) ([]*model.OrderGoods, error) {
	rpcOrderGoods := lo.Map(goods, func(g *model.OrderGoods, index int) *kmodel.OrderGoods {
		return &kmodel.OrderGoods{
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
	})

	resp, err := rpc.commodity.GetCouponAndPrice(ctx, &commodity.GetCouponAndPriceReq{GoodsList: rpcOrderGoods})
	if err = utils.ProcessRpcError("commodity.GetCouponAndPrice", resp, err); err != nil {
		return nil, err
	}

	return lo.Map(resp.AssignedGoodsList, func(item *kmodel.OrderGoods, index int) *model.OrderGoods {
		return &model.OrderGoods{
			OrderID:            item.OrderId,
			MerchantID:         item.MerchantId,
			GoodsID:            item.GoodsId,
			GoodsName:          item.GoodsName,
			StyleID:            item.StyleId,
			StyleName:          item.StyleName,
			GoodsVersion:       item.GoodsVersion,
			StyleHeadDrawing:   item.StyleHeadDrawing,
			OriginPrice:        decimal.NewFromFloat(item.OriginPrice),
			SalePrice:          decimal.NewFromFloat(item.SalePrice),
			SingleFreightPrice: decimal.NewFromFloat(item.SingleFreightPrice),
			PurchaseQuantity:   item.PurchaseQuantity,
			TotalAmount:        decimal.NewFromFloat(item.TotalAmount),
			FreightAmount:      decimal.NewFromFloat(item.FreightAmount),
			DiscountAmount:     decimal.NewFromFloat(item.DiscountAmount),
			PaymentAmount:      decimal.NewFromFloat(item.PaymentAmount),
			SinglePrice:        decimal.NewFromFloat(item.SinglePrice),
			CouponId:           item.CouponId,
			CouponName:         item.CouponName,
		}
	}), nil
}

func stockToSkuBuyInfo(stocks *model.OrderStock) []*kmodel.SkuBuyInfo {
	return lo.Map(stocks.Stocks, func(item *model.Stock, index int) *kmodel.SkuBuyInfo {
		return &kmodel.SkuBuyInfo{
			SkuID: item.SkuID,
			Count: item.Count,
		}
	})
}
