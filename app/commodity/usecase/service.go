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
	loginData, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpu failed: %w", err)
	}
	spu.CreatorId = loginData

	id, err = us.svc.CreateSpu(ctx, spu)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpu failed: %w", err)
	}
	return id, nil
}

func (us *useCase) CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) (int64, error) {
	id, err := us.svc.CreateSpuImage(ctx, spuImage)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpuImage failed: %w", err)
	}
	return id, nil
}

func (us *useCase) DeleteSpu(ctx context.Context, spuId int64) error {
	exists, err := us.db.IsExistSku(ctx, spuId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}
	if exists {
		return fmt.Errorf("usecase.DeleteSpu failed: spu-%d‘s sku already exists", spuId)
	}

	ret, err := us.db.GetSpuBySpuId(ctx, spuId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}

	err = us.svc.IdentifyUser(ctx, ret.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpu identify user failed: %w", err)
	}

	if err = us.db.DeleteSpu(ctx, spuId); err != nil {
		return fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}

	return nil
}

func (us *useCase) UpdateSpu(ctx context.Context, spu *model.Spu) error {
	ret, err := us.db.GetSpuBySpuId(ctx, spu.SpuId)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpu failed: %w", err)
	}

	err = us.svc.IdentifyUser(ctx, ret.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpu identify user failed: %w", err)
	}

	spu.GoodsHeadDrawingUrl = utils.GenerateFileName(constants.SpuDirDest, spu.SpuId)
	if err = us.svc.UpdateSpu(ctx, spu, ret); err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}
	return nil
}

func (us *useCase) UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error {
	spu, err := us.svc.GetSpuFromImageId(ctx, spuImage.ImageID)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage failed: %w", err)
	}

	err = us.svc.IdentifyUser(ctx, spu.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage identify user failed: %w", err)
	}

	if err = us.svc.UpdateSpuImage(ctx, spuImage); err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}
	return nil
}

func (us *useCase) DeleteSpuImage(ctx context.Context, imageId int64) error {
	spu, err := us.svc.GetSpuFromImageId(ctx, imageId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage failed: %w", err)
	}

	err = us.svc.IdentifyUser(ctx, spu.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage identify user failed: %w", err)
	}

	if err = us.svc.DeleteSpuImage(ctx, imageId); err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage failed: %w", err)
	}
	return nil
}

func (us *useCase) CreateSku(ctx context.Context, sku *model.Sku) (skuID int64, err error) {
	if err = us.svc.SetCreatorID(ctx, sku); err != nil {
		return 0, fmt.Errorf("usecase.CreateSku failed: %w", err)
	}

	skuID, err = us.svc.CreateSku(ctx, sku)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSku failed: %w", err)
	}
	return skuID, nil
}

func (us *useCase) UpdateSku(ctx context.Context, sku *model.Sku) (err error) {
	ret, err := us.db.GetSkuBySkuId(ctx, sku.SkuID)
	if err != nil {
		return fmt.Errorf("service.UpdateSku: get sku by sku id failed: %w", err)
	}

	if err := us.svc.IdentifyUser(ctx, ret.CreatorID); err != nil {
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

	return nil
}

func (us *useCase) ViewSkuImage(ctx context.Context, sku *model.Sku, pageNum *int64, pageSize *int64) (images []*model.SkuImage, err error) {
	pNum, pSize := us.svc.NormalizePagination(pageNum, pageSize)

	if pNum < 1 || pSize < 1 {
		return nil, fmt.Errorf("usecase.ViewSkuImage failed: invalid PageNum or PageSize")
	}

	images, err = us.db.ViewSkuImage(ctx, sku, pNum, pSize)
	if err != nil {
		return nil, fmt.Errorf("usecase.ViewSkuImage failed: %w", err)
	}

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

	skus, err = us.db.ViewSku(ctx, skuIDs, pNum, pSize)
	if err != nil {
		return nil, fmt.Errorf("usecase.ViewSku failed: %w", err)
	}
	return skus, nil
}

func (us *useCase) UploadSkuAttr(ctx context.Context, attr *model.AttrValue, sku *model.Sku) (err error) {
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
	return nil
}

func (us *useCase) ListSkuInfo(ctx context.Context, ids []int64, pageNum int64, pageSize int64) (skuInfos []*model.Sku, err error) {
	if pageNum < 1 || pageSize < 1 {
		return nil, fmt.Errorf("usecase.ListSkuInfo failed: invalid PageNum or PageSize")
	}

	skuInfos, err = us.db.ListSkuInfo(ctx, ids, int(pageNum), int(pageSize))
	if err != nil {
		return nil, fmt.Errorf("usecase.ListSkuInfo failed: %w", err)
	}

	return skuInfos, nil
}
