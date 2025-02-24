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

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

func (us *useCase) CreateCategory(ctx context.Context, category *model.Category) (id int64, err error) {
	return 0, nil
}

func (us *useCase) CreateSpu(ctx context.Context, spu *model.Spu) (id int64, err error) {
	loginData, err := contextLogin.GetStreamLoginData(ctx)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpu failed: %w", err)
	}
	spu.CreatorId = loginData

	if err = us.svc.Verify(us.svc.VerifyForSaleStatus(spu.ForSale)); err != nil {
		return 0, fmt.Errorf("usecase.CreateSpu verify failed: %w", err)
	}

	id, err = us.svc.CreateSpu(ctx, spu)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpu failed: %w", err)
	}
	return id, nil
}

func (us *useCase) CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) (int64, error) {
	_, err := us.db.GetSpuBySpuId(ctx, spuImage.SpuID)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpuImage failed: %w", err)
	}
	id, err := us.svc.CreateSpuImage(ctx, spuImage)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpuImage failed: %w", err)
	}
	return id, nil
}

func (us *useCase) DeleteSpu(ctx context.Context, spuId int64) error {
	ret, err := us.svc.MatchDeleteSpuCondition(ctx, spuId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}

	err = us.svc.IdentifyUser(ctx, ret.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpu identify user failed: %w", err)
	}

	if err = us.svc.DeleteSpu(ctx, spuId, ret.GoodsHeadDrawingUrl); err != nil {
		return fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}

	if err = us.svc.DeleteAllSpuImages(ctx, spuId); err != nil {
		return fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}

	return nil
}

func (us *useCase) UpdateSpu(ctx context.Context, spu *model.Spu) error {
	ret, err := us.db.GetSpuBySpuId(ctx, spu.SpuId)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpu failed: %w", err)
	}

	err = us.svc.IdentifyUserInStreamCtx(ctx, ret.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpu identify user failed: %w", err)
	}

	if err = us.svc.Verify(us.svc.VerifyForSaleStatus(spu.ForSale)); err != nil {
		return fmt.Errorf("usecase.UpdateSpu verify failed: %w", err)
	}

	spu.GoodsHeadDrawingUrl = utils.GenerateFileName(constants.SpuDirDest, spu.SpuId)
	if err = us.svc.UpdateSpu(ctx, spu, ret); err != nil {
		return fmt.Errorf("usecase.UpdateSpu failed: %w", err)
	}
	return nil
}

func (us *useCase) UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error {
	spu, img, err := us.svc.GetSpuFromImageId(ctx, spuImage.ImageID)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage failed: %w", err)
	}

	err = us.svc.IdentifyUserInStreamCtx(ctx, spu.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage identify user failed: %w", err)
	}

	spuImage.Url = utils.GenerateFileName(constants.SpuImageDirDest, img.ImageID)
	if err = us.svc.UpdateSpuImage(ctx, spuImage, img); err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}
	return nil
}

func (us *useCase) DeleteSpuImage(ctx context.Context, imageId int64) error {
	spu, img, err := us.svc.GetSpuFromImageId(ctx, imageId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage failed: %w", err)
	}

	err = us.svc.IdentifyUser(ctx, spu.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage identify user failed: %w", err)
	}

	if err = us.svc.DeleteSpuImage(ctx, imageId, img.Url); err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage failed: %w", err)
	}
	return nil
}

func (us *useCase) ViewSpuImages(ctx context.Context, spuId int64, offset, limit int) ([]*model.SpuImage, int64, error) {
	return us.svc.GetSpuImages(ctx, spuId, offset, limit)
}

func (us *useCase) CreateSku(ctx context.Context, sku *model.Sku, ext string) (skuID int64, err error) {
	loginData, err := contextLogin.GetStreamLoginData(ctx)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSku failed: %w", err)
	}
	sku.CreatorID = loginData

	skuID, err = us.svc.CreateSku(ctx, sku, ext)
	if err != nil {
		return -1, fmt.Errorf("usecase.CreateSku failed: %w", err)
	}
	return skuID, nil
}

func (us *useCase) UpdateSku(ctx context.Context, sku *model.Sku) (err error) {
	ret, err := us.db.GetSkuBySkuId(ctx, sku.SkuID)
	if err != nil {
		return fmt.Errorf("service.UpdateSku: get sku by sku id failed: %w", err)
	}

	if err := us.svc.IdentifyUserInStreamCtx(ctx, ret.CreatorID); err != nil {
		return fmt.Errorf("service.UpdateSku: %w", err)
	}

	sku.StyleHeadDrawingUrl = utils.GenerateFileName(constants.SkuDirDest, sku.SkuID)
	if err = us.svc.UpdateSku(ctx, sku, ret); err != nil {
		return fmt.Errorf("usecase.UpdateSku failed: %w", err)
	}

	return nil
}

func (us *useCase) DeleteSku(ctx context.Context, sku *model.Sku) (err error) {
	ret, err := us.db.GetSkuBySkuId(ctx, sku.SkuID)
	if err != nil {
		return fmt.Errorf("service.UpdateSku: get sku by sku id failed: %w", err)
	}

	if err := us.svc.IdentifyUser(ctx, ret.CreatorID); err != nil {
		return fmt.Errorf("service.UpdateSku: %w", err)
	}

	err = us.db.DeleteSku(ctx, sku)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSku failed: %w", err)
	}

	key := fmt.Sprintf("sku:%d", sku.SkuID)
	if us.cache.IsExist(ctx, key) {
		err = us.cache.DeleteSku(ctx, key)
		if err != nil {
			return fmt.Errorf("usecase.DeleteSku failed: %w", err)
		}
	}

	return nil
}

