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

	strategy "github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/tools"
	"github.com/west2-online/DomTok/pkg/errno"
)

// Tips: This function should not be used in the future

const (
	ToolPingName = "ping"
	ToolPingDesc = "当用户想知道服务器是否在线时，可以使用此工具"
)

type ToolPingArgs struct {
	Argument string `json:"argument" desc:"填入角色设定名" required:"true"`
}

var ToolPingRequestBody = schema.NewParamsOneOfByParams(*tools.Reflect(ToolPingArgs{}))

type ToolPing struct {
	tool.InvokableTool

	getServerCaller strategy.GetServerCaller
}

func Ping(strategy strategy.GetServerCaller) *ToolPing {
	return &ToolPing{
		getServerCaller: strategy,
	}
}

func (t *ToolPing) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	s := t.getServerCaller(ToolPingName)
	if s == nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "server is nil")
	}
	resp, err := s.Ping(ctx)
	if err != nil {
		return err.Error(), nil //nolint: nilerr
	}

	return string(resp), nil
}

func (t *ToolPing) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ToolPingName,
		Desc:        ToolPingDesc,
		ParamsOneOf: ToolPingRequestBody,
	}, nil
}
