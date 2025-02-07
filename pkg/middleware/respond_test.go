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

	"github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

type respondResp struct {
	base *model.BaseResp
}

func (resp *respondResp) IsSetBase() bool {
	return resp.base != nil
}

func (resp *respondResp) GetBase() *model.BaseResp {
	return resp.base
}

func (resp *respondResp) SetBase(b *model.BaseResp) {
	resp.base = b
}

type respondResult struct {
	resp *respondResp
}

func (rel *respondResult) IsSetSuccess() bool {
	return rel.resp != nil
}

func (rel *respondResult) GetResult() interface{} {
	return rel.resp
}

func respondNext(ctx context.Context, req, resp interface{}) (err error) {
	return nil
}

func TestRespond(t *testing.T) {
	mid := Respond()
	PatchConvey("Test the middleware respond", t, func() {
		PatchConvey("Test when pass result as nil", func() {
			err := mid(respondNext)(context.Background(), nil, nil)
			So(err, ShouldBeNil)
		})

		PatchConvey("Test when pass result.resp as nil", func() {
			result := &respondResult{resp: nil}
			err := mid(respondNext)(context.Background(), nil, result)
			So(err, ShouldBeNil)
			So(result.resp, ShouldBeNil)
		})

		PatchConvey("Test when pass normal result and nil base", func() {
			result := &respondResult{resp: &respondResp{}}
			err := mid(respondNext)(context.Background(), nil, result)
			So(err, ShouldBeNil)

			res, ok := result.GetResult().(baser)
			So(ok, ShouldBeTrue)
			So(res, ShouldEqual, result.resp)

			base := res.GetBase()
			So(base, ShouldNotBeNil)
			So(base.GetCode(), ShouldEqual, errno.SuccessCode)
			So(base.GetMsg(), ShouldEqual, errno.SuccessMsg)
		})

		PatchConvey("Test when pass normal result and base", func() {
			code, msg := int64(200), "ok"
			result := &respondResult{resp: &respondResp{base: &model.BaseResp{Code: code, Msg: msg}}}
			err := mid(respondNext)(context.Background(), nil, result)
			So(err, ShouldBeNil)

			res, ok := result.GetResult().(baser)
			So(ok, ShouldBeTrue)
			So(res, ShouldEqual, result.resp)

			base := res.GetBase()
			So(base, ShouldNotBeNil)
			So(base.GetCode(), ShouldEqual, code)
			So(base.GetMsg(), ShouldEqual, msg)
		})

		PatchConvey("Test when next return normal error", func() {
			code, msg := int64(200), "ok"
			result := &respondResult{resp: &respondResp{base: &model.BaseResp{Code: code, Msg: msg}}}
			e := fmt.Errorf("test")
			Mock(respondNext).Return(e).Build()

			err := mid(respondNext)(context.Background(), nil, result)
			So(err, ShouldEqual, e)
		})

		PatchConvey("Test when next return errno error, and resp.Base != nil", func() {
			code, msg := int64(200), "ok"
			result := &respondResult{resp: &respondResp{base: &model.BaseResp{Code: code, Msg: msg}}}
			e := errno.NewErrNo(int64(1000), errno.SuccessMsg)
			Mock(respondNext).Return(e).Build()

			err := mid(respondNext)(context.Background(), nil, result)
			So(err, ShouldBeNil)

			So(result.resp.base.GetCode(), ShouldEqual, code)
			So(result.resp.base.GetMsg(), ShouldEqual, msg)
		})

		PatchConvey("Test when next return errno error, and resp.Base == nil", func() {
			code, msg := int64(200), "ok"
			result := &respondResult{resp: &respondResp{}}
			e := errno.NewErrNo(int64(1000), "error msg")
			Mock(respondNext).Return(e).Build()

			err := mid(respondNext)(context.Background(), nil, result)
			So(err, ShouldBeNil)

			So(result.resp.base.GetCode(), ShouldNotEqual, code)
			So(result.resp.base.GetMsg(), ShouldNotEqual, msg)

			So(result.resp.base.GetCode(), ShouldEqual, e.ErrorCode)
			So(result.resp.base.GetMsg(), ShouldEqual, e.ErrorMsg)
		})
	})
}
