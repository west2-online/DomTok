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
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

var _db repository.CommodityDB

func initDB() {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	_db = NewCommodityDB(gormDB)
}

func initConfig() bool {
	if !utils.EnvironmentEnable() {
		return false
	}
	logger.Ignore()
	config.Init("commodity-db-spu-test")
	initDB()
	return true
}

func TestCommodityDB_Spu(t *testing.T) {

	if !initConfig() {
		return
	}
	ctx := context.Background()
	limit := 10
	offset := 0
	var uid int64 = 10000

	spuInfo := buildTestSpu(t, uid)
	imgs := make([]*model.SpuImage, 10)
	for i := 0; i < len(imgs); i++ {
		imgs[i] = buildTestSpuImage(t, spuInfo.SpuId)
	}

	Convey("TestCommodityDB_Spu", t, func() {
		Convey("TestCommodityDB_CreateSpu", func() {
			err := _db.CreateSpu(ctx, spuInfo)
			So(err, ShouldBeNil)

			for i := 0; i < len(imgs); i++ {
				err = _db.CreateSpuImage(ctx, imgs[i])
				So(err, ShouldBeNil)
			}

			err = _db.UpdateSpu(ctx, spuInfo)
			So(err, ShouldBeNil)

			err = _db.UpdateSpuImage(ctx, imgs[0])
			So(err, ShouldBeNil)

			ret, err := _db.GetSpuBySpuId(ctx, spuInfo.SpuId)
			So(err, ShouldBeNil)
			matchTestSpu(t, ret, spuInfo)

			for i := 0; i < len(imgs); i++ {
				img, err := _db.GetSpuImage(ctx, imgs[i].ImageID)
				So(err, ShouldBeNil)
				matchTestSpuImage(t, img, imgs[i])
			}

			_, total, err := _db.GetImagesBySpuId(ctx, spuInfo.SpuId, offset, limit)
			So(err, ShouldBeNil)
			So(len(imgs), ShouldEqual, int(total))

			spus, err := _db.GetSpuByIds(ctx, []int64{spuInfo.SpuId})
			So(err, ShouldBeNil)
			for _, s := range spus {
				matchTestSpu(t, s, spuInfo)
			}
		})

		Convey("TestCommodityDB_CreateSpuAgain", func() {
			err := _db.CreateSpu(ctx, spuInfo)
			So(err, ShouldNotBeNil)

			err = _db.CreateSpuImage(ctx, imgs[0])
			So(err, ShouldNotBeNil)
		})

		Convey("TestCommodityDB_DeleteSpu", func() {
			err := _db.DeleteSpu(ctx, spuInfo.SpuId)
			So(err, ShouldBeNil)

			err = _db.DeleteSpuImage(ctx, imgs[0].ImageID)
			So(err, ShouldBeNil)

			_, url, err := _db.DeleteSpuImagesBySpuId(ctx, spuInfo.SpuId)
			So(err, ShouldBeNil)
			So(len(url), ShouldEqual, len(imgs)-1)
		})
	})
}
