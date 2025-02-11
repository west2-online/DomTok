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
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/pkg/utils"
)

type CommodityService struct {
	db    repository.CommodityDB
	sf    *utils.Snowflake
	cache repository.CommodityCache
	// TODO
}

func NewCommodityService(db repository.CommodityDB, sf *utils.Snowflake, cache repository.CommodityCache) *CommodityService {
	if db == nil {
		panic("commodityService's db should not be nil")
	}
	if sf == nil {
		panic("commodityService's snowflake should not be nil")
	}

	if cache == nil {
		panic("commodityService's cache should not be nil")
	}

	svc := &CommodityService{
		db:    db,
		sf:    sf,
		cache: cache,
	}
	svc.init()
	return svc
}

func (s *CommodityService) init() {
}
