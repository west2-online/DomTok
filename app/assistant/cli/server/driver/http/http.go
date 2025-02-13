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
	"time"

	"github.com/cloudwego/hertz/pkg/app/client"

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

type Client struct {
	adapter.ServerCaller

	cli     *client.Client
	BaseUrl string
}

type ClientOption struct {
	BaseUrl string
}

func NewClient(opt *ClientOption) *Client {
	cli, _ := client.NewClient(
		client.WithDialTimeout(time.Second),
		client.WithClientReadTimeout(time.Second),
		client.WithWriteTimeout(time.Second),
	)
	return &Client{cli: cli, BaseUrl: opt.BaseUrl}
}
