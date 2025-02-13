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

package volcengine

import (
	"sync"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	arkmodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/adapter"
	server "github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
	"github.com/west2-online/DomTok/app/assistant/model"
)

type Client struct {
	adapter.AIClient

	cli     *arkruntime.Client
	caller  server.ServerCaller
	baseReq *arkmodel.CreateChatCompletionRequest

	recorder sync.Map
}

type ClientOption struct {
	ApiKey  string
	BaseUrl string
	Region  string
	Model   string
}

func NewClient(opt *ClientOption) *Client {
	cli := arkruntime.NewClientWithApiKey(
		opt.ApiKey,
		arkruntime.WithBaseUrl(opt.BaseUrl),
		arkruntime.WithRegion(opt.Region),
	)
	baseReq := &arkmodel.CreateChatCompletionRequest{
		Model:    opt.Model,
		Messages: *GetPrevMessages(),
		Tools:    *GetTools(),
	}
	return &Client{cli: cli, baseReq: baseReq}
}

func (c *Client) SetServerCaller(caller server.ServerCaller) {
	c.caller = caller
}

func (c *Client) Call(dialog model.IDialog) (err error) {
	defer dialog.Close()

	id := dialog.Unique()
	h := ""
	v, ok := c.recorder.Load(id)
	if ok {
		h, _ = v.(string)
	}

	input := dialog.Message()

	h += input
	dialog.Send(h)

	c.recorder.Store(id, h)

	return nil
}

func (c *Client) ForgetDialog(dialog model.IDialog) {
	c.recorder.Delete(dialog.Unique())
}
