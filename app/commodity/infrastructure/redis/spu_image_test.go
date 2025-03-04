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
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func initTest(t *testing.T) repository.CommodityCache {
	t.Helper()
	config.Init("test-cache")
	logger.Ignore()
	re, err := client.InitRedis(constants.RedisDBCommodity)
	if err != nil {
		panic(err)
	}
	return NewCommodityCache(re)
}

func TestCommodityCache_SetAndGet(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	cache := initTest(t)

	offset := 0
	var spuId int64 = rand.Int64()
	ctx := context.Background()
	images := buildSpuImage(t, spuId)
	key := fmt.Sprintf("spuImgs:%d:%d", spuId, offset)

	Convey("Test cache set and get", t, func() {
		var err error
		imgs := new(model.SpuImages)
		imgs.Images = make([]*model.SpuImage, 0)

		imgs, err = cache.GetSpuImages(ctx, key)
		So(err, ShouldNotBeNil)
		So(imgs, ShouldBeNil)

		cache.SetSpuImages(ctx, key, images)

		imgs, err = cache.GetSpuImages(ctx, key)
		So(err, ShouldBeNil)
		So(len(imgs.Images), ShouldEqual, len(images.Images))
		So(imgs.Total, ShouldEqual, images.Total)
	})
}
