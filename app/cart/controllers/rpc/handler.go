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

	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/app/cart/usecase"
	"github.com/west2-online/DomTok/kitex_gen/cart"
	"github.com/west2-online/DomTok/pkg/base"
)

type CartHandler struct {
	useCase usecase.CartCasePort
}

func NewCartHandler(useCase usecase.CartCasePort) *CartHandler {
	return &CartHandler{useCase: useCase}
}

func (h *CartHandler) AddGoodsIntoCart(ctx context.Context, req *cart.AddGoodsIntoCartRequest) (r *cart.AddGoodsIntoCartResponse, err error) {
	r = new(cart.AddGoodsIntoCartResponse)
	// create model
	good := &model.GoodInfo{
		SkuId:  req.SkuId,
		ShopId: req.ShopId,
		Count:  req.Count,
	}
	// useCase
	err = h.useCase.AddGoodsIntoCart(ctx, good)
	r.Base = base.BuildBaseResp(err)
	return r, nil
}

func (h *CartHandler) ShowCartGoodsList(ctx context.Context, req *cart.ShowCartGoodsListRequest) (r *cart.ShowCartGoodsListResponse, err error) {
	r = new(cart.ShowCartGoodsListResponse)
	return r, nil
}

func (h *CartHandler) UpdateCartGoods(ctx context.Context, req *cart.UpdateCartGoodsRequest) (r *cart.UpdateCartGoodsResponse, err error) {
	r = new(cart.UpdateCartGoodsResponse)
	return r, nil
}

func (h *CartHandler) DeleteCartGoods(ctx context.Context, req *cart.DeleteAllCartGoodsRequest) (r *cart.DeleteAllCartGoodsResponse, err error) {
	r = new(cart.DeleteAllCartGoodsResponse)
	return r, nil
}

func (h *CartHandler) DeleteAllCartGoods(ctx context.Context, req *cart.DeleteAllCartGoodsRequest) (r *cart.DeleteAllCartGoodsResponse, err error) {
	r = new(cart.DeleteAllCartGoodsResponse)
	return r, nil
}

func (h *CartHandler) PayCartGoods(ctx context.Context, req *cart.PayCartGoodsRequest) (r *cart.PayCartGoodsResponse, err error) {
	r = new(cart.PayCartGoodsResponse)
	return r, nil
}
