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

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/cart/domain/model"
	metainfoContext "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/logger"
)

func (u *UseCase) ShowCartGoods(ctx context.Context, pageNum int64) ([]*model.CartGoods, error) {
	if err := u.svc.Verify(u.svc.VerifyPageNum(pageNum)); err != nil {
		return nil, err
	}
	userID, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return nil, fmt.Errorf("ShowCartGoods get user info error: %w", err)
	}

	e, res, err := u.svc.TryGetCartFromCache(ctx, userID, pageNum)
	if err != nil {
		return nil, fmt.Errorf("ShowCartGoods get cart from cache error: %w", err)
	}
	if e {
		return res, nil
	}

	exist, cartData, err := u.DB.GetCartByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ShowCartGoods DB get cart error: %w", err)
	}
	if !exist {
		return []*model.CartGoods{}, nil
	}

	var cartJson *model.CartJson
	if err = sonic.Unmarshal([]byte(cartData.SkuJson), &cartJson); err != nil {
		return nil, fmt.Errorf("ShowCartGoods unmarshal error: %w", err)
	}
	cartGoods := model.ConvertCartJsonToCartGoods(cartJson)
	goods, err := u.Rpc.GetGoodsInfo(ctx, cartGoods)
	if err != nil {
		return nil, fmt.Errorf("ShowCartGoods RPC error: %w", err)
	}

	go func() {
		err = u.svc.TrySetCartCache(ctx, userID, cartData.SkuJson, pageNum)
		if err != nil {
			logger.Errorf("ShowCartGoods set cart cache error: %v", err)
		}
	}()

	return goods, nil
}
