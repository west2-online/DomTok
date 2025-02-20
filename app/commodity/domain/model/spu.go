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

<<<<<<<< HEAD:app/commodity/infrastructure/mysql/model.go
package mysql

import (
	"time"

	"github.com/west2-online/DomTok/pkg/constants"
)

type Category struct {
	Id        int64
	Name      string
	CreatorId int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	// gorm.Model
========
package model

type Spu struct {
	SpuId               int64
	Name                string
	CreatorId           int64
	Description         string
	CategoryId          int64
	GoodsHeadDrawing    []byte
	Price               float64
	ForSale             int
	Shipping            float64
	CreatedAt           int64
	UpdatedAt           int64
	DeletedAt           int64
	GoodsHeadDrawingUrl string
>>>>>>>> upstream/main:app/commodity/domain/model/spu.go
}

func (Category) TableName() string {
	return constants.CategoryTableName
}
