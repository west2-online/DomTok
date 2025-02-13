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
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

const (
	roleDescription = `你叫DomTok，是一个电商平台的AI助手`
)

var (
	prevMessages *[]*model.ChatCompletionMessage
	tools        *[]*model.Tool
)

// GetPrevMessages returns the prevMessages
func GetPrevMessages() *[]*model.ChatCompletionMessage {
	if prevMessages != nil {
		return prevMessages
	}
	RebuildPrevMessages()
	return prevMessages
}

// RebuildPrevMessages rebuilds the prevMessages
func RebuildPrevMessages() {
	prevMessages = &[]*model.ChatCompletionMessage{
		{
			Role: model.ChatMessageRoleSystem,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(roleDescription),
			},
		},
	}
}

// GetTools returns the tools
func GetTools() *[]*model.Tool {
	if tools != nil {
		return tools
	}
	RebuildTools()
	return tools
}

// RebuildTools rebuilds the tools
func RebuildTools() {
	tools = nil
	RebuildFunctions()
	ts := make([]*model.Tool, len(*functions))
	for i, f := range *functions {
		ts[i] = f.AsTool()
	}
	tools = &ts
}
