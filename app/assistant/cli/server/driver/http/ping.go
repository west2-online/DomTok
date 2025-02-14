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

	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/DomTok/pkg/logger"
)

// Tips: This function should not be used in the future

const (
	_PingPath   = "/ping"
	_PingMethod = consts.MethodGet
)

func (c *Client) Ping(ctx context.Context) ([]byte, error) {
	req, resp := protocol.AcquireRequest(), protocol.AcquireResponse()
	req.SetRequestURI(c.buildUrl(_PingPath))
	req.SetMethod(_PingMethod)

	err := c.do(ctx, req, resp)
	// TODO: remove this line
	logger.Infof(`{
Stage: "http.Ping",
req: %v,
resp: %v,
err: %v,
}`, req, string(resp.Body()), err)

	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
