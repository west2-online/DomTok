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

<<<<<<<< HEAD:app/commodity/infrastructure/mysql/db.go
package mysql
========
package redis
>>>>>>>> upstream/main:app/commodity/infrastructure/redis/ini.go

import (
	"github.com/redis/go-redis/v9"

<<<<<<<< HEAD:app/commodity/infrastructure/mysql/db.go
	"gorm.io/gorm"

	_ "github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
)

// commodityDB impl domain.CommodityDB defined domain
type commodityDB struct {
	client *gorm.DB
}

func NewCommodityDB(client *gorm.DB) repository.CommodityDB {
	return &commodityDB{client: client}
}

func (db *commodityDB) CreateCategory(ctx context.Context, name string) error {
	return nil
========
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
)

type commodityCache struct {
	client *redis.Client
}

func NewCommodityCache(client *redis.Client) repository.CommodityCache {
	return &commodityCache{client: client}
>>>>>>>> upstream/main:app/commodity/infrastructure/redis/ini.go
}
