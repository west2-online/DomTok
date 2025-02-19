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
	ctx = contextLogin.WithLoginData(ctx, 1) //TODO: delete

	loginData, err := contextLogin.GetLoginData(ctx)
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
	ctx = contextLogin.WithLoginData(ctx, 1) //TODO: delete
	id, err := us.svc.CreateSpuImage(ctx, spuImage)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpuImage failed: %w", err)
	}
	return id, nil
}

func (us *useCase) DeleteSpu(ctx context.Context, spuId int64) error {
	ctx = contextLogin.WithLoginData(ctx, 1) //TODO: delete
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
		return fmt.Errorf("usecase.DeleteSpuImages failed: %w", err)
	}

	return nil
}

func (us *useCase) UpdateSpu(ctx context.Context, spu *model.Spu) error {
	ctx = contextLogin.WithLoginData(ctx, 1) //TODO: delete
	ret, err := us.db.GetSpuBySpuId(ctx, spu.SpuId)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpu failed: %w", err)
	}

	err = us.svc.IdentifyUser(ctx, ret.CreatorId)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpu identify user failed: %w", err)
	}

	if err = us.svc.Verify(us.svc.VerifyForSaleStatus(spu.ForSale)); err != nil {
		return fmt.Errorf("usecase.UpdateSpu verify failed: %w", err)
	}

	spu.GoodsHeadDrawingUrl = utils.GenerateFileName(constants.SpuDirDest, spu.SpuId)
	if err = us.svc.UpdateSpu(ctx, spu, ret); err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}
	return nil
}

func (us *useCase) UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error {
	ctx = contextLogin.WithLoginData(ctx, 1) //TODO: delete
	spu, img, err := us.svc.GetSpuFromImageId(ctx, spuImage.ImageID)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpuImage failed: %w", err)
	}

	err = us.svc.IdentifyUser(ctx, spu.CreatorId)
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
	ctx = contextLogin.WithLoginData(ctx, 1) //TODO: delete
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
