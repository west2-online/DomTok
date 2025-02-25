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
	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/kitex_gen/commodity/commodityservice"
	"github.com/west2-online/DomTok/kitex_gen/user/userservice"
	"github.com/west2-online/DomTok/pkg/utils"
)

type orderRpcImpl struct {
	user      userservice.Client
	commodity commodityservice.Client
}

// TODO 等address 接口
func (rpc *orderRpcImpl) GetAddressInfo(ctx context.Context, addressId int64) (string, error) {
	return "", nil
}

// TODO 等 sku 接口完善后补齐
func (rpc *orderRpcImpl) QueryGoodsInfo(ctx context.Context, goods []*model.BaseOrderGoods) ([]*model.OrderGoods, error) {
	skuids := lo.Map(goods, func(item *model.BaseOrderGoods, index int) int64 {
		return item.StyleID
	})
	
	// TODO 新参数应包含 versionid
	skuReq := commodity.ListSkuInfoReq{
		SkuIDs:   skuids,
		PageNum:  1,
		PageSize: int64(len(skuids)),
	}

	skuInfoResp, err := rpc.commodity.ListSkuInfo(ctx, &skuReq)
	if err = utils.ProcessRpcError("commodity.ListSkuInfo", skuInfoResp, err); err != nil {
		return nil, err
	}

	// TODO ListSpuInfo
	//rpc.commodity.ListSkuInfo(ctx, &req)

	return nil, nil
}

// LockStock 预扣除商品数量
func (rpc *orderRpcImpl) LockStock(ctx context.Context, goods []*model.Stock) error {
	return nil
	// TODO
}

// ReleaseStock 再支付失败时进行回滚释放商品库存
func (rpc *orderRpcImpl) ReleaseStock(ctx context.Context, goods []*model.Stock) error {
	return nil
	// TODO
}

func (rpc *orderRpcImpl) DescSkuStock(ctx context.Context, goods []*model.OrderGoods) error {

}
