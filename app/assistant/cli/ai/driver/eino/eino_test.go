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
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	. "github.com/bytedance/mockey"
	model2 "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
	"github.com/west2-online/DomTok/app/assistant/model"
)

func ReceiveDialog(dialog *model.Dialog, timeout time.Duration, cancel context.CancelFunc, t *testing.T,
) (string, error) {
	ticker := time.NewTicker(timeout)
	res := ""
	for {
		select {
		case <-dialog.NotifyOnClosed():
			return res, nil
		case msg := <-dialog.NotifyOnMessage():
			res += msg
		case <-ticker.C:
			cancel()
			return res, fmt.Errorf("timeout")
		}
	}
}

func RunCallTest(c *Client, t *testing.T) (error, error, string) {
	dialog := model.NewDialog("", "")
	wg := sync.WaitGroup{}
	wg.Add(2)

	var err1 error
	ctx, cancel := context.WithCancel(context.Background())
	var res string
	var err error
	go func(res *string, err *error, wg *sync.WaitGroup) {
		*res, *err = ReceiveDialog(dialog, 1000*time.Millisecond, cancel, t)
		wg.Done()
	}(&res, &err, &wg)
	go func(err *error, wg *sync.WaitGroup) {
		errChan := make(chan error)
		go func() { errChan <- c.Call(context.Background(), dialog) }()
		select {
		case <-ctx.Done():
			*err = fmt.Errorf("timeout")
		case *err = <-errChan:
		}
		wg.Done()
	}(&err1, &wg)

	wg.Wait()
	return err1, err, res
}

func TestClient_Call(t *testing.T) {
	c := NewClient()
	stream := &schema.StreamReader[*schema.Message]{}

	PatchConvey("Test the volcengine client Call", t, func() {
		PatchConvey("Test when no error occurs", func() {
			MockValue(&c.caller).To(func(functionName string) adapter.ServerCaller { return nil })
			MockValue(&c.builder).To(func(ctx context.Context) (model2.ChatModel, error) { return nil, nil })

			index := 0
			MockGeneric((*schema.StreamReader[*schema.Message]).Recv).To(func() (*schema.Message, error) {
				if index == 0 {
					index++
					return &schema.Message{Content: "completion"}, nil
				}
				return nil, io.EOF
			}).Build()
			MockGeneric((*schema.StreamReader[*schema.Message]).Close).To(func() {}).Build()
			Mock(react.NewAgent).Return(nil, nil).Build()
			Mock((*react.Agent).Stream).Return(stream, nil).Build()
			Mock(GetTools).Return(&[]tool.BaseTool{}).Build()
			Mock((*Client).BuildChatModel).Return(nil, nil).Build()

			err1, err, res := RunCallTest(c, t)
			So(err1, ShouldBeNil)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "completion")
		})

		PatchConvey("Test when stream delta contents", func() {
			MockValue(&c.caller).To(func(functionName string) adapter.ServerCaller { return nil })
			MockValue(&c.builder).To(func(ctx context.Context) (model2.ChatModel, error) { return nil, nil })

			helloWorld := "hello world"
			index := 0
			MockGeneric((*schema.StreamReader[*schema.Message]).Recv).To(func() (*schema.Message, error) {
				if index < len(helloWorld) {
					index++
					return &schema.Message{Content: string(helloWorld[index-1])}, nil
				}
				return &schema.Message{Content: "completion"}, io.EOF
			}).Build()
			MockGeneric((*schema.StreamReader[*schema.Message]).Close).To(func() {}).Build()
			Mock(react.NewAgent).Return(nil, nil).Build()
			Mock((*react.Agent).Stream).Return(stream, nil).Build()
			Mock(GetTools).Return(&[]tool.BaseTool{}).Build()
			Mock((*Client).BuildChatModel).Return(nil, nil).Build()

			err1, err, msg := RunCallTest(c, t)
			So(err1, ShouldBeNil)
			So(err, ShouldBeNil)
			So(msg, ShouldEqual, helloWorld)
		})

		PatchConvey("Test when server category is not set", func() {
			MockValue(&c.caller).To(nil)
			MockValue(&c.builder).To(func(ctx context.Context) (model2.ChatModel, error) { return nil, nil })
			err1, _, _ := RunCallTest(c, t)
			So(err1, ShouldNotBeNil)
		})

		PatchConvey("Test when build chat model is not set", func() {
			MockValue(&c.caller).To(func(functionName string) adapter.ServerCaller { return nil })
			MockValue(&c.builder).To(nil)
			err1, _, _ := RunCallTest(c, t)
			So(err1, ShouldNotBeNil)
		})

		PatchConvey("Test when stream recv failed", func() {
			MockValue(&c.caller).To(func(functionName string) adapter.ServerCaller { return nil })
			MockValue(&c.builder).To(func(ctx context.Context) (model2.ChatModel, error) { return nil, nil })

			MockGeneric((*schema.StreamReader[*schema.Message]).Recv).To(func() (*schema.Message, error) {
				return nil, fmt.Errorf("stream recv failed")
			}).Build()
			MockGeneric((*schema.StreamReader[*schema.Message]).Close).To(func() {}).Build()
			Mock(react.NewAgent).Return(nil, nil).Build()
			Mock((*react.Agent).Stream).Return(stream, nil).Build()
			Mock(GetTools).Return(&[]tool.BaseTool{}).Build()
			Mock((*Client).BuildChatModel).Return(nil, nil).Build()

			err1, _, _ := RunCallTest(c, t)
			So(err1, ShouldNotBeNil)
		})
	})
}
