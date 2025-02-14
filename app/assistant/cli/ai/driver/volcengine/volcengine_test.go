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
	"testing"
	"time"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	arkmodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/utils"
	"github.com/volcengine/volcengine-go-sdk/volcengine"

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
	c := NewClient(&ClientConfig{})
	emptyServerCaller := struct {
		adapter.ServerCaller
	}{}
	c.SetServerCaller(emptyServerCaller)
	completionResp := arkmodel.ChatCompletionResponse{
		Choices: []*arkmodel.ChatCompletionChoice{
			{
				Message: arkmodel.ChatCompletionMessage{
					Content: &arkmodel.ChatCompletionMessageContent{
						StringValue: volcengine.String("completion"),
					},
				},
			},
		},
	}
	stream := utils.ChatCompletionStreamReader{}
	streamResp := arkmodel.ChatCompletionStreamResponse{
		Choices: []*arkmodel.ChatCompletionStreamChoice{
			{
				Delta: arkmodel.ChatCompletionStreamChoiceDelta{
					Content: "stream",
				},
			},
		},
	}
	errInstance := errors.New("an error")
	PatchConvey("Test the volcengine client Call", t, func() {
		PatchConvey("Test when everything is normal(no stream)", func() {
			Mock((*utils.ChatCompletionStreamReader).Recv).Return(streamResp, io.EOF).Build()
			Mock((*utils.ChatCompletionStreamReader).Close).To(func() error { return nil }).Build()
			Mock((*arkruntime.Client).CreateChatCompletion).Return(completionResp, nil).Build()
			Mock((*arkruntime.Client).CreateChatCompletionStream).Return(&stream, nil).Build()
			Mock(c.functionCalling).Return(nil).Build()

			err1, err, res := RunCallTest(c, t)
			So(err1, ShouldBeNil)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "completion")
		})

		PatchConvey("Test when everything is normal(with stream)", func() {
			calls := 0
			Mock((*utils.ChatCompletionStreamReader).Recv).To(func() (arkmodel.ChatCompletionStreamResponse, error) {
				if calls == 0 {
					calls++
					return streamResp, nil
				} else {
					return arkmodel.ChatCompletionStreamResponse{}, io.EOF
				}
			}).Build()
			Mock((*utils.ChatCompletionStreamReader).Close).To(func() error { return nil }).Build()
			Mock((*arkruntime.Client).CreateChatCompletion).Return(completionResp, nil).Build()
			Mock((*arkruntime.Client).CreateChatCompletionStream).Return(&stream, nil).Build()
			Mock(c.functionCalling).Return(nil).Build()

			err1, err, res := RunCallTest(c, t)
			So(err1, ShouldBeNil)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "stream")
		})

		PatchConvey("Test when CreateChatCompletion returns an error", func() {
			Mock((*arkruntime.Client).CreateChatCompletion).
				Return(arkmodel.ChatCompletionResponse{}, errInstance).Build()

			err1, err, res := RunCallTest(c, t)
			So(err1, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "")
		})

		PatchConvey("Test when CreateChatCompletionStream returns an error", func() {
			Mock((*utils.ChatCompletionStreamReader).Close).To(func() error { return nil }).Build()
			Mock((*arkruntime.Client).CreateChatCompletion).Return(completionResp, nil).Build()
			Mock((*arkruntime.Client).CreateChatCompletionStream).
				Return(&utils.ChatCompletionStreamReader{}, errInstance).Build()
			Mock(c.functionCalling).Return(nil).Build()

			err1, err, res := RunCallTest(c, t)
			So(err1, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "")
		})
	})
}
