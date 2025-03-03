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
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/redis/go-redis/v9"
	"github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"

	"github.com/west2-online/DomTok/pkg/utils"
)

func TestPaymentRedis_SetPaymentToken(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	ctx := context.Background()
	mockey.PatchConvey("SetPaymentToken", t, func() {
		err := _cli.SetPaymentToken(ctx, "key", "value", -1)
		convey.So(err, convey.ShouldBeNil)
		v, err := _cli.client.Get(ctx, "key").Result()
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "value")
	})
}

func TestPaymentRedis_IncrRedisKey(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	ctx := context.Background()
	mockey.PatchConvey("IncrRedisKey", t, func() {
		count, err := _cli.IncrRedisKey(ctx, "key", -1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(count, convey.ShouldEqual, 1)
	})
}

func TestPaymentRedis_CheckRedisDayKey(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	ctx := context.Background()
	mockey.PatchConvey("CheckRedisDayKey", t, func() {
		exists, err := _cli.CheckRedisDayKey(ctx, "key")
		convey.So(err, convey.ShouldBeNil)
		convey.So(exists, convey.ShouldBeFalse)

		_cli.client.Set(ctx, "key", "value", -1)
		exists, err = _cli.CheckRedisDayKey(ctx, "key")
		convey.So(err, convey.ShouldBeNil)
		convey.So(exists, convey.ShouldBeTrue)
	})
}

func TestPaymentRedis_SetRedisDayKey(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	ctx := context.Background()
	mockey.PatchConvey("SetRedisDayKey", t, func() {
		err := _cli.SetRedisDayKey(ctx, "key", "value", -1)
		convey.So(err, convey.ShouldBeNil)
		v, err := _cli.client.Get(ctx, "key").Result()
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "value")
	})
}

func TestPaymentRedis_SetRefundToken(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	ctx := context.Background()
	mockey.PatchConvey("SetRefundToken", t, func() {
		err := _cli.SetRefundToken(ctx, "key", "value", -1)
		convey.So(err, convey.ShouldBeNil)
		v, err := _cli.client.Get(ctx, "key").Result()
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "value")
		ttl, err := _cli.client.TTL(ctx, "key").Result()
		convey.So(err, convey.ShouldBeNil)
		convey.So(ttl, convey.ShouldEqual, -1)
	})
}

func TestPaymentRedis_CheckAndDelPaymentToken(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	ctx := context.Background()
	mockey.PatchConvey("CheckAndDelPaymentToken", t, func() {
		_cli.client.Set(ctx, "key", "value", -1)
		exist, err := _cli.CheckAndDelPaymentToken(ctx, "key", "value")
		convey.So(err, convey.ShouldBeNil)
		convey.So(exist, convey.ShouldBeTrue)
		err = _cli.client.Get(ctx, "key").Err()
		convey.So(err, convey.ShouldEqual, redis.Nil)

		exist, err = _cli.CheckAndDelPaymentToken(ctx, "key", "value")
		convey.So(err, convey.ShouldBeNil)
		convey.So(exist, convey.ShouldBeFalse)
	})
}

func TestPaymentRedis_GetTTLAndDelPaymentToken(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	ctx := context.Background()
	mockey.PatchConvey("GetTTLAndDelPaymentToken", t, func() {
		_cli.client.Set(ctx, "key", "value", -1)
		exist, ttl, err := _cli.GetTTLAndDelPaymentToken(ctx, "key", "value")
		convey.So(err, convey.ShouldBeNil)
		convey.So(exist, convey.ShouldBeTrue)
		convey.So(ttl, convey.ShouldEqual, time.Duration(-1)*time.Second)
		err = _cli.client.Get(ctx, "key").Err()
		convey.So(err, convey.ShouldEqual, redis.Nil)

		exist, _, err = _cli.GetTTLAndDelPaymentToken(ctx, "key", "value")
		convey.So(err, convey.ShouldBeNil)
		convey.So(exist, convey.ShouldBeFalse)
	})
}
