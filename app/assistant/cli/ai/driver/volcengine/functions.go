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

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

var functions *[]Function

type Function struct {
	Name        string
	Description string
	Properties  []Property
	Call        func(args string, server adapter.ServerCaller) (string, error)
}

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

type Property struct {
	Type        string
	Name        string
	Description string
	Required    bool
}

func GetFunctions() *[]Function {
	if functions != nil {
		return functions
	}
	RebuildFunctions()
	return functions
}

func RebuildFunctions() {
	fs := BuildFunctions()
	functions = &fs
}

func CallFunction(ctx context.Context, name, args string, server adapter.ServerCaller) (string, error) {
	for _, f := range *GetFunctions() {
		if f.Name == name {
			return f.Call(args, server)
		}
	}
	return "", fmt.Errorf("function %s not found", name)
}

func BuildFunctions() []Function {
	return []Function{
		{
			Name:        "ping",
			Description: "测试服务器是否正在运行",
			Properties:  []Property{},
			Call:        PingFunction,
		},
	}
}

func PingFunction(args string, server adapter.ServerCaller) (string, error) {
	resp, err := server.Ping(context.Background())
	if err != nil {
		return "", err
	}

	return string(resp), nil
}
