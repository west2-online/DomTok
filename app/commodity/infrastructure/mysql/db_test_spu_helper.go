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

package mysql

import (
	"math/rand/v2"
	"testing"

	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
)

func buildTestSpu(t *testing.T, creatorId int64) *model.Spu {
	t.Helper()
	return &model.Spu{
		SpuId:               rand.Int64(),
		Name:                "OPPO",
		CreatorId:           creatorId, // TODO,
		Description:         "desc",
		CategoryId:          rand.Int64(),
		Price:               rand.Float64(),
		GoodsHeadDrawingUrl: "http://example.com",
	}
}

func buildTestSpuImage(t *testing.T, spuId int64) *model.SpuImage {
	t.Helper()
	return &model.SpuImage{
		ImageID: rand.Int64(),
		SpuID:   spuId,
		Url:     "http://example.com",
	}
}

func matchTestSpu(t *testing.T, mock *model.Spu, expected *model.Spu) {
	t.Helper()
	var exp int32 = 4
	So(mock.SpuId, ShouldEqual, expected.SpuId)
	So(mock.CreatorId, ShouldEqual, expected.CreatorId)
	So(mock.Description, ShouldEqual, expected.Description)
	So(mock.CategoryId, ShouldEqual, expected.CategoryId)
	So(decimal.NewFromFloat(mock.Price).Equal(decimal.NewFromFloat(expected.Price).Round(exp)), ShouldBeTrue)
	So(mock.Name, ShouldEqual, expected.Name)
	So(mock.ForSale, ShouldEqual, expected.ForSale)
	So(decimal.NewFromFloat(mock.Shipping).Equal(decimal.NewFromFloat(expected.Shipping).Round(exp)), ShouldBeTrue)
	So(mock.GoodsHeadDrawingUrl, ShouldEqual, expected.GoodsHeadDrawingUrl)
}

func matchTestSpuImage(t *testing.T, mock *model.SpuImage, expected *model.SpuImage) {
	t.Helper()
	So(mock.SpuID, ShouldEqual, expected.SpuID)
	So(mock.Url, ShouldEqual, expected.Url)
	So(mock.ImageID, ShouldEqual, expected.ImageID)
}
