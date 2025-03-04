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

package cache

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/cart/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

// 如果你和前面 redis 测试类似，也可以直接用 initTest(t) 方式来初始化
// 这里写一个示例 initTestRedis，用于返回一个实现了 SetCartCache / GetCartCache 的缓存适配器。
func initTestRedis(t *testing.T) repository.CachePort {
	t.Helper()
	// 初始化配置信息
	config.Init("test-cache")
	logger.Ignore()

	// 这一步根据你项目里 client.InitRedis() 的实现进行调整
	re, err := client.InitRedis(constants.RedisDBCommodity)
	if err != nil {
		panic(fmt.Sprintf("failed to init redis: %v", err))
	}

	// 假设此处 NewCacheAdapter(re) 返回的结构体，包含 SetCartCache / GetCartCache 方法
	// 如果你的项目里是 NewCommodityCache(re) 或其他名称，请替换
	return NewCacheAdapter(re)
}

func TestCacheAdapter_SetAndGetCartCache(t *testing.T) {
	// 如果你的项目通过环境变量来控制测试，按需判断
	if !utils.EnvironmentEnable() {
		return
	}

	// 初始化
	c := initTestRedis(t)
	ctx := context.Background()

	// 构造一个测试用的 key & cart 数据
	cartKey := "test:user:cart:12345"
	cartData := `{"items": [{"skuId": 101, "count": 2}]}`

	Convey("测试 SetCartCache & GetCartCache", t, func() {
		Convey("当向 Redis 写入购物车数据", func() {
			err := c.SetCartCache(ctx, cartKey, cartData)
			So(err, ShouldBeNil)

			Convey("应该能够正常读取相同的数据", func() {
				ret, err := c.GetCartCache(ctx, cartKey)
				So(err, ShouldBeNil)
				So(ret, ShouldEqual, cartData)
			})
		})
	})
}
