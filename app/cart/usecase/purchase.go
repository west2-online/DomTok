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

package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/west2-online/DomTok/app/cart/domain/model"
	metainfoContext "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/logger"
)

func (u *UseCase) PurChaseCartGoods(ctx context.Context, goodsList []*model.CartGoods) (int64, error) {
	orderId, err := u.Rpc.PurchaseCartGoods(ctx, goodsList)
	if err != nil {
		return 0, fmt.Errorf("PurChaseCartGoods: failed to create order: %w", err)
	}
	id, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return -1, fmt.Errorf("cartCase.AddGoodsIntoCart metainfo unmarshal error:%w", err)
	}
	go func() {
		goodsInfo := lo.Map(goodsList, func(item *model.CartGoods, index int) *model.GoodInfo {
			return &model.GoodInfo{
				SkuId:     item.SkuID,
				ShopId:    item.MerchantID,
				VersionId: item.GoodsVersion,
				Count:     item.PurchaseQuantity,
			}
		})
		err = u.svc.DeleteGoods(ctx, id, goodsInfo)
		if err != nil {
			logger.Errorf("PurChaseCartGoods: Failed to delete goods: %v", err)
		}
	}()
	return orderId, nil
}
