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
	"github.com/west2-online/DomTok/pkg/errno"
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
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, err.Error())
	}

	c.markDialog(dialog)

	history, err := c.readHistory(dialog.Unique())
	if err != nil {
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, err.Error())
	}

	chatModel, err := c.BuildChatModel(ctx)
	if err != nil {
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, err.Error())
	}

	ra, err := react.NewAgent(ctx, &react.AgentConfig{
		Model:                 chatModel,
		ToolsConfig:           compose.ToolsNodeConfig{Tools: c.tools},
		MessageModifier:       func(_ context.Context, input []*schema.Message) []*schema.Message { return append(history, input...) },
		StreamToolCallChecker: streamToolCallCheckerStrict,
	})
	if err != nil {
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, err.Error())
	}

	out, err := c.readStreamWithDialog(ctx, ra, dialog)
	if err != nil {
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, err.Error())
	}

	history = append(history, schema.UserMessage(dialog.Message()), schema.AssistantMessage(out, nil))
	c.storeMarkedDialog(dialog, history)

	return nil
}

// checkCallerAndBuilder checks if the caller and builder are set
func (c *Client) checkCallerAndBuilder() error {
	if c.caller == nil {
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "server strategy is not set")
	}

	if c.builder == nil {
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "builder is not set")
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
		return nil, errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "unexpected type transition")
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
		return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, err.Error())
	}
	defer stream.Close()

	out := ""
	for {
		frame, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", errno.NewErrNoWithStack(errno.InternalServiceErrorCode, err.Error())
		}

		if len(frame.Content) != 0 {
			dialog.Send(frame.Content)
			out += frame.Content
		}
	}
	return out, nil
}

func streamToolCallCheckerStrict(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	defer sr.Close()

	frame, err := sr.Recv()
	if err != nil {
		return false, err
	}

	// doubao 1.5 pro 32k:
	// a tool call stream result can be like:
	//   frame.toolCalls = []
	//   frame.content = ""
	// so default checker may not work
	// func firstChunkStreamToolCallChecker(_ context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	//	 defer sr.Close()
	//
	//	 msg, err := sr.Recv()
	//	 if err != nil {
	//	 	 return false, err
	//	 }
	//
	//	 if len(msg.ToolCalls) == 0 {
	//		 return false, nil
	//	 }
	//
	//	 return true, nil
	// }
	// an easy way to fix this is to check if the content is empty:
	//   - ai always response tool calls with empty content
	//     - if content is empty, it is a tool call stream
	//     - if tool calls is not empty, it is a tool call stream
	//   - ai always response message with content, so if content is not empty, it is a message stream
	// only try with doubao 1.5 pro 32k, other models may not work or need modification
	if len(frame.ToolCalls) != 0 || len(frame.Content) == 0 {
		return true, nil
	}

	return false, nil
}
