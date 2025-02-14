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

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	arkmodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/adapter"
	server "github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
	"github.com/west2-online/DomTok/app/assistant/model"
)

// Client is a client struct for volcengine AI service
type Client struct {
	adapter.AIClient

	cli     *arkruntime.Client
	caller  server.ServerCaller
	baseReq *arkmodel.CreateChatCompletionRequest

	recorder sync.Map
}

// ClientConfig is the option for creating a new volcengine client
type ClientConfig struct {
	ApiKey  string
	BaseUrl string
	Region  string
	Model   string
}

// NewClient creates a new volcengine client
func NewClient(opt *ClientConfig) *Client {
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

// SetServerCaller sets the server caller
func (c *Client) SetServerCaller(caller server.ServerCaller) {
	c.caller = caller
}

// Call calls the volcengine AI service
func (c *Client) Call(ctx context.Context, dialog model.IDialog) (err error) {
	defer dialog.Close()

	// use the dialog's unique id as the key
	history, err := c.readHistory(dialog.Unique())
	if err != nil {
		return err
	}

	// build the request
	req := c.buildReq(history)

	// append the user message to the request
	c.appendUserMessage(req, dialog.Message())

	// call the AI service
	// this is the first round of conversation
	resp, err := c.cli.CreateChatCompletion(
		ctx,
		req,
		arkruntime.WithCustomHeader(arkmodel.ClientRequestHeader, dialog.Unique()),
	)
	if err != nil {
		return err
	}

	// function calling or just simple response, we can delay determination
	req.Messages = append(req.Messages, &resp.Choices[0].Message)

	// do function calling
	c.functionCalling(ctx, req, &resp)

	// start the second round of conversation
	// if the first round of conversation is a function calling, the second round will send messages in stream
	// otherwise, the second round will send the first round of conversation
	// tips: if there's no function calling, the stream would not send any message
	stream, err := c.cli.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}

	// a simple impl of judging whether the first round of conversation is a function calling
	isMessageSent := false
	output := ""
	for {
		receive, err := stream.Recv()
		// use io.EOF to judge the end of the stream
		// but in the official document, there's `err == io.EOF` instead of `errors.Is(err, io.EOF)`
		// it may cause some problems (error's unwrapping)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		// if the assistant sends a message and the message is not empty, send the message
		if len(receive.Choices) > 0 {
			output += receive.Choices[0].Delta.Content
			dialog.Send(receive.Choices[0].Delta.Content)
			isMessageSent = true
		}
	}
	_ = stream.Close()

	// if the assistant does not send a message, send the first conversation.
	if !isMessageSent {
		dialog.Send(*resp.Choices[0].Message.Content.StringValue)
		c.recorder.Store(dialog.Unique(), req.Messages)
		return nil
	}

	// if the assistant sends a message, it's clear that the first round of conversation is a function calling
	// therefore, we need to append the assistant message to the request
	c.appendAssistantMessage(req, output)
	c.recorder.Store(dialog.Unique(), req.Messages)

	return nil
}

// ForgetDialog forgets the dialog
func (c *Client) ForgetDialog(dialog model.IDialog) {
	c.recorder.Delete(dialog.Unique())
}

// readHistory reads the history of the dialog, if the history is not found, it returns the base request
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

// buildReq builds the request for the AI service
func (c *Client) buildReq(messages []*arkmodel.ChatCompletionMessage) *arkmodel.CreateChatCompletionRequest {
	return &arkmodel.CreateChatCompletionRequest{
		Model:    c.baseReq.Model,
		Messages: messages,
		Tools:    c.baseReq.Tools,
	}
}

// functionCalling is a function that adapts the function calling
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

// appendUserMessage appends a user message to the request
func (c *Client) appendUserMessage(req *arkmodel.CreateChatCompletionRequest, message string) {
	req.Messages = append(req.Messages, &arkmodel.ChatCompletionMessage{
		Role:    arkmodel.ChatMessageRoleUser,
		Content: &arkmodel.ChatCompletionMessageContent{StringValue: volcengine.String(message)},
	})
}

// appendAssistantMessage appends an assistant message to the request
func (c *Client) appendAssistantMessage(req *arkmodel.CreateChatCompletionRequest, message string) {
	req.Messages = append(req.Messages, &arkmodel.ChatCompletionMessage{
		Role:    arkmodel.ChatMessageRoleAssistant,
		Content: &arkmodel.ChatCompletionMessageContent{StringValue: volcengine.String(message)},
	})
}
