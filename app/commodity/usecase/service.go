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
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	modelKitex "github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/errno"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
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

func (uc *useCase) ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*modelKitex.CategoryInfo, err error) {
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

func (us *useCase) IncrLockStock(ctx context.Context, infos []*modelKitex.SkuBuyInfo) error {
	if !us.svc.Cached(ctx, infos) {
		return errno.Errorf(errno.RedisKeyNotExist, "useCase.IncrLockStock failed")
	}
	err := us.cache.IncrLockStockNum(ctx, infos)
	if err != nil {
		return fmt.Errorf("usecase.IncrLockStock failed: %w", err)
	}
	err = us.db.IncrLockStock(ctx, infos)
	if err != nil {
		return fmt.Errorf("usecase.IncrLockStock failed: %w", err)
	}
	return nil
}

func (us *useCase) DecrLockStock(ctx context.Context, infos []*modelKitex.SkuBuyInfo) error {
	if !us.svc.Cached(ctx, infos) {
		return errno.Errorf(errno.RedisKeyNotExist, "useCase.DecrLockStock failed")
	}
	err := us.cache.DecrLockStockNum(ctx, infos)
	if err != nil {
		return fmt.Errorf("usecase.IncrLockStock failed: %w", err)
	}

	// TODO: mq配置db后续处理
	err = us.db.DecrLockStock(ctx, infos)
	if err != nil {
		return fmt.Errorf("usecase.IncrLockStock failed: %w", err)
	}
	return err
}

func (us *useCase) DecrStock(ctx context.Context, infos []*modelKitex.SkuBuyInfo) error {
	return us.db.DecrStock(ctx, infos)
}
