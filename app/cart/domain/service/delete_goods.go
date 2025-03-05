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

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/cart/domain/model"
)

func (svc *CartService) DeleteGoods(ctx context.Context, uid int64, info []*model.GoodInfo) error {
	_, cartModel, err := svc.DB.GetCartByUserId(ctx, uid)
	if err != nil {
		return fmt.Errorf("CartService.appendCart get cartJsonStr err:%w", err)
	}
	cartJson := new(model.CartJson)
	err = sonic.UnmarshalString(cartModel.SkuJson, cartJson)
	if err != nil {
		return fmt.Errorf("CartService.appendCart unmarshalString err:%w", err)
	}
	for _, infoItem := range info {
		cartJson.DeleteSku(infoItem)
	}
	cartJsonStr, err := sonic.MarshalString(cartJson)
	if err != nil {
		return fmt.Errorf("CartService.appendCart marshalString err:%w", err)
	}
	err = svc.DB.SaveCart(ctx, uid, cartJsonStr)
	if err != nil {
		return fmt.Errorf("CartService.appendCart save cart err:%w", err)
	}
	return nil
}
