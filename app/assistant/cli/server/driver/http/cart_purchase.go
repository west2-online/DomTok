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

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/DomTok/app/gateway/model/api/cart"
)

const (
	_CartPurchasePath   = "/api/v1/cart/purchase"
	_CartPurchaseMethod = consts.MethodPost
)

func (c *Client) CartPurchase(ctx context.Context, params *cart.PurChaseCartGoodsRequest) ([]byte, error) {
	body, err := sonic.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, resp := protocol.AcquireRequest(), protocol.AcquireResponse()
	req.SetRequestURI(c.buildUrl(_CartPurchasePath))
	req.SetMethod(_CartPurchaseMethod)
	req.SetBody(body)
	req.Header.Set("Content-Type", "application/json") //nolint

	err = c.do(ctx, req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
