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
	"bytes"
	"context"
	"time"

	"github.com/samber/lo"
	"fmt"

	"github.com/west2-online/DomTok/app/commodity/controllers/rpc/pack"
	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/usecase"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	kmodel "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/base"
)

type CommodityHandler struct {
	useCase usecase.CommodityUseCase
}

func (c CommodityHandler) ListSpuInfo(ctx context.Context, req *commodity.ListSpuInfoReq) (r *commodity.ListSpuInfoResp, err error) {
	r = new(commodity.ListSpuInfoResp)
	spus, err := c.useCase.ListSpuInfo(ctx, req.SpuIDs)
	if err != nil {
		return nil, err
	}
	r.Base = base.BuildBaseResp(nil)
	r.Spus = pack.BuildSpus(spus)
	return r, err
}

func (c CommodityHandler) CreateSpuImage(streamServer commodity.CommodityService_CreateSpuImageServer) (err error) {
	resp := new(commodity.CreateSpuImageResp)
	req, err := streamServer.Recv()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	for i := 0; i < int(req.BufferCount); i++ {
		data, err := streamServer.Recv()
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return streamServer.SendAndClose(resp)
		}
		req.Data = bytes.Join([][]byte{req.Data, data.Data}, []byte(""))
	}

	id, err := c.useCase.CreateSpuImage(streamServer.Context(), &model.SpuImage{
		Data:  req.Data,
		SpuID: req.SpuID,
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	resp.Base = base.BuildBaseResp(nil)
	resp.ImageID = id
	return streamServer.SendAndClose(resp)
}

func (c CommodityHandler) UpdateSpuImage(streamServer commodity.CommodityService_UpdateSpuImageServer) (err error) {
	resp := new(commodity.UpdateSpuImageResp)
	req, err := streamServer.Recv()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	for i := 0; i < int(req.BufferCount); i++ {
		data, err := streamServer.Recv()
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return streamServer.SendAndClose(resp)
		}
		req.Data = bytes.Join([][]byte{req.Data, data.Data}, []byte(""))
	}

	err = c.useCase.UpdateSpuImage(streamServer.Context(), &model.SpuImage{
		Data:    req.Data,
		ImageID: req.ImageID,
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	resp.Base = base.BuildBaseResp(nil)
	return streamServer.SendAndClose(resp)
}

func (c CommodityHandler) DeleteSpuImage(ctx context.Context, req *commodity.DeleteSpuImageReq) (r *commodity.DeleteSpuImageResp, err error) {
	resp := new(commodity.DeleteSpuImageResp)
	err = c.useCase.DeleteSpuImage(ctx, req.GetSpuImageID())
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, err
	}

	resp.Base = base.BuildBaseResp(nil)
	return resp, nil
}

func (c CommodityHandler) CreateCoupon(ctx context.Context, req *commodity.CreateCouponReq) (r *commodity.CreateCouponResp, err error) {
	r = new(commodity.CreateCouponResp)
	coupon := &model.Coupon{
		Name:           req.Name,
		TypeInfo:       int64(req.TypeInfo),
		ConditionCost:  *req.ConditionCost,
		DiscountAmount: *req.DiscountAmount,
		Discount:       *req.Discount,
		RangeType:      int64(req.RangeType),
		RangeId:        req.RangeID,
		Description:    *req.Description,
		ExpireTime:     time.Unix(req.ExpireTime, 0),
		DeadlineForGet: time.Unix(req.ExpireTime, 0),
	}
	couponId, err := c.useCase.CreateCoupon(ctx, coupon)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildBaseResp(nil)
	r.CouponID = couponId
	return r, nil
}

func (c CommodityHandler) DeleteCoupon(ctx context.Context, req *commodity.DeleteCouponReq) (r *commodity.DeleteCouponResp, err error) {
	r = new(commodity.DeleteCouponResp)
	coupon := &model.Coupon{
		Id: req.CouponID,
	}
	err = c.useCase.DeleteCoupon(ctx, coupon)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildBaseResp(nil)
	return r, nil
}

func (c CommodityHandler) CreateUserCoupon(ctx context.Context, req *commodity.CreateUserCouponReq) (r *commodity.CreateUserCouponResp, err error) {
	r = new(commodity.CreateUserCouponResp)
	userCoupon := &model.UserCoupon{
		CouponId:      req.CouponID,
		RemainingUses: req.RemainingUse,
	}
	err = c.useCase.CreateUserCoupon(ctx, userCoupon)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildBaseResp(nil)
	return r, nil
}

func (c CommodityHandler) ViewCoupon(ctx context.Context, req *commodity.ViewCouponReq) (r *commodity.ViewCouponResp, err error) {
	r = new(commodity.ViewCouponResp)
	coupons, err := c.useCase.GetCreatorCoupons(ctx, req.PageNum)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildBaseResp(nil)
	r.CouponInfo = pack.BuildCoupons(coupons)
	return r, nil
}

func (c CommodityHandler) ViewUserAllCoupon(ctx context.Context, req *commodity.ViewUserAllCouponReq) (r *commodity.ViewUserAllCouponResp, err error) {
	r = new(commodity.ViewUserAllCouponResp)
	coupons, err := c.useCase.SearchUserCoupons(ctx, req.PageNum)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildBaseResp(nil)
	r.Coupons = pack.BuildCoupons(coupons)
	return r, nil
}

func (c CommodityHandler) UseUserCoupon(ctx context.Context, req *commodity.UseUserCouponReq) (r *commodity.UseUserCouponResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) CreateSpu(streamServer commodity.CommodityService_CreateSpuServer) (err error) {
	resp := new(commodity.CreateSpuResp)

	req, err := streamServer.Recv()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	for i := 0; i < int(req.BufferCount); i++ {
		fileData, err := streamServer.Recv()
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return streamServer.SendAndClose(resp)
		}
		req.GoodsHeadDrawing = bytes.Join([][]byte{req.GoodsHeadDrawing, fileData.GoodsHeadDrawing}, []byte(""))
	}

	id, err := c.useCase.CreateSpu(streamServer.Context(), &model.Spu{
		Name:             req.Name,
		Description:      req.Description,
		CategoryId:       req.CategoryID,
		Price:            req.Price,
		ForSale:          int(req.ForSale),
		Shipping:         req.Shipping,
		GoodsHeadDrawing: req.GoodsHeadDrawing,
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	resp.Base = base.BuildBaseResp(nil)
	resp.SpuID = id
	return streamServer.SendAndClose(resp)
}

func (c CommodityHandler) UpdateSpu(streamServer commodity.CommodityService_UpdateSpuServer) (err error) {
	resp := new(commodity.UpdateSpuResp)

	req, err := streamServer.Recv()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	for i := 0; i < int(*req.BufferCount); i++ {
		fileData, err := streamServer.Recv()
		if err != nil {
			resp.Base = base.BuildBaseResp(err)
			return streamServer.SendAndClose(resp)
		}
		req.GoodsHeadDrawing = bytes.Join([][]byte{req.GoodsHeadDrawing, fileData.GoodsHeadDrawing}, []byte(""))
	}

	err = c.useCase.UpdateSpu(streamServer.Context(), &model.Spu{
		SpuId:            req.SpuID,
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		CategoryId:       req.GetCategoryID(),
		Price:            req.GetPrice(),
		ForSale:          int(req.GetForSale()),
		Shipping:         req.GetShipping(),
		GoodsHeadDrawing: req.GetGoodsHeadDrawing(),
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	resp.Base = base.BuildBaseResp(nil)
	return streamServer.SendAndClose(resp)
}

func (c CommodityHandler) ViewSpu(ctx context.Context, req *commodity.ViewSpuReq) (r *commodity.ViewSpuResp, err error) {
	r = new(commodity.ViewSpuResp)
	res, total, err := c.useCase.ViewSpus(ctx, req)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, err
	}
	r.Base = base.BuildBaseResp(nil)
	r.Total = total
	r.Spus = pack.BuildSpus(res)
	return r, err
}

func (c CommodityHandler) DeleteSpu(ctx context.Context, req *commodity.DeleteSpuReq) (r *commodity.DeleteSpuResp, err error) {
	r = new(commodity.DeleteSpuResp)
	err = c.useCase.DeleteSpu(ctx, req.GetSpuID())
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, err
	}
	r.Base = base.BuildBaseResp(nil)
	return r, nil
}

func (c CommodityHandler) ViewSpuImage(ctx context.Context, req *commodity.ViewSpuImageReq) (r *commodity.ViewSpuImageResp, err error) {
	r = new(commodity.ViewSpuImageResp)
	offset := req.GetPageNum() * req.GetPageSize()
	imgs, total, err := c.useCase.ViewSpuImages(ctx, req.GetSpuID(), int(offset), int(req.GetPageSize()))
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, err
	}

	r.Base = base.BuildBaseResp(nil)
	r.Images = pack.BuildImages(imgs)
	r.Total = total
	return
}

func (c CommodityHandler) CreateSku(streamServer commodity.CommodityService_CreateSkuServer) (err error) {
	r := new(commodity.CreateSkuResp)

	req, err := streamServer.Recv()
	if err != nil {
		logger.Errorf("rpc.CreateSpu: receive error: %v", err)
		r.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(r)
	}

	for i := 0; i < int(req.BufferCount); i++ {
		fileData, err := streamServer.Recv()
		if err != nil {
			logger.Errorf("rpc.CreateSku: receive error: %v", err)
			r.Base = base.BuildBaseResp(err)
			return streamServer.SendAndClose(r)
		}
		req.StyleHeadDrawing = bytes.Join([][]byte{req.StyleHeadDrawing, fileData.StyleHeadDrawing}, []byte(""))
	}

	id, err := c.useCase.CreateSku(streamServer.Context(), &model.Sku{
		Name:             req.Name,
		Stock:            req.Stock,
		Description:      req.Description,
		StyleHeadDrawing: req.StyleHeadDrawing,
		Price:            req.Price,
		ForSale:          int(req.ForSale),
		SpuID:            req.SpuID,
	}, req.Ext)
	if err != nil {
		logger.Errorf("rpc.CreateSku: create sku error: %v", err)
		r.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(r)
	}

	r.Base = base.BuildBaseResp(nil)
	r.SkuID = id
	return streamServer.SendAndClose(r)
}

func (c CommodityHandler) UpdateSku(streamServer commodity.CommodityService_UpdateSkuServer) (rr error) {
	r := new(commodity.UpdateSkuResp)

	req, err := streamServer.Recv()
	if err != nil {
		logger.Errorf("rpc.UpdateSku: receive error: %v", err)
		r.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(r)
	}

	for i := 0; i < int(*req.BufferCount); i++ {
		fileData, err := streamServer.Recv()
		if err != nil {
			r.Base = base.BuildBaseResp(err)
			return streamServer.SendAndClose(r)
		}
		req.StyleHeadDrawing = bytes.Join([][]byte{req.StyleHeadDrawing, fileData.StyleHeadDrawing}, []byte(""))
	}

	err = c.useCase.UpdateSku(streamServer.Context(), &model.Sku{
		SkuID:            req.SkuID,
		Stock:            req.GetStock(),
		Description:      req.GetDescription(),
		StyleHeadDrawing: req.GetStyleHeadDrawing(),
		Price:            req.GetPrice(),
		ForSale:          int(req.GetForSale()),
	}, req.Ext)
	if err != nil {
		logger.Errorf("rpc.UpdateSku: update sku error: %v", err)
		r.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(r)
	}

	r.Base = base.BuildBaseResp(nil)
	return streamServer.SendAndClose(r)
}

func (c CommodityHandler) DeleteSku(ctx context.Context, req *commodity.DeleteSkuReq) (r *commodity.DeleteSkuResp, err error) {
	r = new(commodity.DeleteSkuResp)

	sku := &model.Sku{
		SkuID: req.SkuID,
	}

	if err = c.useCase.DeleteSku(ctx, sku); err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}

	r.Base = base.BuildBaseResp(nil)
	return
}

