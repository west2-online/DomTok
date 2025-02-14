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

package repository

import (
	"context"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	kmodel "github.com/west2-online/DomTok/kitex_gen/model"
)

type CommodityDB interface {
	IsCategoryExist(ctx context.Context, Id int64) (bool, error)
	CategoryCreatorId(ctx context.Context, Id int64) (int64)
	CreateCategory(ctx context.Context, entity *model.Category) error
	DeleteCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	ViewCategory(ctx context.Context, pageNum, pageSize int) (resp []*kmodel.CategoryInfo, err error)
}

type CommodityCache interface{}
