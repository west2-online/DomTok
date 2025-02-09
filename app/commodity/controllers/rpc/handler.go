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
	"github.com/west2-online/DomTok/app/commodity/entities"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/kitex_gen/model"
)

type UseCasePort interface {
	CreateCategory(ctx context.Context, category *entities.Category) (id int64, err error)
	DeleteCategory(ctx context.Context, category *entities.Category) (err error)
	UpdateCategory(ctx context.Context, category *entities.Category) (err error)
	ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*model.CategoryInfo, err error)
}

type CommodityHandler struct {
	useCase UseCasePort
}

func NewCommodityHandler(useCase UseCasePort) *CommodityHandler {
	return &CommodityHandler{useCase}
}

func (h *CommodityHandler) CreateCoupon(ctx context.Context, req *commodity.CreateCouponReq) (r *commodity.CreateCouponResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) DeleteCoupon(ctx context.Context, req *commodity.DeleteCouponReq) (r *commodity.DeleteCouponResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) CreateUserCoupon(ctx context.Context, req *commodity.CreateCouponReq) (r *commodity.CreateUserCouponResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) ViewCoupon(ctx context.Context, req *commodity.ViewCouponReq) (r *commodity.ViewCouponResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) ViewUserAllCoupon(ctx context.Context, req *commodity.ViewCouponReq) (r *commodity.ViewUserAllCouponResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) UseUserCoupon(ctx context.Context, req *commodity.UseUserCouponReq) (r *commodity.UseUserCouponResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) CreateSpu(ctx context.Context, req *commodity.CreateSpuReq) (r *commodity.CreateSpuResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) UpdateSpu(ctx context.Context, req *commodity.UpdateSkuReq) (r *commodity.UpdateSpuResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) ViewSpu(ctx context.Context, req *commodity.ViewSpuReq) (r *commodity.ViewSpuResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) DeleteSpu(ctx context.Context, req *commodity.DeleteSpuReq) (r *commodity.DeleteSpuResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) ViewSpuImage(ctx context.Context, req *commodity.ViewSpuImageReq) (r *commodity.ViewSpuImageResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) CreateSku(ctx context.Context, req *commodity.CreateSkuReq) (r *commodity.CreateSkuResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) UpdateSku(ctx context.Context, req *commodity.UpdateSkuReq) (r *commodity.UpdateSkuResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) DeleteSku(ctx context.Context, req *commodity.DeleteSkuReq) (r *commodity.DeleteSkuResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) ViewSkuImage(ctx context.Context, req *commodity.ViewSkuImageReq) (r *commodity.ViewSkuImageResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) ViewSku(ctx context.Context, req *commodity.ViewSkuReq) (r *commodity.ViewSkuResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) UploadSkuAttr(ctx context.Context, req *commodity.UploadSkuAttrReq) (r *commodity.UploadSkuAttrResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) ListSkuInfo(ctx context.Context, req *commodity.ListSkuInfoReq) (r *commodity.ListSkuInfoResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) ViewHistory(ctx context.Context, req *commodity.ViewHistoryPriceReq) (r *commodity.ViewHistoryPriceResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) DescSkuLockStock(ctx context.Context, req *commodity.DescSkuLockStockReq) (r *commodity.DescSkuLockStockResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) IncrSkuLockStock(ctx context.Context, req *commodity.IncrSkuLockStockReq) (r *commodity.IncrSkuLockStockResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) DescSkuStock(ctx context.Context, req *commodity.DescSkuStockReq) (r *commodity.DescSkuStockResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *CommodityHandler) CreateCategory(ctx context.Context, req *commodity.CreateCategoryReq) (resp *commodity.CreateCategoryResp, err error) {
	resp = new(commodity.CreateCategoryResp)
	return resp, nil
}

func (h *CommodityHandler) DeleteCategory(ctx context.Context, req *commodity.DeleteCategoryReq) (resp *commodity.DeleteCategoryResp, err error) {
	resp = new(commodity.DeleteCategoryResp)
	return resp, nil
}

func (h *CommodityHandler) UpdateCategory(ctx context.Context, req *commodity.UpdateCategoryReq) (resp *commodity.UpdateCategoryResp, err error) {
	resp = new(commodity.UpdateCategoryResp)
	return resp, nil
}

func (h *CommodityHandler) ViewCategory(ctx context.Context, req *commodity.ViewCategoryReq) (*commodity.ViewCategoryResp, error) {
	resp := new(commodity.ViewCategoryResp)
	return resp, nil
}