func (c CommodityHandler) ViewSku(ctx context.Context, req *commodity.ViewSkuReq) (*commodity.ViewSkuResp, error) {
	r := new(commodity.ViewSkuResp)

	var (
		isSpuId bool
		sku     model.Sku
	)

	if req.SkuID != nil {
		sku.SkuID = *req.SkuID
		isSpuId = false
	}
	if req.SpuID != nil {
		sku.SpuID = *req.SpuID
		isSpuId = true
	}
	if req.SkuID == nil && req.SpuID == nil {
		err := fmt.Errorf("ViewSku failed: skuID and spuID are both nil")
		r.Base = base.BuildBaseResp(err)
		return r, err
	}

	Skus, total, err := c.useCase.ViewSku(ctx, &sku, req.PageNum, req.PageSize, isSpuId)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, err
	}

	r.Base = base.BuildBaseResp(nil)
	r.Skus = pack.BuildSkus(Skus)
	r.Total = total
	return r, nil
}

func (c CommodityHandler) UploadSkuAttr(ctx context.Context, req *commodity.UploadSkuAttrReq) (r *commodity.UploadSkuAttrResp, err error) {
	r = new(commodity.UploadSkuAttrResp)

	attr := &model.AttrValue{
		SaleAttr:  req.SaleAttr,
		SaleValue: req.SaleValue,
	}

	sku := &model.Sku{
		SkuID: *req.SkuID,
	}

	if err = c.useCase.UploadSkuAttr(ctx, attr, sku); err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}

	r.Base = base.BuildBaseResp(nil)
	return
}

