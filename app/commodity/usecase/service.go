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
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/utils"
)

func (uc *useCase) CreateCategory(ctx context.Context, category *model.Category) (int64, error) {
	exist, err := uc.db.IsCategoryExistByName(ctx, category.Name)
	if err != nil {
		return 0, fmt.Errorf("check category exist failed: %w", err)
	}
	if exist {
		return 0, errno.NewErrNo(errno.ServiceUserExist, "category  exist")
	}

	if err = uc.svc.CreateCategory(ctx, category); err != nil {
		return 0, fmt.Errorf("create category failed: %w", err)
	}

	return category.Id, nil
}

func (uc *useCase) DeleteCategory(ctx context.Context, category *model.Category) (err error) {
	// 判断是否存在
	exist, err := uc.db.IsCategoryExistById(ctx, category.Id)
	if err != nil {
		return fmt.Errorf("check category exist failed: %w", err)
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceUserNotExist, "category does not exist")
	}
	// 判断用户是否有权限
	err = uc.svc.IdentifyUser(ctx, category.CreatorId)
	if err != nil {
		return errno.NewErrNo(errno.AuthInvalidCode, " Get login data fail")
	}
	err = uc.db.DeleteCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("delete category failed: %w", err)
	}
	return nil
}

func (uc *useCase) UpdateCategory(ctx context.Context, category *model.Category) (err error) {
	// 判断是否存在
	exist, err := uc.db.IsCategoryExistById(ctx, category.Id)
	if err != nil {
		return fmt.Errorf("check category exist failed: %w", err)
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceUserNotExist, "category does not exist")
	}
	// 判断用户是否有权限
	err = uc.svc.IdentifyUser(ctx, category.CreatorId)
	if err != nil {
		return errno.NewErrNo(errno.AuthInvalidCode, " Get login data fail")
	}
	err = uc.db.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("update category failed: %w", err)
	}
	return err
}

func (uc *useCase) ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*model.CategoryInfo, err error) {
	resp, err = uc.db.ViewCategory(ctx, pageNum, pageSize)
	if err != nil {
		return nil, errno.Errorf(errno.ServiceListCategoryFailed, "failed to view categories: %v", err)
	}
	return resp, nil
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
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
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

func (us *useCase) ViewSpus(ctx context.Context, req *commodity.ViewSpuReq) ([]*model.Spu, int64, error) {
	ids, total, err := us.es.SearchItems(ctx, constants.SpuTableName, req)
	if err != nil {
		return nil, 0, fmt.Errorf("usecase.ViewSpus failed: %w", err)
	}

	res, err := us.db.GetSpuByIds(ctx, ids)
	if err != nil {
		return nil, 0, fmt.Errorf("usecase.ViewSpus failed: %w", err)
	}
	return res, total, err
}

func (us *useCase) ListSpuInfo(ctx context.Context, ids []int64) ([]*model.Spu, error) {
	return us.db.GetSpuByIds(ctx, ids)
}

func (us *useCase) IncrLockStock(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := us.cache.IsHealthy(ctx)
	if err != nil {
		err = us.db.IncrLockStock(ctx, infos)
		if err != nil {
			return fmt.Errorf("usecase.IncrLockStock failed: %w", err)
		}
		return err
	} else {
		return us.svc.IncrLockStockInNX(ctx, infos)
	}
}

func (us *useCase) DecrLockStock(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := us.cache.IsHealthy(ctx)
	if err != nil {
		err = us.db.DecrLockStock(ctx, infos)
		if err != nil {
			return fmt.Errorf("usecase.DecrLockStock failed: %w", err)
		}
		return err
	} else {
		return us.svc.DecrLockStockInNX(ctx, infos)
	}
}

func (us *useCase) DecrStock(ctx context.Context, infos []*model.SkuBuyInfo) error {
	err := us.cache.IsHealthy(ctx)
	if err != nil {
		err = us.db.DecrStock(ctx, infos)
		if err != nil {
			return fmt.Errorf("usecase.DecrStock failed: %w", err)
		}
		return err
	} else {
		return us.svc.DecrStockInNX(ctx, infos)
	}
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

func (us *useCase) UpdateSku(ctx context.Context, sku *model.Sku, ext string) (err error) {
	ret, err := us.db.GetSkuBySkuId(ctx, sku.SkuID)
	if err != nil {
		return fmt.Errorf("service.UpdateSku: get sku by sku id failed: %w", err)
	}

	if err := us.svc.IdentifyUserInStreamCtx(ctx, ret.CreatorID); err != nil {
		return fmt.Errorf("service.UpdateSku: %w", err)
	}

	sku.StyleHeadDrawingUrl = utils.GenerateFileName(constants.SkuDirDest, sku.SkuID) + ext
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

	if err = us.svc.DeleteSku(ctx, sku); err != nil {
		return fmt.Errorf("usecase.DeleteSku failed: %w", err)
	}
	return nil
}

func (us *useCase) ViewSku(ctx context.Context, sku *model.Sku, pageNum *int64, pageSize *int64, isSpuId bool) (skus []*model.Sku, total int64, err error) {
	pNum, pSize := us.svc.NormalizePagination(pageNum, pageSize)
	if pNum < 1 || pSize < 1 {
		return nil, -1, fmt.Errorf("usecase.ViewSku failed: invalid PageNum or PageSize")
	}

	var skuIDs []*int64
	if isSpuId {
		ids, err := us.svc.GetSkuIdBySpuID(ctx, sku.SpuID, pNum, pSize)
		if err != nil {
			return nil, -1, fmt.Errorf("usecase.ViewSku failed: %w", err)
		}
		skuIDs = ids
	} else {
		skuIDs = []*int64{&sku.SkuID}
	}

	skus, total, err = us.svc.ViewSku(ctx, skuIDs, pNum, pSize)
	if err != nil {
		return nil, -1, fmt.Errorf("usecase.ViewSku failed: %w", err)
	}
	return skus, total, nil
}

func (us *useCase) UploadSkuAttr(ctx context.Context, attr *model.AttrValue, sku *model.Sku) (err error) {
	ret, err := us.db.GetSkuBySkuId(ctx, sku.SkuID)
	if err != nil {
		return fmt.Errorf("service.UpdateSku: get sku by sku id failed: %w", err)
	}

	if err := us.svc.IdentifyUser(ctx, ret.CreatorID); err != nil {
		return fmt.Errorf("service.UpdateSku: %w", err)
	}

	sku.HistoryID = ret.HistoryID

	err = us.svc.UploadSkuAttr(ctx, attr, sku)
	if err != nil {
		return fmt.Errorf("usecase.UploadSkuAttr failed: %w", err)
	}

	return nil
}

func (us *useCase) ListSkuInfo(ctx context.Context, ids []int64, pageNum int64, pageSize int64) (skuInfos []*model.Sku, total int64, err error) {
	if pageNum < 1 || pageSize < 1 {
		return nil, -1, fmt.Errorf("usecase.ListSkuInfo failed: invalid PageNum or PageSize")
	}

	skuInfos, total, err = us.svc.ListSkuInfo(ctx, ids, pageNum, pageSize)
	if err != nil {
		return nil, -1, fmt.Errorf("usecase.ListSkuInfo failed: %w", err)
	}
	return skuInfos, total, nil
}

func (us *useCase) CreateSkuImage(ctx context.Context, skuImage *model.SkuImage, data []byte) (int64, error) {
	ret, err := us.db.GetSkuBySkuId(ctx, skuImage.SkuID)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSkuImage failed: %w", err)
	}
	if err := us.svc.IdentifyUserInStreamCtx(ctx, ret.CreatorID); err != nil {
		return 0, fmt.Errorf("usecase.CreateSkuImage failed: %w", err)
	}

	id, err := us.svc.CreateSkuImage(ctx, skuImage, data)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSkuImage failed: %w", err)
	}
	return id, nil
}

