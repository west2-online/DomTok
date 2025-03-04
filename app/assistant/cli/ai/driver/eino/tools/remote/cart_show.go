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

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/west2-online/DomTok/app/gateway/model/api/cart"

	strategy "github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/tools"
	"github.com/west2-online/DomTok/pkg/errno"
)

const (
	ToolCartShowName = "show_cart"
	ToolCartShowDesc = "当用户想查看购物车或从购物车下单时，可以使用此工具"
)

type ToolCartShowArgs struct {
	PageNum int64 `json:"page_num" desc:"填入页码" required:"true"`
}

var ToolCartShowRequestBody = schema.NewParamsOneOfByParams(*tools.Reflect(ToolCartShowArgs{}))

type ToolCartShow struct {
	tool.InvokableTool

	getServerCaller strategy.GetServerCaller
}

func CartShow(strategy strategy.GetServerCaller) *ToolCartShow {
	return &ToolCartShow{
		getServerCaller: strategy,
	}
}

func (t *ToolCartShow) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	s := t.getServerCaller(ToolPingName)
	if s == nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "server is nil")
	}
	params := &ToolCartShowArgs{}
	if err := sonic.Unmarshal([]byte(argumentsInJSON), params); err != nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "unmarshal arguments failed")
	}
	resp, err := s.CartShow(ctx, &cart.ShowCartGoodsListRequest{PageNum: params.PageNum})
	if err != nil {
		return err.Error(), nil //nolint: nilerr
	}

	return string(resp), nil
}

func (t *ToolCartShow) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ToolCartShowName,
		Desc:        ToolCartShowDesc,
		ParamsOneOf: ToolCartShowRequestBody,
	}, nil
}
