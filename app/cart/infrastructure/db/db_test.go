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

package db

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/cart/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

var _db repository.PersistencePort

func initDB() {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	_db = NewDBAdapter(gormDB)
}

func initConfig() bool {
	if !utils.EnvironmentEnable() {
		return false
	}
	logger.Ignore()
	config.Init("cart-db-test")
	initDB()
	return true
}

func TestDBAdapter_Cart(t *testing.T) {
	// 如果环境或配置不满足，则直接跳过
	if !initConfig() {
		return
	}

	// 上下文，避免测试中频繁创建
	ctx := context.Background()

	// 这里随便定义一个测试用的 userID
	userID := int64(12345)

	// 这里是你的购物车 JSON，实际可根据需求来定义
	initialSkuJson := `{"items": [{"skuId": 101, "count": 2}]}`

	Convey("测试 Cart 相关数据库操作", t, func() {
		Convey("创建购物车(CreateCart)", func() {
			err := _db.CreateCart(ctx, userID, initialSkuJson)
			So(err, ShouldBeNil)
		})

		Convey("查询购物车(GetCartByUserId)", func() {
			exists, cart, err := _db.GetCartByUserId(ctx, userID)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
			So(cart, ShouldNotBeNil)

			// 验证字段是否和刚插入的相同
			So(cart.UserId, ShouldEqual, userID)
			So(cart.SkuJson, ShouldEqual, initialSkuJson)
		})

		Convey("更新(保存)购物车(SaveCart)", func() {
			newSkuJson := `{"items": [{"skuId": 101, "count": 3}, {"skuId": 102, "count": 1}]}`
			err := _db.SaveCart(ctx, userID, newSkuJson)
			So(err, ShouldBeNil)

			// 再次查询，确认已更新
			exists, cart, err := _db.GetCartByUserId(ctx, userID)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
			So(cart.SkuJson, ShouldEqual, newSkuJson)
		})

		Convey("删除购物车(DeleteCart)", func() {
			err := _db.DeleteCart(ctx, userID)
			So(err, ShouldBeNil)

			exists, cart, err := _db.GetCartByUserId(ctx, userID)
			So(err, ShouldBeNil)
			// 由于删除了，故应该不存在
			So(exists, ShouldBeFalse)
			So(cart, ShouldBeNil)
		})
	})
}
