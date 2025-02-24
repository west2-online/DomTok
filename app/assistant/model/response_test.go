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

package model

import (
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestResponse_SetMeta(t *testing.T) {
	PatchConvey("Test Response.SetMeta", t, func() {
		r := NewResponse()
		r.SetMeta("key", "value")
		So(r.Meta["key"], ShouldEqual, "value")
	})
}

func TestResponse_SetData(t *testing.T) {
	PatchConvey("Test Response.SetData", t, func() {
		r := NewResponse()
		r.SetData("data")
		So(r.Data, ShouldEqual, "data")
	})
}

func TestResponse_GetMeta(t *testing.T) {
	PatchConvey("Test Response.GetMeta", t, func() {
		r := NewResponse()
		r.SetMeta("key", "value")
		So(r.GetMeta("key"), ShouldEqual, "value")
		So(r.GetMeta("key2"), ShouldBeNil)
	})
}

func TestResponse_GetData(t *testing.T) {
	PatchConvey("Test Response.GetData", t, func() {
		r := NewResponse()
		r.SetData("data")
		So(r.GetData(), ShouldEqual, "data")
	})
}

func TestResponse_Marshal(t *testing.T) {
	PatchConvey("Test Response.Marshal", t, func() {
		r := NewResponse()
		r.SetMeta("key", "value")
		r.SetData("data")
		b, err := r.Marshal()
		So(err, ShouldBeNil)
		So(string(b), ShouldEqual, `{"meta":{"key":"value"},"data":"data"}`)
	})
}

func TestResponse_MustMarshal(t *testing.T) {
	PatchConvey("Test Response.MustMarshal", t, func() {
		r := NewResponse()
		r.SetMeta("key", "value")
		r.SetData("data")
		b := r.MustMarshal()
		So(string(b), ShouldEqual, `{"meta":{"key":"value"},"data":"data"}`)
	})

	PatchConvey("Test Response.MustMarshal with nil meta", t, func() {
		r := NewResponse()
		r.Meta = nil
		r.SetData("data")
		b := r.MustMarshal()
		So(string(b), ShouldEqual, `{"meta":{},"data":"data"}`)
	})

	PatchConvey("Test Response.MustMarshal with nil data", t, func() {
		r := NewResponse()
		r.SetMeta("key", "value")
		b := r.MustMarshal()
		So(string(b), ShouldEqual, `{"meta":{"key":"value"},"data":{}}`)
	})
}

func TestNewConnectSuccess(t *testing.T) {
	PatchConvey("Test NewConnectSuccess", t, func() {
		r := NewConnectSuccess("id", "tz")
		So(r, ShouldNotBeNil)
		So(r.DialogID, ShouldEqual, "id")
		So(r.TZ, ShouldEqual, "tz")
	})
}

func TestNewDeltaContent(t *testing.T) {
	PatchConvey("Test NewDeltaContent", t, func() {
		r := NewDeltaContent("delta", 1, 2)
		So(r, ShouldNotBeNil)
		So(r.Delta, ShouldEqual, "delta")
		So(r.Index, ShouldEqual, 1)
		So(r.Turn, ShouldEqual, 2)
	})
}

func TestNewErrorData(t *testing.T) {
	PatchConvey("Test NewErrorData", t, func() {
		r := NewErrorData(1, "error")
		So(r, ShouldNotBeNil)
		So(r.Code, ShouldEqual, 1)
		So(r.Error, ShouldEqual, "error")
	})
}

func TestNewDialogOp(t *testing.T) {
	PatchConvey("Test NewDialogOp", t, func() {
		r := NewDialogOp("id", 1)
		So(r, ShouldNotBeNil)
		So(r.Content, ShouldEqual, "id")
		So(r.Turn, ShouldEqual, 1)
	})
}