func (c CommodityHandler) ListSkuInfo(ctx context.Context, req *commodity.ListSkuInfoReq) (r *commodity.ListSkuInfoResp, err error) {
	r = new(commodity.ListSkuInfoResp)

	SkuInfos, total, err := c.useCase.ListSkuInfo(ctx, req.SkuIDs, req.PageNum, req.PageSize)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}

	r.Base = base.BuildBaseResp(nil)
	r.SkuInfos = pack.BuildSkuInfos(SkuInfos)
	r.Total = total
	return
}

func (c CommodityHandler) CreateSkuImage(streamServer commodity.CommodityService_CreateSkuImageServer) (err error) {
	resp := new(commodity.CreateSkuImageResp)
	req, err := streamServer.Recv()
	if err != nil {
		logger.Errorf("rpc.CreateSkuImage: receive error: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	for i := 0; i < int(req.BufferCount); i++ {
		data, err := streamServer.Recv()
		if err != nil {
			logger.Errorf("rpc.CreateSkuImage: receive error: %v", err)
			resp.Base = base.BuildBaseResp(err)
			return streamServer.SendAndClose(resp)
		}
		req.Data = bytes.Join([][]byte{req.Data, data.Data}, []byte(""))
	}
	id, err := c.useCase.CreateSkuImage(streamServer.Context(), &model.SkuImage{
		SkuID: req.SkuID,
	}, req.Data)
	if err != nil {
		logger.Errorf("rpc.CreateSkuImage: create sku image error: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	resp.Base = base.BuildBaseResp(nil)
	resp.ImageID = id
	return streamServer.SendAndClose(resp)
}

func (c CommodityHandler) UpdateSkuImage(streamServer commodity.CommodityService_UpdateSkuImageServer) (err error) {
	resp := new(commodity.UpdateSkuImageResp)
	req, err := streamServer.Recv()
	if err != nil {
		logger.Errorf("rpc.UpdateSkuImage: receive error: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	for i := 0; i < int(req.BufferCount); i++ {
		data, err := streamServer.Recv()
		if err != nil {
			logger.Errorf("rpc.UpdateSkuImage: receive error: %v", err)
			resp.Base = base.BuildBaseResp(err)
			return streamServer.SendAndClose(resp)
		}
		req.Data = bytes.Join([][]byte{req.Data, data.Data}, []byte(""))
	}

	err = c.useCase.UpdateSkuImage(streamServer.Context(), &model.SkuImage{
		ImageID: req.ImageID,
	}, req.Data)
	if err != nil {
		logger.Errorf("rpc.UpdateSkuImage: update sku image error: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return streamServer.SendAndClose(resp)
	}

	resp.Base = base.BuildBaseResp(nil)
	return streamServer.SendAndClose(resp)
}

func (c CommodityHandler) DeleteSkuImage(ctx context.Context, req *commodity.DeleteSkuImageReq) (r *commodity.DeleteSkuImageResp, err error) {
	r = new(commodity.DeleteSkuImageResp)

	if err = c.useCase.DeleteSkuImage(ctx, req.SkuImageID); err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}

	r.Base = base.BuildBaseResp(nil)
	return
}

func (c CommodityHandler) ViewSkuImage(ctx context.Context, req *commodity.ViewSkuImageReq) (r *commodity.ViewSkuImageResp, err error) {
	r = new(commodity.ViewSkuImageResp)

	sku := &model.Sku{
		SkuID: req.SkuID,
	}

	var (
		images []*model.SkuImage
		total  int64
	)

	if images, total, err = c.useCase.ViewSkuImages(ctx, sku, req.PageNum, req.PageSize); err != nil {
		r.Base = base.BuildBaseResp(err)
		return
	}

	r.Base = base.BuildBaseResp(nil)
	r.Images = pack.BuildSkuImages(images)
	r.Total = total
	return
}

func (c CommodityHandler) ViewHistory(ctx context.Context, req *commodity.ViewHistoryPriceReq) (r *commodity.ViewHistoryPriceResp, err error) {
	// TODO implement me
	panic("implement me")
}

func (c CommodityHandler) DescSkuLockStock(ctx context.Context, req *commodity.DescSkuLockStockReq) (r *commodity.DescSkuLockStockResp, err error) {
	r = new(commodity.DescSkuLockStockResp)
	infos := make([]*model.SkuBuyInfo, 0)
	for _, info := range req.Infos {
		infos = append(infos, &model.SkuBuyInfo{
			SkuID: info.SkuID,
			Count: info.Count,
		})
	}
	err = c.useCase.DecrLockStock(ctx, infos)
	if err != nil {
		return r, err
	}
	r.Base = base.BuildBaseResp(nil)
	return r, nil
}

func (c CommodityHandler) IncrSkuLockStock(ctx context.Context, req *commodity.IncrSkuLockStockReq) (r *commodity.IncrSkuLockStockResp, err error) {
	r = new(commodity.IncrSkuLockStockResp)
	infos := make([]*model.SkuBuyInfo, 0)
	for _, info := range req.Infos {
		infos = append(infos, &model.SkuBuyInfo{
			SkuID: info.SkuID,
			Count: info.Count,
		})
	}
	err = c.useCase.IncrLockStock(ctx, infos)
	if err != nil {
		return r, err
	}
	r.Base = base.BuildBaseResp(nil)
	return r, nil
}

func (c CommodityHandler) DescSkuStock(ctx context.Context, req *commodity.DescSkuStockReq) (r *commodity.DescSkuStockResp, err error) {
	r = new(commodity.DescSkuStockResp)
	infos := make([]*model.SkuBuyInfo, 0)
	for _, info := range req.Infos {
		infos = append(infos, &model.SkuBuyInfo{
			SkuID: info.SkuID,
			Count: info.Count,
		})
	}
	err = c.useCase.DecrStock(ctx, infos)
	if err != nil {
		return r, err
	}
	r.Base = base.BuildBaseResp(nil)
	return r, nil
}

func (c CommodityHandler) CreateCategory(ctx context.Context, req *commodity.CreateCategoryReq) (r *commodity.CreateCategoryResp, err error) {
	r = new(commodity.CreateCategoryResp)
	category := model.Category{
		Name: req.Name,
	}
	id, err := c.useCase.CreateCategory(ctx, &category)
	r.Base = base.BuildBaseResp(err)
	r.CategoryID = id
	return
}

func (c CommodityHandler) DeleteCategory(ctx context.Context, req *commodity.DeleteCategoryReq) (r *commodity.DeleteCategoryResp, err error) {
	r = new(commodity.DeleteCategoryResp)
	category := model.Category{
		Id: req.CategoryID,
	}
	err = c.useCase.DeleteCategory(ctx, &category)
	r.Base = base.BuildBaseResp(err)
	return
}

func (c CommodityHandler) ViewCategory(ctx context.Context, req *commodity.ViewCategoryReq) (r *commodity.ViewCategoryResp, err error) {
	r = new(commodity.ViewCategoryResp)
	cInfos, err := c.useCase.ViewCategory(ctx, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return r, err
	}

	r.CategoryInfo = lo.Map(cInfos, func(item *model.CategoryInfo, index int) *kmodel.CategoryInfo {
		return &kmodel.CategoryInfo{
			CategoryID: item.CategoryID,
			Name:       item.Name,
		}
	})
	return r, nil
}

func (c CommodityHandler) UpdateCategory(ctx context.Context, req *commodity.UpdateCategoryReq) (r *commodity.UpdateCategoryResp, err error) {
	r = new(commodity.UpdateCategoryResp)
	category := model.Category{
		Id:   req.CategoryID,
		Name: req.Name,
	}
	err = c.useCase.UpdateCategory(ctx, &category)
	r.Base = base.BuildBaseResp(err)
	return
}

func NewCommodityHandler(useCase usecase.CommodityUseCase) *CommodityHandler {
	return &CommodityHandler{useCase}
}
