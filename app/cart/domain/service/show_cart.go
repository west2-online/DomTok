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

package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/cart/domain/model"
)

func (svc *CartService) TryGetCartFromCache(ctx context.Context, uid int64, pageNum int64) (bool, []*model.CartGoods, error) {
	if pageNum != 1 && pageNum != 2 {
		return false, nil, nil
	}
	key := strconv.FormatInt(uid, 10)
	if svc.Cache.IsKeyExist(ctx, key) {
		cartStr, err := svc.Cache.GetCartCache(ctx, key)
		if err != nil {
			return false, nil, fmt.Errorf("TryGetCartFromCache get cart from cache error: %w", err)
		}

		var cartJson *model.CartJson
		if err = sonic.Unmarshal([]byte(cartStr), &cartJson); err != nil {
			return false, nil, fmt.Errorf("TryGetCartFromCache unmarshal error: %w", err)
		}

		cartGoods := model.ConvertCartJsonToCartGoods(cartJson)
		goods, err := svc.Rpc.GetGoodsInfo(ctx, cartGoods)
		if err != nil {
			return false, nil, fmt.Errorf("TryGetCartFromCache RPC error: %w", err)
		}
		return true, goods, nil
	}
	return false, nil, nil
}

func (svc *CartService) TrySetCartCache(ctx context.Context, uid int64, json string, pageNum int64) error {
	if pageNum != 1 && pageNum != 2 {
		return nil
	}
	key := strconv.FormatInt(uid, 10)
	err := svc.Cache.SetCartCache(ctx, key, json)
	if err != nil {
		return fmt.Errorf("TrySetCartCache set cart cache error: %w", err)
	}
	return nil
}
