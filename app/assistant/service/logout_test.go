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

	"github.com/west2-online/DomTok/app/assistant/cli/ai/adapter"
	"github.com/west2-online/DomTok/app/assistant/model"
)

func TestCore_Logout(t *testing.T) {
	PatchConvey("Test Core.Logout", t, func() {
		// prepare
		type EmptyCli struct {
			adapter.AIClient
		}
		c := Core{}
		c.ai = &EmptyCli{}
		ctx := context.Background()
		ctx = context.WithValue(ctx, CtxKeyID, "id")
		forgot := false
		Mock((*EmptyCli).ForgetDialog).To(func(dialog model.IDialog) {
			forgot = true
		}).Build()

		PatchConvey("success", func() {
			err := c.Logout(ctx)
			So(err, ShouldBeNil)
			So(forgot, ShouldBeTrue)
		})

		PatchConvey("no id", func() {
			ctx = context.Background()
			err := c.Logout(ctx)
			So(err, ShouldBeNil)
			So(forgot, ShouldBeFalse)
		})

		PatchConvey("id type error", func() {
			ctx = context.WithValue(ctx, CtxKeyID, 1)
			err := c.Logout(ctx)
			So(err, ShouldBeNil)
			So(forgot, ShouldBeFalse)
		})
	})
}
