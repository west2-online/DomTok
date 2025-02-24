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

package service

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCore_Login(t *testing.T) {
	PatchConvey("Test Core.Login", t, func() {
		// prepare
		c := Core{}
		ctx := context.Background()

		PatchConvey("success", func() {
			ctx = context.WithValue(ctx, CtxKeyID, "id")
			ctx = context.WithValue(ctx, CtxKeyAccessToken, "token")
			err := c.Login(ctx)
			So(err, ShouldBeNil)
		})

		PatchConvey("no id", func() {
			ctx = context.WithValue(ctx, CtxKeyAccessToken, "token")
			err := c.Login(ctx)
			So(err, ShouldNotBeNil)
		})

		PatchConvey("no access token", func() {
			ctx = context.WithValue(ctx, CtxKeyID, "id")
			err := c.Login(ctx)
			So(err, ShouldNotBeNil)
		})

		PatchConvey("id type error", func() {
			ctx = context.WithValue(ctx, CtxKeyID, 1)
			ctx = context.WithValue(ctx, CtxKeyAccessToken, "token")
			err := c.Login(ctx)
			So(err, ShouldNotBeNil)
		})

		PatchConvey("access token type error", func() {
			ctx = context.WithValue(ctx, CtxKeyID, "id")
			ctx = context.WithValue(ctx, CtxKeyAccessToken, 1)
			err := c.Login(ctx)
			So(err, ShouldNotBeNil)
		})
	})
}
