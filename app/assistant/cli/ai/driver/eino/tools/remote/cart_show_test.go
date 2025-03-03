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
	"fmt"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/bytedance/sonic"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/west2-online/DomTok/app/gateway/model/api/cart"

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

func TestToolCartShow_InvokableRun(t *testing.T) {
	f := CartShow(nil)
	type MockServerCaller struct {
		adapter.ServerCaller
	}
	fakeServerCaller := &MockServerCaller{}
	args := ToolCartShowArgs{PageNum: 4}
	argsBytes, _ := sonic.Marshal(args)
	PatchConvey("Test ToolCartShow", t, func() {
		PatchConvey("success", func() {
			MockValue(&f.getServerCaller).To(func(_ string) adapter.ServerCaller { return fakeServerCaller })
			Mock((*MockServerCaller).CartShow).To(func(_ context.Context, params *cart.ShowCartGoodsListRequest) ([]byte, error) {
				return []byte(fmt.Sprintf("page_num: %d", params.PageNum)), nil
			}).Build()

			resp, err := f.InvokableRun(context.Background(), string(argsBytes))
			So(err, ShouldBeNil)
			So(resp, ShouldEqual, "page_num: 4")
		})

		PatchConvey("if server caller is nil", func() {
			MockValue(&f.getServerCaller).To(func(_ string) adapter.ServerCaller { return nil })
			Mock((*MockServerCaller).CartShow).Return([]byte("pong"), nil).Build()

			_, err := f.InvokableRun(context.Background(), string(argsBytes))

			So(err, ShouldNotBeNil)
		})

		PatchConvey("if server caller returns error", func() {
			MockValue(&f.getServerCaller).To(func(_ string) adapter.ServerCaller { return fakeServerCaller })
			Mock((*MockServerCaller).CartShow).Return(nil, errors.New("dial error")).Build()

			resp, err := f.InvokableRun(context.Background(), string(argsBytes))

			So(err, ShouldBeNil)
			So(resp, ShouldEqual, "dial error")
		})
	})
}
