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
	"github.com/west2-online/DomTok/pkg/errno"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
)

func (us *useCase) CreateCategory(ctx context.Context, category *model.Category) (id int64, err error) {
	return 0, nil
}

func (us *useCase) CreateSpu(ctx context.Context, spu *model.Spu) (id int64, err error) {
	if err = us.svc.Verify(us.svc.VerifyFileType(spu.GoodsHeadDrawingName)); err != nil {
		return 0, fmt.Errorf("usecase.CreateSpu failed: %w", err)
	}

	loginData, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpu failed: %w", err)
	}
	spu.CreatorId = loginData.UserId

	id, err = us.svc.CreateSpu(ctx, spu)
	if err != nil {
		return 0, fmt.Errorf("usecase.CreateSpu failed: %w", err)
	}
	return id, nil
}

func (us *useCase) CreateSpuImage(ctx context.Context, spuImage *model.SpuImage) (int64, error) {
	if err := us.svc.Verify(us.svc.VerifyFileType(spuImage.Url)); err != nil {
		return 0, fmt.Errorf("usecase.CreateSpuImage failed: %w", err)
	}

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

	loginData, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("usecase.DeleteSpu failed: %w", err)
	}

	if loginData.UserId != ret.CreatorId {
		return errno.AuthNoOperatePermission
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
	loginData, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpu failed: %w", err)
	}
	if loginData.UserId != ret.CreatorId {
		return errno.AuthNoOperatePermission
	}
	if err = us.svc.UpdateSpu(ctx, spu); err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}
	return nil
}

func (us *useCase) UpdateSpuImage(ctx context.Context, spuImage *model.SpuImage) error {
	image, err := us.db.GetSpuImage(ctx, spuImage.ImageID)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}

	spu, err := us.db.GetSpuBySpuId(ctx, image.SpuID)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}
	// 上面这部分封装一下
	loginData, err := contextLogin.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}
	if loginData.UserId != spu.CreatorId {
		return errno.AuthNoOperatePermission
	}
	if err = us.svc.UpdateSpuImage(ctx, image); err != nil {
		return fmt.Errorf("usecase.UpdateSpuImage failed: %w", err)
	}
	return nil
}
