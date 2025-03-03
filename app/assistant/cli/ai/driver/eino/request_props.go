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

package eino

import (
	"github.com/cloudwego/eino/components/tool"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/tools/remote"
)

const (
	roleDescription = `你叫DomTok，是一个电商平台的AI助手`
)

var tools *[]tool.BaseTool

// GetPersona returns the prevMessages
func GetPersona() string {
	return roleDescription
}

// GetTools returns the tools
func GetTools(caller model.GetServerCaller) *[]tool.BaseTool {
	if tools != nil {
		return tools
	}
	RebuildTools(caller)
	return tools
}

// RebuildTools rebuilds the tools
func RebuildTools(caller model.GetServerCaller) {
	ts := BuildTools(caller)
	tools = &ts
}

// BuildTools builds the tools
func BuildTools(strategy model.GetServerCaller) []tool.BaseTool {
	return []tool.BaseTool{
		remote.CartShow(strategy),
		remote.OrderCreate(strategy),
		remote.OrderList(strategy),
		remote.OrderView(strategy),
	}
}
