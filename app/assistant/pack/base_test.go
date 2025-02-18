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

package pack

import (
	"encoding/json"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/assistant/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

func Test_ResponseFactory_ConnectSuccess(t *testing.T) {
	PatchConvey("Test ResponseFactory.ConnectSuccess", t, func() {
		resp := ResponseFactory.ConnectSuccess("extra")
		v := model.Response{}
		_ = json.Unmarshal(resp, &v)
		So(v.Meta[MetaType], ShouldEqual, MetaTypePing)
		So(v.Meta[MetaExtra], ShouldEqual, "extra")
	})
}

func Test_ResponseFactory_Command(t *testing.T) {
	PatchConvey("Test ResponseFactory.Command", t, func() {
		resp := ResponseFactory.Command("params")
		v := model.Response{}
		_ = json.Unmarshal(resp, &v)
		So(v.Meta[MetaType], ShouldEqual, MetaTypeCommand)
		So(v.Data, ShouldEqual, "params")
	})
}

func Test_ResponseFactory_Error(t *testing.T) {
	PatchConvey("Test ResponseFactory.Error", t, func() {
		err := errno.ConvertErr(nil)
		e := model.NewErrorData(err.ErrorCode, err.ErrorMsg)
		resp := ResponseFactory.Error(err)
		v := model.Response{}
		_ = json.Unmarshal(resp, &v)
		d := map[string]interface{}{}
		data, _ := json.Marshal(e)
		_ = json.Unmarshal(data, &d)
		So(v.Meta[MetaType], ShouldEqual, MetaTypeError)
		So(v.Data, ShouldResemble, d)
	})
}

func Test_ResponseFactory_Message(t *testing.T) {
	PatchConvey("Test ResponseFactory.Message", t, func() {
		resp := ResponseFactory.Message("params")
		v := model.Response{}
		_ = json.Unmarshal(resp, &v)
		So(v.Meta[MetaType], ShouldEqual, MetaTypeMessage)
		So(v.Data, ShouldEqual, "params")
	})
}
