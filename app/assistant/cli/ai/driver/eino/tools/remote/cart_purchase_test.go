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
	"github.com/bytedance/sonic"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
	"github.com/west2-online/DomTok/app/gateway/model/api/order"
)

func TestToolCartPurchase_InvokableRun(t *testing.T) {
	f := CartPurchase(nil)
	type MockServerCaller struct {
		adapter.ServerCaller
	}
	fakeServerCaller := &MockServerCaller{}
	args := ToolCartPurchaseArgs{
		BaseOrderGoods: []_CartPurchaseBaseOrderGoods{{
			MerchantID: 1,
			GoodsID:    2,
			SkuID:      3,
		}},
	}
	argsBytes, _ := sonic.Marshal(args)
	PatchConvey("Test OrderCreate", t, func() {
		PatchConvey("success", func() {
			mp := map[string]interface{}{}
			MockValue(&f.getServerCaller).To(func(_ string) adapter.ServerCaller { return fakeServerCaller })
			Mock((*MockServerCaller).CartPurchase).To(func(_ context.Context, params *order.CreateOrderReq) ([]byte, error) {
				mp["base_order_goods"] = params.BaseOrderGoods
				return nil, nil
			}).Build()

			_, err := f.InvokableRun(context.Background(), string(argsBytes))
			So(err, ShouldBeNil)
			So(mp["base_order_goods"], ShouldResemble, ConvertArgsOrderGoodsToRequestGoods(args.BaseOrderGoods...))
		})

		PatchConvey("if server caller is nil", func() {
			MockValue(&f.getServerCaller).To(func(_ string) adapter.ServerCaller { return nil })
			Mock((*MockServerCaller).CartPurchase).Return([]byte("pong"), nil).Build()

			_, err := f.InvokableRun(context.Background(), string(argsBytes))

			So(err, ShouldNotBeNil)
		})

		PatchConvey("if server caller returns error", func() {
			MockValue(&f.getServerCaller).To(func(_ string) adapter.ServerCaller { return fakeServerCaller })
			Mock((*MockServerCaller).CartPurchase).Return(nil, errors.New("dial error")).Build()

			resp, err := f.InvokableRun(context.Background(), string(argsBytes))

			So(err, ShouldBeNil)
			So(resp, ShouldEqual, "dial error")
		})
	})
}
