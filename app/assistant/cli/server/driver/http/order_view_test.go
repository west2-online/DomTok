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

package http

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/protocol"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/gateway/model/api/order"
)

func TestClient_OrderView(t *testing.T) {
	c := &Client{}
	PatchConvey("HTTP Client.OrderView", t, func() {
		PatchConvey("success resp", func() {
			mp := make(map[string]string)
			Mock((*Client).do).To(func(ctx context.Context, req *protocol.Request, resp *protocol.Response) error {
				args := req.URI().QueryArgs()
				args.VisitAll(func(key, value []byte) {
					if string(key) == "orderID" {
						mp[string(key)] = string(value)
					}
				})
				resp.SetBody([]byte("orderID=" + mp["orderID"]))
				return nil
			}).Build()

			resp, err := c.OrderView(context.Background(), &order.ViewOrderReq{
				OrderID: 1,
			})

			So(err, ShouldBeNil)
			So(string(resp), ShouldEqual, "orderID=1")
		})

		PatchConvey("error resp", func() {
			Mock((*Client).do).To(func(ctx context.Context, req *protocol.Request, resp *protocol.Response) error {
				return errors.ErrTimeout
			}).Build()

			_, err := c.OrderView(context.Background(), &order.ViewOrderReq{})
			So(err, ShouldEqual, errors.ErrTimeout)
		})
	})
}
