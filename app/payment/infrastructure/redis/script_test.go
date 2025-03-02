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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/pkg/utils"
)

func TestPaymentRedis_LoadScript(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	err := _cli.client.ScriptFlush(context.Background()).Err()
	if err != nil {
		t.Fatal(err)
	}

	mockey.PatchConvey("LoadScript", t, func() {
		err := _cli.loadScript()
		convey.So(err, convey.ShouldBeNil)
		for _, value := range scripts {
			exists, err := _cli.client.ScriptExists(context.Background(), value.Hash).Result()
			convey.So(err, convey.ShouldBeNil)
			convey.So(exists[0], convey.ShouldBeTrue)
		}
	})

	err = _cli.client.ScriptFlush(context.Background()).Err()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPaymentRedis_ExecScript(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	_EnvSetup()
	scripts = map[ScriptKey]*_Script{
		CheckAndDelScript: {
			Cmd: `return {KEYS[1], ARGV[1]}`,
		},
	}
	err := _cli.client.ScriptFlush(context.Background()).Err()
	if err != nil {
		t.Fatal(err)
	}
	err = _cli.loadScript()
	if err != nil {
		t.Fatal(err)
	}

	mockey.PatchConvey("ExecScript", t, func() {
		res, err := _cli.execScript(context.Background(), CheckAndDelScript, []string{"key"}, "value")
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldResemble, []interface{}{"key", "value"})
	})

	err = _cli.client.ScriptFlush(context.Background()).Err()
	if err != nil {
		t.Fatal(err)
	}
}
