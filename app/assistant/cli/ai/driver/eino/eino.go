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
	strategy "github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/model"
)

// Client is a client struct for calling the AI
type Client struct {
	adapter.AIClient

	persona string

	caller  strategy.GetServerCaller
	builder strategy.BuildChatModel

	tools []tool.BaseTool

	recorder sync.Map
}

// NewClient creates a new eino client
func NewClient() *Client {
	cli := &Client{}
	cli.persona = GetPersona()
	callbacks.InitCallbackHandlers([]callbacks.Handler{&LoggerCallback{}})
	return cli
}

// SetServerStrategy sets the server strategy
func (c *Client) SetServerStrategy(strategy strategy.GetServerCaller) {
	c.caller = strategy
	c.tools = *GetTools(strategy)
}

// SetBuilder sets the build chat model
func (c *Client) SetBuilder(buildChatModel strategy.BuildChatModel) {
	c.builder = buildChatModel
}

func (c *Client) BuildChatModel(ctx context.Context) (components.ChatModel, error) {
	return c.builder(ctx)
}

// Call calls the AI service
func (c *Client) Call(ctx context.Context, dialog model.IDialog) (err error) {
	defer dialog.Close()

	err = c.checkCallerAndBuilder()
	if err != nil {
		return fmt.Errorf("failed to continue: %w", err)
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
		Model:           chatModel,
		ToolsConfig:     compose.ToolsNodeConfig{Tools: c.tools},
		MessageModifier: func(_ context.Context, input []*schema.Message) []*schema.Message { return append(history, input...) },
	})
	if err != nil {
		return fmt.Errorf("create agent failed: %w", err)
	}

	out, err := c.readStreamWithDialog(ctx, ra, dialog)
	if err != nil {
		return fmt.Errorf("read stream failed: %w", err)
	}

	history = append(history, schema.UserMessage(dialog.Message()), schema.AssistantMessage(out, nil))
	c.storeMarkedDialog(dialog, history)

	return nil
}

// checkCallerAndBuilder checks if the caller and builder are set
func (c *Client) checkCallerAndBuilder() error {
	if c.caller == nil {
		return fmt.Errorf("server category is not set")
	}

	if c.builder == nil {
		return fmt.Errorf("build chat model is not set")
	}

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

func (c *Client) readStreamWithDialog(
	ctx context.Context,
	ra *react.Agent,
	dialog model.IDialog,
) (string, error) {
	stream, err := ra.Stream(ctx, []*schema.Message{schema.UserMessage(dialog.Message())})
	if err != nil {
		return "", fmt.Errorf("stream failed: %w", err)
	}
	defer stream.Close()

	out := ""
	for {
		frame, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("stream recv failed: %w", err)
		}

		if len(frame.Content) != 0 {
			dialog.Send(frame.Content)
			out += frame.Content
		}
	}
	return out, nil
}
