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
	kmodel "github.com/west2-online/DomTok/kitex_gen/model"
	kcontext "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (uc *useCase) CreateCategory(ctx context.Context, category *entities.Category) (int64, error) {
	exist, err := uc.db.IsCategoryExist(ctx, category.Id)
	if err != nil {
		return 0, fmt.Errorf("check category exist failed: %w", err)
	}
	if exist {
		return 0, errno.NewErrNo(errno.ServiceUserNotExist, "category  exist")
	}

	if err = uc.svc.CreateCategory(ctx, category); err != nil {
		return 0, fmt.Errorf("create category failed: %w", err)
	}

	return category.Id, nil
}

func (uc *useCase) DeleteCategory(ctx context.Context, category *entities.Category) (err error) {
	// 判断是否存在
	exist, err := uc.db.IsCategoryExist(ctx, category.Id)
	if err != nil {
		return fmt.Errorf("check category exist failed: %w", err)
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceUserNotExist, "category does not exist")
	}
	// 判断用户是否有权限
	LoginData, err := kcontext.GetLoginData(ctx)
	if err != nil {
		return errno.NewErrNo(errno.AuthInvalidCode, " Get login data fail")
	}
	category.CreatorId = uc.db.CategoryCreatorId(ctx, category.Id)
	if LoginData.UserId != category.CreatorId {
		return errno.NewErrNo(errno.AuthNoOperatePermissionCode, " You are not authorized to delete this category")
	}
	err = uc.db.DeleteCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("delete category failed: %w", err)
	}
	return nil
}

func (uc *useCase) UpdateCategory(ctx context.Context, category *entities.Category) (err error) {
	// 判断是否存在
	exist, err := uc.db.IsCategoryExist(ctx, category.Id)
	if err != nil {
		return fmt.Errorf("check category exist failed: %w", err)
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceUserNotExist, "category does not exist")
	}
	// 判断用户是否有权限
	LoginData, err := kcontext.GetLoginData(ctx)
	if err != nil {
		return errno.NewErrNo(errno.AuthInvalidCode, " Get login data fail")
	}
	category.CreatorId = uc.db.CategoryCreatorId(ctx, category.Id)
	if LoginData.UserId != category.CreatorId {
		return errno.NewErrNo(errno.AuthNoOperatePermissionCode, " You are not authorized to delete this category")
	}
	err = uc.db.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("update category failed: %w", err)
	}
	return err
}

func (uc *useCase) ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*kmodel.CategoryInfo, err error) {
	resp, err = uc.db.ViewCategory(ctx, pageNum, pageSize)
	if err != nil {
		return nil, errno.Errorf(errno.ServiceListCategoryFailed, "failed to view categories: %v", err)
	}
	return resp, nil
}
