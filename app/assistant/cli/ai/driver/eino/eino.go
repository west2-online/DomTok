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
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/cloudwego/eino/callbacks"
	components "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/adapter"
	category "github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/model"
	"github.com/west2-online/DomTok/pkg/logger"
)

// Client is a client struct for calling the AI
type Client struct {
	adapter.AIClient

	persona string

	caller  category.GetServerCaller
	builder category.BuildChatModel

	tools []tool.BaseTool

	recorder sync.Map
}

// NewClient creates a new eino client
func NewClient() *Client {
	cli := &Client{}
	cli.persona = GetPersona()
	callbacks.InitCallbackHandlers([]callbacks.Handler{agentLogger})
	return cli
}

// SetServerCategory sets the server category
func (c *Client) SetServerCategory(category category.GetServerCaller) {
	c.caller = category
	c.tools = *GetTools(category)
}

// SetBuilder sets the build chat model
func (c *Client) SetBuilder(buildChatModel category.BuildChatModel) {
	c.builder = buildChatModel
}

func (c *Client) BuildChatModel(ctx context.Context) (components.ChatModel, error) {
	return c.builder(ctx)
}

// Call calls the AI service
func (c *Client) Call(ctx context.Context, dialog model.IDialog) (err error) {
	defer dialog.Close()

	if c.caller == nil {
		return fmt.Errorf("server category is not set")
	}

	if c.builder == nil {
		return fmt.Errorf("build chat model is not set")
	}

	c.markDialog(dialog)

	history, err := c.readHistory(dialog.Unique())
	if err != nil {
		return fmt.Errorf("read history failed: %w", err)
	}

	chatModel, err := c.BuildChatModel(ctx)
	if err != nil {
		return fmt.Errorf("build chat model failed: %w", err)
	}

	ra, err := react.NewAgent(ctx, &react.AgentConfig{
		Model: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: c.tools,
		},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			return append(history, input...)
		},
	})
	if err != nil {
		return fmt.Errorf("create agent failed: %w", err)
	}

	history = append(history, schema.UserMessage(dialog.Message()))
	stream, err := ra.Stream(ctx, []*schema.Message{schema.UserMessage(dialog.Message())})
	if err != nil {
		return fmt.Errorf("stream failed: %w", err)
	}
	defer stream.Close()

	out := ""
	for {
		frame, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("stream recv failed: %w", err)
		}

		if len(frame.Content) != 0 {
			dialog.Send(frame.Content)
			out += frame.Content
		}
	}

	history = append(history, schema.AssistantMessage(out, nil))
	c.storeMarkedDialog(dialog, history)

	return nil
}

// markDialog marks the dialog
func (c *Client) markDialog(dialog model.IDialog) {
	_, exist := c.recorder.Load(dialog.Unique())
	if !exist {
		c.recorder.Store(dialog.Unique(), make([]*schema.Message, 0))
	}
}

// storeMarkedDialog stores the dialog that has been marked
func (c *Client) storeMarkedDialog(dialog model.IDialog, messages []*schema.Message) {
	v, exist := c.recorder.Load(dialog.Unique())
	if !exist || v == nil {
		return
	}

	c.recorder.Store(dialog.Unique(), messages)
}

// ForgetDialog forgets the dialog
func (c *Client) ForgetDialog(dialog model.IDialog) {
	c.recorder.Delete(dialog.Unique())
}

// readHistory reads the history of the dialog, if the history is not found, it returns the base request
func (c *Client) readHistory(key string) ([]*schema.Message, error) {
	v, ok := c.recorder.Load(key)
	if !ok {
		return []*schema.Message{schema.SystemMessage(c.persona)}, nil
	}

	history, ok := v.([]*schema.Message)
	if !ok {
		return nil, fmt.Errorf("unexpected type transition")
	}

	if len(history) == 0 {
		history = []*schema.Message{schema.SystemMessage(c.persona)}
	}
	return history, nil
}

var agentLogger = &LoggerCallback{}

type LoggerCallback struct {
	callbacks.HandlerBuilder
}

func (cb *LoggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	logger.Infof("[AI-Agent] input: %#v", input)
	return ctx
}

func (cb *LoggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	logger.Infof("[AI-Agent] output: %#v", output)
	return ctx
}

func (cb *LoggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	logger.Errorf("[AI-Agent Stream] error: %v", err)
	return ctx
}

func (cb *LoggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput],
) context.Context {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Fatalf("internal error: %v", err)
			}
		}()

		defer output.Close() // remember to close the stream in defer

		for {
			_, err := output.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				logger.Errorf("[AI-Agent Stream] error: %v", err)
				return
			}
		}
	}()
	return ctx
}

func (cb *LoggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput],
) context.Context {
	defer input.Close()
	return ctx
}
