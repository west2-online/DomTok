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

package base

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestBuildBaseResp(t *testing.T) {
	Convey("TestBuildBaseResp", t, func() {
		nilError := BuildBaseResp(nil)
		So(nilError.Code, ShouldEqual, int64(errno.SuccessCode))
		So(nilError.Code, ShouldEqual, int64(errno.SuccessCode))
		So(nilError.Msg, ShouldEqual, errno.Success.ErrorMsg)

		normalError := BuildBaseResp(fmt.Errorf("ok"))
		So(normalError.Code, ShouldEqual, int64(errno.InternalServiceErrorCode))
		So(normalError.Msg, ShouldEqual, "ok")

		errnoError := BuildBaseResp(errno.NewErrNo(200, "ok"))
		So(errnoError.Code, ShouldEqual, int64(200))
		So(errnoError.Msg, ShouldEqual, "ok")
	})
}

func TestBuildSuccessResp(t *testing.T) {
	Convey("TestBuildSuccessResp", t, func() {
		r := BuildSuccessResp()
		So(r.Code, ShouldEqual, int64(errno.SuccessCode))
		So(r.Msg, ShouldEqual, errno.Success.ErrorMsg)
	})
}

func TestLogError(t *testing.T) {
	LogError(nil)
	LogError(fmt.Errorf("ok"))
	LogError(errno.Success)
	// LogError(errno.NewErrNoWithStack(200, "ok")) // have tested
}

func TestBuildRespAndLog(t *testing.T) {
	Convey("Test BuildRespAndLog", t, func() {
		nilError := BuildBaseResp(nil)
		So(nilError.Code, ShouldEqual, int64(errno.SuccessCode))
		So(nilError.Msg, ShouldEqual, errno.Success.ErrorMsg)

		normalError := BuildBaseResp(fmt.Errorf("ok"))
		So(normalError.Code, ShouldEqual, int64(errno.InternalServiceErrorCode))
		So("ok", ShouldEqual, normalError.Msg)

		errnoError := BuildBaseResp(errno.NewErrNo(200, "ok"))
		So(errnoError.Code, ShouldEqual, int64(200))
		So(errnoError.Msg, ShouldEqual, "ok")
	})
}

func TestBuildTypeList(t *testing.T) {
	Convey("Test BuildTypeList", t, func() {
		f := func(i int64) string {
			return strconv.FormatInt(i, 10)
		}
		ints := []int64{1, 2, 3, 4}
		strs := BuildTypeList(ints, f)
		So(len(strs), ShouldEqual, len(ints))
		So([]string{"1", "2", "3", "4"}, ShouldResemble, strs)
	})
}
