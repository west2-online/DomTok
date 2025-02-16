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

package service

import (
	"context"

	"github.com/west2-online/DomTok/app/cart/domain/repository"
)

type CartService struct {
	DB    repository.PersistencePort
	Cache repository.CachePort
	MQ    repository.MqPort
}

func NewCartService(db repository.PersistencePort, cache repository.CachePort, mq repository.MqPort) *CartService {
	svc := &CartService{
		DB:    db,
		Cache: cache,
		MQ:    mq,
	}
	svc.init()
	return svc
}

func (svc *CartService) init() {
	svc.initConsumer()
}

func (svc *CartService) initConsumer() {
	go svc.ConsumeAddGoods(context.Background())
}
