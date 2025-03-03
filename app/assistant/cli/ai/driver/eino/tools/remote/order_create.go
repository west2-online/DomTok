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
	"github.com/west2-online/DomTok/app/gateway/model/api/order"
	"github.com/west2-online/DomTok/app/gateway/model/model"

	strategy "github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/tools"
	"github.com/west2-online/DomTok/pkg/errno"
)

const (
	ToolOrderCreateName = "create_order"
	ToolOrderCreateDesc = "当用户想要下单时，可以使用此工具，目前只支持从购物车下单"
)

type ToolOrderCreateArgs struct {
	AddressID      int64                        `json:"address_id" desc:"填入用户地址ID" required:"true"`
	BaseOrderGoods []_OrderCreateBaseOrderGoods `json:"base_order_goods" desc:"填入订单商品信息" required:"true"`
}

type _OrderCreateBaseOrderGoods struct {
	MerchantID       int64 `json:"merchant_id" desc:"填入商家ID" required:"true"`
	GoodsID          int64 `json:"goods_id" desc:"填入商品ID" required:"true"`
	SkuID            int64 `json:"sku_id" desc:"填入SKU ID" required:"true"`
	GoodsVersionID   int64 `json:"goods_version_id" desc:"填入商品版本ID" required:"true"`
	PurchaseQuantity int64 `json:"purchase_quantity" desc:"填入购买数量" required:"true"`
}

var ToolOrderCreateRequestBody = schema.NewParamsOneOfByParams(*tools.Reflect(ToolOrderCreateArgs{}))

type ToolOrderCreate struct {
	tool.InvokableTool

	getServerCaller strategy.GetServerCaller
}

func OrderCreate(strategy strategy.GetServerCaller) *ToolOrderCreate {
	return &ToolOrderCreate{
		getServerCaller: strategy,
	}
}

func (t *ToolOrderCreate) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	s := t.getServerCaller(ToolOrderCreateName)
	if s == nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "server is nil")
	}
	params := &ToolOrderCreateArgs{}
	if err := sonic.Unmarshal([]byte(argumentsInJSON), params); err != nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "unmarshal arguments failed")
	}
	resp, err := s.OrderCreate(ctx, &order.CreateOrderReq{
		AddressID:      params.AddressID,
		BaseOrderGoods: ConvertArgsOrderGoodsToRequestGoods(params.BaseOrderGoods...),
	})
	if err != nil {
		return err.Error(), nil //nolint: nilerr
	}

	return string(resp), nil
}

func (t *ToolOrderCreate) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ToolOrderCreateName,
		Desc:        ToolOrderCreateDesc,
		ParamsOneOf: ToolOrderCreateRequestBody,
	}, nil
}

func (b *_OrderCreateBaseOrderGoods) ToRequestGoods() *model.BaseOrderGoods {
	return &model.BaseOrderGoods{
		MerchantID:       b.MerchantID,
		GoodsID:          b.GoodsID,
		GoodsVersion:     b.GoodsVersionID,
		PurchaseQuantity: b.PurchaseQuantity,
	}
}

func ConvertArgsOrderGoodsToRequestGoods(gs ..._OrderCreateBaseOrderGoods) []*model.BaseOrderGoods {
	res := make([]*model.BaseOrderGoods, len(gs))
	for i, g := range gs {
		res[i] = g.ToRequestGoods()
	}
	return res
}
