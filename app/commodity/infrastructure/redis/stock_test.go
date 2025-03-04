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
	"math/rand/v2"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/pkg/utils"
)

func TestCommodityCache_Stock(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	cache := initTest(t)

	var stockNum int64 = 50
	var lock int64 = 15

	id := rand.Int64()
	ctx := context.Background()
	key := cache.GetStockKey(id)

	lockKey := cache.GetLockStockKey(id)
	infos := buildSkuBuyInfo(t, id)

	Convey("TestCommodityCache_Stock", t, func() {
		Convey("TestCommodityCache_GetAndSetStock", func() {
			stock, err := cache.GetLockStockNum(ctx, key)
			So(err, ShouldNotBeNil)
			So(stock, ShouldEqual, 0)

			cache.SetLockStockNum(ctx, key, stockNum)

			stock, err = cache.GetLockStockNum(ctx, key)
			So(err, ShouldBeNil)
			So(stock, ShouldEqual, stockNum)

			lockstock, err := cache.GetLockStockNum(ctx, lockKey)
			So(err, ShouldNotBeNil)
			So(lockstock, ShouldEqual, 0)

			cache.SetLockStockNum(ctx, lockKey, lock)

			lockstock, err = cache.GetLockStockNum(ctx, lockKey)
			So(err, ShouldBeNil)
			So(lockstock, ShouldEqual, lock)
		})
		Convey("TestCommodityCache_IncrAndDecrLockStock", func() {
			err := cache.DecrLockStockNum(ctx, infos)
			So(err, ShouldBeNil)

			lockStock, err := cache.GetLockStockNum(ctx, lockKey)
			So(err, ShouldBeNil)
			So(lockStock, ShouldEqual, lock-infos[0].Count)

			err = cache.IncrLockStockNum(ctx, infos)
			So(err, ShouldBeNil)

			lockStock, err = cache.GetLockStockNum(ctx, lockKey)
			So(err, ShouldBeNil)
			So(lockStock, ShouldEqual, lock)
		})

		Convey("TestCommodityCache_DecrStockNum", func() {
			err := cache.DecrStockNum(ctx, infos)
			So(err, ShouldBeNil)

			stock, err := cache.GetLockStockNum(ctx, key)
			So(err, ShouldBeNil)
			So(stock, ShouldEqual, stockNum-infos[0].Count)

			lockstock, err := cache.GetLockStockNum(ctx, lockKey)
			So(err, ShouldBeNil)
			So(lockstock, ShouldEqual, lock-infos[0].Count)
		})
	})
}
