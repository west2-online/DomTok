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

	"github.com/west2-online/DomTok/app/commodity/usecase"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
)

type CommodityHandler struct {
	useCase usecase.CommodityUseCase
}

func (c CommodityHandler) CreateCoupon(ctx context.Context, req *commodity.CreateCouponReq) (r *commodity.CreateCouponResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) DeleteCoupon(ctx context.Context, req *commodity.DeleteCouponReq) (r *commodity.DeleteCouponResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) CreateUserCoupon(ctx context.Context, req *commodity.CreateCouponReq) (r *commodity.CreateUserCouponResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ViewCoupon(ctx context.Context, req *commodity.ViewCouponReq) (r *commodity.ViewCouponResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ViewUserAllCoupon(ctx context.Context, req *commodity.ViewCouponReq) (r *commodity.ViewUserAllCouponResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) UseUserCoupon(ctx context.Context, req *commodity.UseUserCouponReq) (r *commodity.UseUserCouponResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) CreateSpu(ctx context.Context, req *commodity.CreateSpuReq) (r *commodity.CreateSpuResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) UpdateSpu(ctx context.Context, req *commodity.UpdateSkuReq) (r *commodity.UpdateSpuResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ViewSpu(ctx context.Context, req *commodity.ViewSpuReq) (r *commodity.ViewSpuResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) DeleteSpu(ctx context.Context, req *commodity.DeleteSpuReq) (r *commodity.DeleteSpuResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ViewSpuImage(ctx context.Context, req *commodity.ViewSpuImageReq) (r *commodity.ViewSpuImageResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) CreateSku(ctx context.Context, req *commodity.CreateSkuReq) (r *commodity.CreateSkuResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) UpdateSku(ctx context.Context, req *commodity.UpdateSkuReq) (r *commodity.UpdateSkuResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) DeleteSku(ctx context.Context, req *commodity.DeleteSkuReq) (r *commodity.DeleteSkuResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ViewSkuImage(ctx context.Context, req *commodity.ViewSkuImageReq) (r *commodity.ViewSkuImageResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ViewSku(ctx context.Context, req *commodity.ViewSkuReq) (r *commodity.ViewSkuResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) UploadSkuAttr(ctx context.Context, req *commodity.UploadSkuAttrReq) (r *commodity.UploadSkuAttrResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ListSkuInfo(ctx context.Context, req *commodity.ListSkuInfoReq) (r *commodity.ListSkuInfoResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ViewHistory(ctx context.Context, req *commodity.ViewHistoryPriceReq) (r *commodity.ViewHistoryPriceResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) DescSkuLockStock(ctx context.Context, req *commodity.DescSkuLockStockReq) (r *commodity.DescSkuLockStockResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) IncrSkuLockStock(ctx context.Context, req *commodity.IncrSkuLockStockReq) (r *commodity.IncrSkuLockStockResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) DescSkuStock(ctx context.Context, req *commodity.DescSkuStockReq) (r *commodity.DescSkuStockResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) CreateCategory(ctx context.Context, req *commodity.CreateCategoryReq) (r *commodity.CreateCategoryResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) DeleteCategory(ctx context.Context, req *commodity.DeleteCategoryReq) (r *commodity.DeleteCategoryResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) ViewCategory(ctx context.Context, req *commodity.ViewCategoryReq) (r *commodity.ViewCategoryResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) UpdateCategory(ctx context.Context, req *commodity.UpdateCategoryReq) (r *commodity.UpdateCategoryResp, err error) {
	// TODO implement me
	panic("implement me")
}

func NewCommodityHandler(useCase usecase.CommodityUseCase) *CommodityHandler {
	return &CommodityHandler{useCase}
}
