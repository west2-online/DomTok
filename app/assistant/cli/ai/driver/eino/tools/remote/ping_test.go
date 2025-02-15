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

package remote

import (
	"context"
	"errors"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

func TestToolPing_InvokableRun(t *testing.T) {
	f := Ping(nil)
	type MockServerCaller struct {
		adapter.ServerCaller
	}
	fakeServerCaller := &MockServerCaller{}
	PatchConvey("TestToolPing_InvokableRun", t, func() {
		PatchConvey("success", func() {
			MockValue(&f.server).To(fakeServerCaller)
			Mock((*MockServerCaller).Ping).Return([]byte("pong"), nil).Build()

			resp, err := f.InvokableRun(context.Background(), "")
			So(err, ShouldBeNil)
			So(resp, ShouldEqual, "pong")
		})

		PatchConvey("if server caller is nil", func() {
			MockValue(&f.server).To(nil)
			Mock((*MockServerCaller).Ping).Return([]byte("pong"), nil).Build()

			_, err := f.InvokableRun(context.Background(), "")

			So(err, ShouldNotBeNil)
		})

		PatchConvey("if server caller returns error", func() {
			MockValue(&f.server).To(fakeServerCaller)
			Mock((*MockServerCaller).Ping).Return(nil, errors.New("")).Build()

			_, err := f.InvokableRun(context.Background(), "")

			So(err, ShouldNotBeNil)
		})
	})
}