func (us *useCase) UpdateSkuImage(ctx context.Context, skuImage *model.SkuImage, data []byte) (err error) {
	sku, img, err := us.svc.GetSkuFromImageId(ctx, skuImage.ImageID)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSkuImage failed: %w", err)
	}

	if err := us.svc.IdentifyUserInStreamCtx(ctx, sku.CreatorID); err != nil {
		return fmt.Errorf("usecase.UpdateSkuImage failed: %w", err)
	}

	skuImage.Url = utils.GenerateFileName(constants.SkuImageDirDest, skuImage.ImageID)
	err = us.svc.UpdateSkuImage(ctx, skuImage, img, data)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSkuImage failed: %w", err)
	}
	return nil
}

func (us *useCase) ViewSkuImages(ctx context.Context, sku *model.Sku, pageNum *int64, pageSize *int64) (images []*model.SkuImage, total int64, err error) {
	pNum, pSize := us.svc.NormalizePagination(pageNum, pageSize)

	if pNum < 1 || pSize < 1 {
		return nil, -1, fmt.Errorf("usecase.ViewSkuImage failed: invalid PageNum or PageSize")
	}

	images, total, err = us.svc.ViewSkuImages(ctx, sku, pNum, pSize)
	if err != nil {
		return nil, -1, fmt.Errorf("usecase.ViewSkuImage failed: %w", err)
	}
	return images, total, nil
}

func (us *useCase) DeleteSkuImage(ctx context.Context, imageId int64) (err error) {
	sku, img, err := us.svc.GetSkuFromImageId(ctx, imageId)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSkuImage failed: %w", err)
	}

	if err := us.svc.IdentifyUser(ctx, sku.CreatorID); err != nil {
		return fmt.Errorf("usecase.DeleteSkuImage failed: %w", err)
	}

	if err = us.svc.DeleteSkuImage(ctx, imageId, img.Url); err != nil {
		return fmt.Errorf("usecase.DeleteSkuImage failed: %w", err)
	}
	return nil
}
