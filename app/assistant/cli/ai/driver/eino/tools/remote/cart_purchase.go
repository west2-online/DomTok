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

package remote

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	strategy "github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/tools"
	"github.com/west2-online/DomTok/app/gateway/model/api/cart"
	"github.com/west2-online/DomTok/app/gateway/model/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

const (
	ToolCartPurchaseName = "cart_purchase"
	ToolCartPurchaseDesc = "当用户想要下单时，可以使用此工具。从购物车中购买商品，需要购物车中的商品信息(调用Tool: show_cart)"
)

type ToolCartPurchaseArgs struct {
	BaseOrderGoods []_CartPurchaseBaseOrderGoods `json:"base_order_goods" desc:"填入订单商品信息" required:"true"`
}

type _CartPurchaseBaseOrderGoods struct {
	MerchantID       int64 `json:"merchant_id"       desc:"填入商家ID"             required:"true"`
	GoodsID          int64 `json:"goods_id"          desc:"填入商品ID"             required:"true"`
	SkuID            int64 `json:"sku_id"            desc:"填入SKU ID"           required:"true"`
	GoodsVersion     int64 `json:"goods_version"     desc:"填入商品版本"             required:"true"`
	PurchaseQuantity int64 `json:"purchase_quantity" desc:"填入购买数量"             required:"true"`
}

var ToolCartPurchaseRequestBody = schema.NewParamsOneOfByParams(*tools.Reflect(ToolCartPurchaseArgs{}))

type ToolCartPurchase struct {
	tool.InvokableTool

	getServerCaller strategy.GetServerCaller
}

func CartPurchase(strategy strategy.GetServerCaller) *ToolCartPurchase {
	return &ToolCartPurchase{
		getServerCaller: strategy,
	}
}

func (t *ToolCartPurchase) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	s := t.getServerCaller(ToolCartPurchaseName)
	if s == nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "server is nil")
	}
	params := &ToolCartPurchaseArgs{}
	if err := sonic.Unmarshal([]byte(argumentsInJSON), params); err != nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "unmarshal arguments failed")
	}
	resp, err := s.CartPurchase(ctx, &cart.PurChaseCartGoodsRequest{
		CartGoods: ConvertArgsOrderGoodsToRequestGoods(params.BaseOrderGoods...),
	})
	if err != nil {
		return err.Error(), nil //nolint: nilerr
	}

	return string(resp), nil
}

func (t *ToolCartPurchase) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ToolCartPurchaseName,
		Desc:        ToolCartPurchaseDesc,
		ParamsOneOf: ToolCartPurchaseRequestBody,
	}, nil
}

func (b *_CartPurchaseBaseOrderGoods) ToRequestGoods() *model.CartGoods {
	return &model.CartGoods{
		MerchantId:       b.MerchantID,
		GoodsId:          b.GoodsID,
		SkuId:            b.SkuID,
		GoodsVersion:     b.GoodsVersion,
		PurchaseQuantity: b.PurchaseQuantity,
	}
}

func ConvertArgsOrderGoodsToRequestGoods(gs ..._CartPurchaseBaseOrderGoods) []*model.CartGoods {
	res := make([]*model.CartGoods, len(gs))
	for i, g := range gs {
		res[i] = g.ToRequestGoods()
	}
	return res
}
