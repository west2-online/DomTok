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

package redis

import (
	"math/rand/v2"
	"testing"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
)

func buildSpuImage(t *testing.T, spuId int64) *model.SpuImages {
	imgs := []*model.SpuImage{
		{
			ImageID: rand.Int64(),
			Url:     "http://example/com",
			SpuID:   spuId,
		},
		{
			ImageID: rand.Int64(),
			Url:     "http://example/com",
			SpuID:   spuId,
		},
		{
			ImageID: rand.Int64(),
			Url:     "http://example/com",
			SpuID:   spuId,
		},
	}
	return &model.SpuImages{
		Images: imgs,
		Total:  int64(len(imgs)),
	}
}
