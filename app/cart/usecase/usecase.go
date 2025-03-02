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

	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/app/cart/domain/repository"
	"github.com/west2-online/DomTok/app/cart/domain/service"
)

type CartCasePort interface {
	AddGoodsIntoCart(ctx context.Context, goods *model.GoodInfo) error
	ShowCartGoods(ctx context.Context, pageNum int64) ([]*model.CartGoods, error)
}

type UseCase struct {
	DB    repository.PersistencePort
	Cache repository.CachePort
	MQ    repository.MqPort
	svc   *service.CartService
}

func NewCartCase(db repository.PersistencePort, cache repository.CachePort, mq repository.MqPort, svc *service.CartService) *UseCase {
	return &UseCase{
		DB:    db,
		Cache: cache,
		MQ:    mq,
		svc:   svc,
	}
}
