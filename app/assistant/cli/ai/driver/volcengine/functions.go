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
	"context"
	"fmt"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/volcengine/tools/remote"
	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

var functions *[]Function

type Function struct {
	Name        string
	Description string
	Properties  []Property
	Call        func(ctx context.Context, args string, server adapter.ServerCaller) (string, error)
}

// AsTool converts the function to a tool
func (f Function) AsTool() *model.Tool {
	type Params struct {
		Type       string
		Properties map[string]interface{}
		Required   []string
	}

	params := Params{
		Type:       "object",
		Properties: make(map[string]interface{}),
		Required:   make([]string, 0),
	}
	for _, p := range f.Properties {
		params.Properties[p.Name] = map[string]interface{}{
			"type":        p.Type,
			"description": p.Description,
		}
		if p.Required {
			params.Required = append(params.Required, p.Name)
		}
	}

	return &model.Tool{
		Type: model.ToolTypeFunction,
		Function: &model.FunctionDefinition{
			Name:        f.Name,
			Description: f.Description,
			Parameters:  params,
		},
	}
}

// Property represents a function parameter
// 打算让ServerCaller自行依赖gateway的request结构体
// 因此这里的Property不需要考虑嵌套
type Property struct {
	Type        string
	Name        string
	Description string
	Required    bool
}

// GetFunctions returns the functions
func GetFunctions() *[]Function {
	if functions != nil {
		return functions
	}
	RebuildFunctions()
	return functions
}

// RebuildFunctions rebuilds the functions
func RebuildFunctions() {
	fs := BuildFunctions()
	functions = &fs
}

// CallFunction calls a function by name
func CallFunction(ctx context.Context, name, args string, server adapter.ServerCaller) (string, error) {
	for _, f := range *GetFunctions() {
		if f.Name == name {
			return f.Call(ctx, args, server)
		}
	}
	return "", fmt.Errorf("function %s not found", name)
}

// BuildFunctions builds the functions
func BuildFunctions() []Function {
	return []Function{
		{
			Name:        "ping",
			Description: "测试服务器是否正在运行",
			Properties:  []Property{},
			Call:        remote.Ping,
		},
	}
}
