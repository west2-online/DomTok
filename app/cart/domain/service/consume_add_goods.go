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
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
)

func (s *CartService) ConsumeAddGoods(ctx context.Context) {
	msgCh := s.MQ.Consume(ctx,
		constants.KafkaCartTopic,
		constants.KafkaConsumerNum,
		constants.KafkaCartAddGoodsGroupId,
		constants.DefaultConsumerChanCap)
	go func() {
		for msg := range msgCh {
			req := new(model.AddGoodsMsg)
			err := sonic.Unmarshal(msg.V, req)
			if err != nil {
				logger.Errorf("CartService.ConsumeAddGoods: Unmarshal err: %v", err)
			}
			err = s.addGoodsIntoCart(ctx, req.Uid, req.Goods)
			if err != nil {
				logger.Errorf("CartService.ConsumeAddGoods: addGoodsIntoCart err: %v", err)
			}
		}
	}()
}

func (s *CartService) addGoodsIntoCart(ctx context.Context, uid int64, goods *model.GoodInfo) error {
	exist, _, err := s.DB.GetCartByUserId(ctx, uid)
	if err != nil {
		return fmt.Errorf("CartService.AddGoodsIntoCart is cart exist err:%w", err)
	}
	cartJson := new(model.CartJson)

	// 不存在该用户记录，添加用户购物车
	if !exist {
		err = s.createCart(ctx, uid, goods, cartJson)
		// 存在该用户记录，追加其购物车
	} else {
		err = s.appendCart(ctx, uid, goods, cartJson)
	}
	if err != nil {
		return err
	}

	// 接下来要更新缓存
	key := strconv.FormatInt(uid, 10)
	cacheCart := cartJson.GetRecentNStores(constants.RedisCartStoreNum)
	cacheCartJsonStr, err := sonic.MarshalString(cacheCart)
	if err != nil {
		return fmt.Errorf("CartService.AddGoodsIntoCart marshalString err:%w", err)
	}
	err = s.Cache.SetCartCache(ctx, key, cacheCartJsonStr)
	if err != nil {
		return fmt.Errorf("CartService.AddGoodsIntoCart cache err:%w", err)
	}
	return nil
}

func (s *CartService) createCart(ctx context.Context, uid int64, goods *model.GoodInfo, cartJson *model.CartJson) error {
	cartJson.InsertSku(goods)
	cartJsonStr, err := sonic.MarshalString(cartJson)
	if err != nil {
		return fmt.Errorf("CartService.createCart marshalString err:%w", err)
	}
	err = s.DB.CreateCart(ctx, uid, cartJsonStr)
	if err != nil {
		return fmt.Errorf("CartService.createCart create cart err:%w", err)
	}
	return nil
}

func (s *CartService) appendCart(ctx context.Context, uid int64, goods *model.GoodInfo, cartJson *model.CartJson) error {
	_, cartModel, err := s.DB.GetCartByUserId(ctx, uid)
	if err != nil {
		return fmt.Errorf("CartService.appendCart get cartJsonStr err:%w", err)
	}
	err = sonic.UnmarshalString(cartModel.SkuJson, cartJson)
	if err != nil {
		return fmt.Errorf("CartService.appendCart unmarshalString err:%w", err)
	}
	cartJson.InsertSku(goods)
	cartJsonStr, err := sonic.MarshalString(cartJson)
	if err != nil {
		return fmt.Errorf("CartService.appendCart marshalString err:%w", err)
	}
	err = s.DB.SaveCart(ctx, uid, cartJsonStr)
	if err != nil {
		return fmt.Errorf("CartService.appendCart save cart err:%w", err)
	}
	return nil
}
