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
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	arkmodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/adapter"
	server "github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
	"github.com/west2-online/DomTok/app/assistant/model"
)

type Client struct {
	adapter.AIClient

	cli     *arkruntime.Client
	caller  server.ServerCaller
	baseReq *arkmodel.CreateChatCompletionRequest

	recorder sync.Map
}

type ClientOption struct {
	ApiKey  string
	BaseUrl string
	Region  string
	Model   string
}

func NewClient(opt *ClientOption) *Client {
	cli := arkruntime.NewClientWithApiKey(
		opt.ApiKey,
		arkruntime.WithBaseUrl(opt.BaseUrl),
		arkruntime.WithRegion(opt.Region),
	)
	baseReq := &arkmodel.CreateChatCompletionRequest{
		Model:    opt.Model,
		Messages: *GetPrevMessages(),
		Tools:    *GetTools(),
	}
	return &Client{cli: cli, baseReq: baseReq}
}

func (c *Client) SetServerCaller(caller server.ServerCaller) {
	c.caller = caller
}

func (c *Client) Call(ctx context.Context, dialog model.IDialog) (err error) {
	defer dialog.Close()

	history, err := c.readHistory(dialog.Unique())
	if err != nil {
		return err
	}

	req := c.buildReq(history)

	c.appendUserMessage(req, dialog.Message())

	hlog.Info(dialog.Unique())
	for _, m := range req.Messages {
		hlog.Info(m.Role, " ", *m.Content.StringValue)
	}

	resp, err := c.cli.CreateChatCompletion(
		ctx,
		req,
		arkruntime.WithCustomHeader(arkmodel.ClientRequestHeader, dialog.Unique()),
	)
	if err != nil {
		return err
	}

	req.Messages = append(req.Messages, &resp.Choices[0].Message)
	c.functionCalling(ctx, req, &resp)

	stream, err := c.cli.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}

	isMessageSent := false
	output := ""
	for {
		receive, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if len(receive.Choices) > 0 {
			output += receive.Choices[0].Delta.Content
			dialog.Send(receive.Choices[0].Delta.Content)
			isMessageSent = true
		}
	}
	_ = stream.Close()

	// If the assistant does not send a message, send the first choice.
	if !isMessageSent {
		dialog.Send(*resp.Choices[0].Message.Content.StringValue)
		c.recorder.Store(dialog.Unique(), req.Messages)
		return nil
	}

	c.appendAssistantMessage(req, output)
	c.recorder.Store(dialog.Unique(), req.Messages)

	return nil
}

func (c *Client) ForgetDialog(dialog model.IDialog) {
	c.recorder.Delete(dialog.Unique())
}

func (c *Client) readHistory(key string) ([]*arkmodel.ChatCompletionMessage, error) {
	v, ok := c.recorder.Load(key)
	if !ok {
		return c.baseReq.Messages, nil
	}

	history, ok := v.([]*arkmodel.ChatCompletionMessage)
	if !ok {
		return nil, fmt.Errorf("unexpected type transition")
	}

	if len(history) == 0 {
		history = c.baseReq.Messages
	}

	return history, nil
}

func (c *Client) buildReq(messages []*arkmodel.ChatCompletionMessage) *arkmodel.CreateChatCompletionRequest {
	return &arkmodel.CreateChatCompletionRequest{
		Model:    c.baseReq.Model,
		Messages: messages,
		Tools:    c.baseReq.Tools,
	}
}

func (c *Client) functionCalling(
	ctx context.Context,
	req *arkmodel.CreateChatCompletionRequest,
	resp *arkmodel.ChatCompletionResponse,
) {
	for _, toolCall := range resp.Choices[0].Message.ToolCalls {
		toolCallResult, err := CallFunction(
			ctx,
			toolCall.Function.Name,
			toolCall.Function.Arguments,
			c.caller,
		)
		if err != nil {
			toolCallResult = err.Error()
		}

		req.Messages = append(req.Messages, &arkmodel.ChatCompletionMessage{
			Role:       arkmodel.ChatMessageRoleTool,
			ToolCallID: toolCall.ID,
			Content: &arkmodel.ChatCompletionMessageContent{
				StringValue: &toolCallResult,
			},
		})
	}
}

func (c *Client) appendUserMessage(req *arkmodel.CreateChatCompletionRequest, message string) {
	req.Messages = append(req.Messages, &arkmodel.ChatCompletionMessage{
		Role:    arkmodel.ChatMessageRoleUser,
		Content: &arkmodel.ChatCompletionMessageContent{StringValue: volcengine.String(message)},
	})
}

func (c *Client) appendAssistantMessage(req *arkmodel.CreateChatCompletionRequest, message string) {
	req.Messages = append(req.Messages, &arkmodel.ChatCompletionMessage{
		Role:    arkmodel.ChatMessageRoleAssistant,
		Content: &arkmodel.ChatCompletionMessageContent{StringValue: volcengine.String(message)},
	})
}
