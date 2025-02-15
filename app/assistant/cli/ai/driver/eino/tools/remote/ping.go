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

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
	"github.com/west2-online/DomTok/pkg/logger"
)

// Tips: This function should not be used in the future

type ToolPing struct {
	tool.InvokableTool

	server adapter.ServerCaller
}

const (
	ToolPingName = "ping"
	ToolPingDesc = "当用户想知道服务器是否在线时，可以使用此工具"
)

func Ping(server adapter.ServerCaller) *ToolPing {
	return &ToolPing{
		server: server,
	}
}

func (t *ToolPing) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	if t.server == nil {
		panic("server is nil")
	}
	resp, err := t.server.Ping(ctx)
	if err != nil {
		return "", err
	}

	// TODO: remove this line
	logger.Infof(`{
Stage: "remote.Ping",
args: %v,
resp: %v,
opts: %v,
err: %v,
}`, argumentsInJSON, string(resp), opts, err)
	return string(resp), nil
}

func (t *ToolPing) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: ToolPingName,
		Desc: ToolPingDesc,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"argument": {
				Type:     schema.String,
				Desc:     "填入角色设定名",
				Required: true,
			},
		}),
	}, nil
}
