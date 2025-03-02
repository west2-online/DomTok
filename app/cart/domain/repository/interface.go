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

	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/pkg/kafka"
)

type PersistencePort interface {
	CreateCart(ctx context.Context, uid int64, cart string) error
	GetCartByUserId(ctx context.Context, uid int64) (bool, *model.Cart, error)
	SaveCart(ctx context.Context, uid int64, cart string) error
}

type CachePort interface {
	SetCartCache(ctx context.Context, key string, cart string) error
	GetCartCache(ctx context.Context, key string) (string, error)
	IsKeyExist(ctx context.Context, key string) bool
}

type MqPort interface {
	SendAddGoods(ctx context.Context, uid int64, goods *model.GoodInfo) error
	ConsumeAddGoods(ctx context.Context) <-chan *kafka.Message
}

type RpcPort interface {
	GetGoodsInfo(ctx context.Context, cartGoods []*model.CartGoods) ([]*model.CartGoods, error)
}
