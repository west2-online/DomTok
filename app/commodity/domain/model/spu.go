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
}

// SpuEs : SpuId 和 Category 不能是int64, 存到es里会有精度损失, ref: https://www.cnblogs.com/ahfuzhang/p/16922292.html
type SpuES struct {
	SpuId      string  `json:"spu_id,omitempty"`
	Name       string  `json:"name,omitempty"`
	CategoryId string  `json:"category_id,omitempty"`
	Price      float64 `json:"price,omitempty"`
	Shipping   bool    `json:"shipping,omitempty"`
}