func (us *useCase) ViewSkuImage(ctx context.Context, sku *model.Sku, pageNum *int64, pageSize *int64) (images []*model.SkuImage, err error) {
	pNum, pSize := us.svc.NormalizePagination(pageNum, pageSize)

	if pNum < 1 || pSize < 1 {
		return nil, fmt.Errorf("usecase.ViewSkuImage failed: invalid PageNum or PageSize")
	}
	offset := pNum * pSize
	key := fmt.Sprintf("skuImgs:%d:%d", sku.SkuID, offset)
	if us.cache.IsExist(ctx, key) {
		ret, err := us.cache.GetSkuImages(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetSkuImages failed: %w", err)
		}
		return ret, nil
	}

	images, err = us.db.ViewSkuImage(ctx, sku, pNum, pSize)
	if err != nil {
		return nil, fmt.Errorf("usecase.ViewSkuImage failed: %w", err)
	}

	us.cache.SetSkuImages(ctx, key, images)

	return images, nil
}

func (us *useCase) ViewSku(ctx context.Context, sku *model.Sku, pageNum *int64, pageSize *int64, isSpuId bool) (skus []*model.Sku, err error) {
	pNum, pSize := us.svc.NormalizePagination(pageNum, pageSize)
	if pNum < 1 || pSize < 1 {
		return nil, fmt.Errorf("usecase.ViewSku failed: invalid PageNum or PageSize")
	}

	var skuIDs []*int64
	if isSpuId {
		ids, err := us.svc.GetSkuIdBySpuID(ctx, sku.SpuID, pNum, pSize)
		if err != nil {
			return nil, fmt.Errorf("usecase.ViewSku failed: %w", err)
		}
		skuIDs = ids
	} else {
		skuIDs = []*int64{&sku.SkuID}
	}

	var remainingIDs []*int64
	for _, id := range skuIDs {
		key := fmt.Sprintf("sku:%d", *id)
		if us.cache.IsExist(ctx, key) {
			s, err := us.cache.GetSku(ctx, key)
			if err != nil {
				return nil, fmt.Errorf("usecase.ViewSku failed: %w", err)
			}
			skus = append(skus, s)
		} else {
			remainingIDs = append(remainingIDs, id)
		}
	}
	if len(remainingIDs) == 0 {
		return skus, nil
	}

	skuIDs = remainingIDs

	result, err := us.db.ViewSku(ctx, skuIDs, pNum, pSize)
	if err != nil {
		return nil, fmt.Errorf("usecase.ViewSku failed: %w", err)
	}

	for _, s := range result {
		key := fmt.Sprintf("sku:%d", s.SkuID)
		us.cache.SetSku(ctx, key, s)
	}

	skus = append(skus, result...)
	return skus, nil
}

func (us *useCase) UploadSkuAttr(ctx context.Context, attr *model.AttrValue, sku *model.Sku) (err error) {
	key := fmt.Sprintf("sku:%d", sku.SkuID)

	ret, err := us.db.GetSkuBySkuId(ctx, sku.SkuID)
	if err != nil {
		return fmt.Errorf("service.UpdateSku: get sku by sku id failed: %w", err)
	}

	if err := us.svc.IdentifyUser(ctx, ret.CreatorID); err != nil {
		return fmt.Errorf("service.UpdateSku: %w", err)
	}

	if err = us.db.UploadSkuAttr(ctx, sku, attr); err != nil {
		return fmt.Errorf("usecase.UploadSkuAttr failed: %w", err)
	}

	if us.cache.IsExist(ctx, key) {
		ret, err := us.cache.GetSku(ctx, key)
		if err != nil {
			return fmt.Errorf("usecase.UploadSkuAttr failed: %w", err)
		}
		ret.SaleAttr = append(ret.SaleAttr, attr)
		us.cache.SetSku(ctx, key, ret)
	}

	return nil
}

func (us *useCase) ListSkuInfo(ctx context.Context, ids []int64, pageNum int64, pageSize int64) (skuInfos []*model.Sku, err error) {
	if pageNum < 1 || pageSize < 1 {
		return nil, fmt.Errorf("usecase.ListSkuInfo failed: invalid PageNum or PageSize")
	}

	var remainingIDs []int64
	for i := (pageNum - 1) * pageSize; i < pageNum*pageSize; i++ {
		key := fmt.Sprintf("sku:%d", ids[i])
		if us.cache.IsExist(ctx, key) {
			s, err := us.cache.GetSku(ctx, key)
			if err != nil {
				return nil, fmt.Errorf("usecase.ListSkuInfo failed: %w", err)
			}
			skuInfos = append(skuInfos, s)
		} else {
			remainingIDs = append(remainingIDs, ids[i])
		}
	}

	if len(remainingIDs) == 0 {
		return skuInfos, nil
	}

	ids = remainingIDs

	result, err := us.db.ListSkuInfo(ctx, ids, int(pageNum), int(pageSize))
	if err != nil {
		return nil, fmt.Errorf("usecase.ListSkuInfo failed: %w", err)
	}

	for _, s := range result {
		key := fmt.Sprintf("sku:%d", s.SkuID)
		us.cache.SetSku(ctx, key, s)
	}

	skuInfos = append(skuInfos, result...)

	return skuInfos, nil
}
