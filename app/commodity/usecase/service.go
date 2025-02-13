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
