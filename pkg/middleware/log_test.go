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

package middleware

import (
	"context"
	"fmt"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/pkg/errno"
)

func errorLogNext(ctx context.Context, req, resp interface{}) (err error) {
	return nil
}

func TestErrorLog(t *testing.T) {
	mid := ErrorLog()
	PatchConvey("Test the middleware ErrorLog", t, func() {
		PatchConvey("Test when next return nil", func() {
			Mock(errorLogNext).Return(nil).Build()
			err := mid(errorLogNext)(context.Background(), nil, nil)
			So(err, ShouldEqual, nil)
		})

		PatchConvey("Test when next return normal error", func() {
			e := fmt.Errorf("test")
			Mock(errorLogNext).Return(e).Build()
			err := mid(errorLogNext)(context.Background(), nil, nil)
			So(err, ShouldEqual, e)
		})

		PatchConvey("Test when next return errno", func() {
			e := errno.NewErrNo(errno.SuccessCode, "ok")
			Mock(errorLogNext).Return(e).Build()
			err := mid(errorLogNext)(context.Background(), nil, nil)
			So(err, ShouldEqual, e)
		})

		PatchConvey("Test when next return use fmt.Errorf pack errno", func() {
			en := errno.NewErrNo(errno.SuccessCode, "ok")
			e := fmt.Errorf("test`s %w", en)
			Mock(errorLogNext).Return(e).Build()
			err := mid(errorLogNext)(context.Background(), nil, nil)
			So(err, ShouldEqual, e)
		})
	})
}
