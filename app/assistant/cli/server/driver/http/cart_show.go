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
	"net/url"
	"strconv"

	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/DomTok/app/gateway/model/api/cart"
)

const (
	_CartShowPath   = "/api/v1/cart/show"
	_CartShowMethod = consts.MethodGet
)

func (c *Client) CartShow(ctx context.Context, params *cart.ShowCartGoodsListRequest) ([]byte, error) {
	uri, err := url.Parse(c.buildUrl(_CartShowPath))
	if err != nil {
		return nil, err
	}
	query := uri.Query()
	query.Set("page_num", strconv.FormatInt(params.PageNum, 10))
	uri.RawQuery = query.Encode()

	req, resp := protocol.AcquireRequest(), protocol.AcquireResponse()
	req.SetRequestURI(uri.String())
	req.SetMethod(_CartShowMethod)

	err = c.do(ctx, req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
