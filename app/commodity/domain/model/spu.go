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

type SpuES struct {
	SpuId      int64
	Name       string
	CategoryId int64
	Price      float64
	Shipping   bool
}

// ConvertIntToBool 将输入的int转换为bool，1为出售true,2为暂不出售false
func (s *Spu) ConvertIntToBool(input int) bool {
	return input == 1
}

func (s *SpuES) ConvertBoolToInt(input bool) int {
	if input {
		return 1
	}
	return 2
}
