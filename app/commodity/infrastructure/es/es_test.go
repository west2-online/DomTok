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
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

var _es repository.CommodityElastic

func initES() {
	elastic, err := client.NewEsCommodityClient()
	if err != nil {
		panic(err)
	}
	_es = NewCommodityElastic(elastic)
}

func initConfig() bool {
	if !utils.EnvironmentEnable() {
		return false
	}
	logger.Ignore()
	config.Init("es-test")
	initES()
	return true
}

func TestCommodityElastic_CreateAndDeleteItem(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	if !initConfig() {
		return
	}
	ctx := context.Background()
	indexName := constants.SpuTableName
	var creatorId int64 = 10000
	var pageSize int64 = 10
	var pageNum int64 = 0

	infos := make([]*model.Spu, 30)
	for i := range infos {
		infos[i] = buildTestSpu(t, creatorId)
	}

	Convey("TestCommodityElastic_CreateAndDeleteItem", t, func() {
		Convey("TestCommodityElastic_ExistIndex", func() {
			exists := _es.IsExist(ctx, indexName)
			So(exists, ShouldBeFalse)
			err := _es.CreateIndex(ctx, indexName)
			So(err, ShouldBeNil)
			exists = _es.IsExist(ctx, indexName)
			So(exists, ShouldBeTrue)
		})

		Convey("TestCommodityElastic_CreateItem", func() {
			for i := range infos {
				err := _es.AddItem(ctx, indexName, infos[i])
				So(err, ShouldBeNil)
			}

			for i := range infos {
				err := _es.UpdateItem(ctx, indexName, infos[i])
				So(err, ShouldBeNil)
			}

			_, _, err := _es.SearchItems(ctx, indexName, &commodity.ViewSpuReq{
				PageSize: &pageSize,
				PageNum:  &pageNum,
			})
			So(err, ShouldBeNil)

			for _, info := range infos {
				err = _es.RemoveItem(ctx, indexName, info.SpuId)
				So(err, ShouldBeNil)
			}
		})
	})
}
