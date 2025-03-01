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

func ReceiveDialog(dialog *model.Dialog, timeout time.Duration, cancel context.CancelFunc,
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

func RunCallTest(c *Client) (error, error, string) {
	dialog := model.NewDialog("", "")
	wg := sync.WaitGroup{}
	wg.Add(2)

	var err1 error
	ctx, cancel := context.WithCancel(context.Background())
	var res string
	var err error
	go func(res *string, err *error, wg *sync.WaitGroup) {
		*res, *err = ReceiveDialog(dialog, 1000*time.Millisecond, cancel)
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

func ClientCallNormalize(c *Client, stream *schema.StreamReader[*schema.Message]) {
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
}

func TestClient_Call(t *testing.T) {
	c := NewClient()
	stream := &schema.StreamReader[*schema.Message]{}

	PatchConvey("Test Call", t, func() {
		ClientCallNormalize(c, stream)
		PatchConvey("Test when no error occurs", func() {
			err1, err, msg := RunCallTest(c)
			So(err1, ShouldBeNil)
			So(err, ShouldBeNil)
			So(msg, ShouldEqual, "completion")
		})

		PatchConvey("Test when stream delta contents", func() {
			helloWorld := "hello world"
			index := 0
			MockGeneric((*schema.StreamReader[*schema.Message]).Recv).To(func() (*schema.Message, error) {
				if index < len(helloWorld) {
					index++
					return &schema.Message{Content: string(helloWorld[index-1])}, nil
				}
				return &schema.Message{Content: "completion"}, io.EOF
			}).Build()

			err1, err, msg := RunCallTest(c)
			So(err1, ShouldBeNil)
			So(err, ShouldBeNil)
			So(msg, ShouldEqual, helloWorld)
		})

		PatchConvey("Test when server category is not set", func() {
			MockValue(&c.caller).To(nil)
			err1, _, _ := RunCallTest(c)
			So(err1, ShouldNotBeNil)
		})

		PatchConvey("Test when build chat model is not set", func() {
			MockValue(&c.builder).To(nil)
			err1, _, _ := RunCallTest(c)
			So(err1, ShouldNotBeNil)
		})

		PatchConvey("Test when read history failed", func() {
			Mock((*Client).readHistory).Return(nil, fmt.Errorf("read history failed")).Build()
			err1, _, _ := RunCallTest(c)
			So(err1, ShouldNotBeNil)
		})

		PatchConvey("Test when build chat model failed", func() {
			Mock((*Client).BuildChatModel).Return(nil, fmt.Errorf("build chat model failed")).Build()
			err1, _, _ := RunCallTest(c)
			So(err1, ShouldNotBeNil)
		})

		PatchConvey("Test when create agent failed", func() {
			Mock(react.NewAgent).Return(nil, fmt.Errorf("create agent failed")).Build()
			err1, _, _ := RunCallTest(c)
			So(err1, ShouldNotBeNil)
		})

		PatchConvey("Test when create stream failed", func() {
			Mock((*react.Agent).Stream).Return(nil, fmt.Errorf("create stream failed")).Build()
			err1, _, _ := RunCallTest(c)
			So(err1, ShouldNotBeNil)
		})

		PatchConvey("Test when stream recv failed", func() {
			MockGeneric((*schema.StreamReader[*schema.Message]).Recv).To(func() (*schema.Message, error) {
				return nil, fmt.Errorf("stream recv failed")
			}).Build()

			err1, _, _ := RunCallTest(c)
			So(err1, ShouldNotBeNil)
		})
	})
}

func TestClient_SetServerCategory(t *testing.T) {
	cli := NewClient()
	defaultCaller := cli.caller
	PatchConvey("Test SetServerStrategy", t, func() {
		Mock(GetTools).Return(&[]tool.BaseTool{}).Build()
		PatchConvey("Test when server category is set", func() {
			cli.SetServerStrategy(func(functionName string) adapter.ServerCaller { return nil })
			So(cli.caller, ShouldNotEqual, defaultCaller)
		})
	})
}

func TestClient_SetBuilder(t *testing.T) {
	cli := NewClient()
	defaultBuilder := cli.builder
	PatchConvey("Test SetBuilder", t, func() {
		PatchConvey("Test when build chat model is set", func() {
			cli.SetBuilder(func(ctx context.Context) (model2.ChatModel, error) { return nil, nil })
			So(cli.builder, ShouldNotEqual, defaultBuilder)
		})
	})
}

func TestClient_BuildChatModel(t *testing.T) {
	cli := NewClient()
	type X struct {
		model2.ChatModel
	}
	m := struct {
		model2.ChatModel
		x X
	}{}
	testBuilder := func(ctx context.Context) (model2.ChatModel, error) { return m, nil }
	PatchConvey("Test BuildChatModel", t, func() {
		cli.SetBuilder(testBuilder)
		PatchConvey("Test when build chat model is set", func() {
			chatModel, err := cli.BuildChatModel(context.Background())
			So(chatModel, ShouldEqual, m)
			So(err, ShouldBeNil)
		})
	})
}

func IsDialogExist(cli *Client, dialog *model.Dialog) bool {
	exist := false
	cli.recorder.Range(func(key, value interface{}) bool {
		if key == dialog.Unique() {
			exist = true
		}
		return true
	})

	return exist
}

func TestClient_ForgetDialog(t *testing.T) {
	cli := NewClient()
	d1 := model.NewDialog("1", "")
	d2 := model.NewDialog("2", "")
	PatchConvey("Test ForgetDialog", t, func() {
		PatchConvey("Test when dialog is not exist", func() {
			cli.ForgetDialog(d1)
			So(IsDialogExist(cli, d1), ShouldBeFalse)
			So(IsDialogExist(cli, d2), ShouldBeFalse)
		})

		PatchConvey("Test when dialog is exist", func() {
			cli.recorder.Store(d1.Unique(), nil)
			cli.recorder.Store(d2.Unique(), nil)
			cli.ForgetDialog(d1)
			So(IsDialogExist(cli, d1), ShouldBeFalse)
			So(IsDialogExist(cli, d2), ShouldBeTrue)
		})
	})
}

func TestClient_Recorder(t *testing.T) {
	cli := NewClient()
	d1 := model.NewDialog("1", "")
	d2 := model.NewDialog("2", "")
	PatchConvey("Test Recorder", t, func() {
		PatchConvey("Test when logic is normal", func() {
			cli.markDialog(d1)
			cli.storeMarkedDialog(d1, nil)
			So(IsDialogExist(cli, d1), ShouldBeTrue)

			cli.ForgetDialog(d1)
			So(IsDialogExist(cli, d1), ShouldBeFalse)
		})

		PatchConvey("Test when dialog is not marked", func() {
			cli.storeMarkedDialog(d1, nil)
			So(IsDialogExist(cli, d1), ShouldBeFalse)

			cli.ForgetDialog(d1)
			So(IsDialogExist(cli, d1), ShouldBeFalse)
		})

		PatchConvey("Test dialog cannot impact each other", func() {
			cli.markDialog(d1)
			cli.storeMarkedDialog(d1, nil)
			So(IsDialogExist(cli, d1), ShouldBeTrue)

			cli.ForgetDialog(d2)
			So(IsDialogExist(cli, d1), ShouldBeTrue)
		})

		PatchConvey("Test when forget dialog first", func() {
			cli.ForgetDialog(d1)
			cli.storeMarkedDialog(d1, nil)
			So(IsDialogExist(cli, d1), ShouldBeFalse)
		})
	})
}

func TestClient_ReadHistory(t *testing.T) {
	cli := NewClient()
	PatchConvey("Test ReadHistory", t, func() {
		PatchConvey("Test when dialog is not exist", func() {
			history, err := cli.readHistory("1")
			So(history, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		PatchConvey("Test when dialog is exist but unexpected type", func() {
			cli.recorder.Store("1", "1")
			history, err := cli.readHistory("1")
			So(history, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}
