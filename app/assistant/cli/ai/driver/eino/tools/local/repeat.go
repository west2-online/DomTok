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

package local

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// Tips: This function should not be used in the future

type ToolRepeat struct {
	tool.InvokableTool
}

const (
	ToolRepeatName = "repeat"
	ToolRepeatDesc = "重复用户的输入"
)

type _RepeatArgs struct {
	Message string `json:"message"`
}

func Repeat() *ToolRepeat {
	return &ToolRepeat{}
}

func (t *ToolRepeat) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	args := &_RepeatArgs{}
	err := json.Unmarshal([]byte(argumentsInJSON), args)
	if err != nil {
		return "", err
	}
	return args.Message, nil
}

func (t *ToolRepeat) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: ToolRepeatName,
		Desc: ToolRepeatDesc,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"message": {
				Type:     schema.String,
				Desc:     "要重复的消息",
				Required: true,
			},
		}),
	}, nil
}
