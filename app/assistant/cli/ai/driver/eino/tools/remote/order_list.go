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
	ToolOrderListName = "list_order"
	ToolOrderListDesc = "当用户想查看订单列表时，可以使用此工具"
)

type ToolOrderListArgs struct {
	Page int64 `json:"page" desc:"填入页码" required:"true"`
	Size int64 `json:"size" desc:"填入每页数量" required:"true"`
}

var ToolOrderListRequestBody = schema.NewParamsOneOfByParams(*tools.Reflect(ToolOrderListArgs{}))

type ToolOrderList struct {
	tool.InvokableTool

	getServerCaller strategy.GetServerCaller
}

func OrderList(strategy strategy.GetServerCaller) *ToolOrderList {
	return &ToolOrderList{
		getServerCaller: strategy,
	}
}

func (t *ToolOrderList) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	s := t.getServerCaller(ToolOrderListName)
	if s == nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "server is nil")
	}
	params := &ToolOrderListArgs{}
	if err := sonic.Unmarshal([]byte(argumentsInJSON), params); err != nil {
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "unmarshal arguments failed")
	}
	resp, err := s.OrderList(ctx, &order.ViewOrderListReq{
		Page: int32(params.Page),
		Size: int32(params.Size),
	})
	if err != nil {
		return err.Error(), nil //nolint: nilerr
	}

	return string(resp), nil
}

func (t *ToolOrderList) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        ToolOrderListName,
		Desc:        ToolOrderListDesc,
		ParamsOneOf: ToolOrderListRequestBody,
	}, nil
}
