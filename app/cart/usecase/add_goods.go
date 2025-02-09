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
	"strconv"

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/cart/entities"
	"github.com/west2-online/DomTok/pkg/constants"
)

func (u *UseCase) AddGoodsIntoCart(ctx context.Context, uid int64, goods *entities.GoodInfo) error {
	exist, _, err := u.DB.GetCartByUserId(ctx, uid)
	if err != nil {
		return fmt.Errorf("cartCase.AddGoodsIntoCart is cart exist err:%w", err)
	}

	cartJson := new(entities.CartJson)

	// 不存在该用户记录，添加用户购物车
	if !exist {
		err = u.createCart(ctx, uid, goods, cartJson)
		// 存在该用户记录，追加其购物车
	} else {
		err = u.appendCart(ctx, uid, goods, cartJson)
	}
	if err != nil {
		return err
	}

	// 接下来要更新缓存
	key := strconv.FormatInt(uid, 10)
	cacheCart := cartJson.GetRecentNStores(constants.RedisCartStoreNum)
	cacheCartJsonStr, err := sonic.MarshalString(cacheCart)
	if err != nil {
		return fmt.Errorf("cartCase.AddGoodsIntoCart marshalString err:%w", err)
	}
	err = u.Cache.SetCartCache(ctx, key, cacheCartJsonStr)
	if err != nil {
		return fmt.Errorf("cartCase.AddGoodsIntoCart cache err:%w", err)
	}
	return nil
}

func (u *UseCase) createCart(ctx context.Context, uid int64, goods *entities.GoodInfo, cartJson *entities.CartJson) error {
	cartJson.InsertSku(goods)
	cartJsonStr, err := sonic.MarshalString(cartJson)
	if err != nil {
		return fmt.Errorf("cartCase.createCart marshalString err:%w", err)
	}
	err = u.DB.CreateCart(ctx, uid, cartJsonStr)
	if err != nil {
		return fmt.Errorf("cartCase.createCart create cart err:%w", err)
	}
	return nil
}

func (u *UseCase) appendCart(ctx context.Context, uid int64, goods *entities.GoodInfo, cartJson *entities.CartJson) error {
	_, cartModel, err := u.DB.GetCartByUserId(ctx, uid)
	if err != nil {
		return fmt.Errorf("cartCase.appendCart get cartJsonStr err:%w", err)
	}
	err = sonic.UnmarshalString(cartModel.SkuJson, cartJson)
	if err != nil {
		return fmt.Errorf("cartCase.appendCart unmarshalString err:%w", err)
	}
	cartJson.InsertSku(goods)
	cartJsonStr, err := sonic.MarshalString(cartJson)
	if err != nil {
		return fmt.Errorf("cartCase.appendCart marshalString err:%w", err)
	}
	err = u.DB.SaveCart(ctx, uid, cartJsonStr)
	if err != nil {
		return fmt.Errorf("cartCase.appendCart save cart err:%w", err)
	}
	return nil
}
