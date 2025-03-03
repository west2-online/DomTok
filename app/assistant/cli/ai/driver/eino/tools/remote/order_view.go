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
	strategy "github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/tools"
	"github.com/west2-online/DomTok/app/gateway/model/api/order"
	"github.com/west2-online/DomTok/pkg/errno"
)

const (
	ToolOrderViewName = "view_order"
	ToolOrderViewDesc = "当用户想查看订单详情时，可以使用此工具"
)

type ToolOrderViewArgs struct {
	OrderID int64 `json:"order_id" desc:"填入订单ID" required:"true"`
}

var ToolOrderViewRequestBody = schema.NewParamsOneOfByParams(*tools.Reflect(ToolOrderViewArgs{}))

type ToolOrderView struct {
	tool.InvokableTool

	getServerCaller strategy.GetServerCaller
}

func OrderView(strategy strategy.GetServerCaller) *ToolOrderView {
	return &ToolOrderView{
		getServerCaller: strategy,
	}
}

func (t *ToolOrderView) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	s := t.getServerCaller(ToolOrderViewName)
	if s == nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "server is nil")
	}
	params := &ToolOrderViewArgs{}
	if err := sonic.Unmarshal([]byte(argumentsInJSON), params); err != nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "unmarshal arguments failed")
	}
	resp, err := s.OrderView(ctx, &order.ViewOrderReq{OrderID: params.OrderID})
	if err != nil {
		return err.Error(), nil //nolint: nilerr
	}

	return string(resp), nil
}

func (t *ToolOrderView) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ToolOrderViewName,
		Desc:        ToolOrderViewDesc,
		ParamsOneOf: ToolOrderViewRequestBody,
	}, nil
}
