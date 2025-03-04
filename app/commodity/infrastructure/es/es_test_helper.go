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

package es

import (
	"math/rand/v2"
	"testing"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
)

func buildTestSpu(t *testing.T, creatorId int64) *model.Spu {
	t.Helper()
	return &model.Spu{
		SpuId:               rand.Int64(),
		Name:                "我现在在这里",
		CreatorId:           creatorId,
		Description:         "desc",
		CategoryId:          rand.Int64(),
		Price:               rand.Float64(),
		GoodsHeadDrawingUrl: "http://example.com",
	}
}
