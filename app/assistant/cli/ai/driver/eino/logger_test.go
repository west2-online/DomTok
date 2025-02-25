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
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/schema"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/pkg/logger"
)

func RunLogger(cb LoggerCallback,
	input callbacks.CallbackInput,
	output callbacks.CallbackOutput,
) {
	ctx := context.Background()
	info := &callbacks.RunInfo{}
	cb.OnError(cb.OnEnd(cb.OnStart(ctx, info, input), info, output), info, fmt.Errorf("error"))
}

func RunStreamLogger(cb LoggerCallback,
	input *schema.StreamReader[callbacks.CallbackInput],
	output *schema.StreamReader[callbacks.CallbackOutput],
) {
	ctx := context.Background()
	info := &callbacks.RunInfo{}
	cb.OnEndWithStreamOutput(
		cb.OnStartWithStreamInput(ctx, info, input), info, output)
}

func TestLoggerCallback(t *testing.T) {
	cb := &LoggerCallback{}
	input := struct{ callbacks.CallbackInput }{}
	inputStream := schema.StreamReader[callbacks.CallbackInput]{}
	output := struct{ callbacks.CallbackOutput }{}
	outputStream := schema.StreamReader[callbacks.CallbackOutput]{}

	PatchConvey("Test LoggerCallback", t, func() {
		Mock(logger.Infof).To(func(string, ...interface{}) {}).Build()
		Mock(logger.Errorf).To(func(string, ...interface{}) {}).Build()
		Mock(logger.Fatalf).To(func(string, ...interface{}) {}).Build()
		Mock(logger.Warnf).To(func(string, ...interface{}) {}).Build()
		Mock(logger.Debugf).To(func(string, ...interface{}) {}).Build()
		PatchConvey("Test the logger callback on completion", func() {
			RunLogger(*cb, input, output)
			// if no panic, it means the logger callback is correct
		})

		PatchConvey("Test the logger callback on input stream", func() {
			x := make(chan struct{})
			MockGeneric((*schema.StreamReader[any]).Recv).Return(nil, io.EOF).Build()
			MockGeneric((*schema.StreamReader[any]).Close).To(func() {
				x <- struct{}{}
			}).Build()
			wg := sync.WaitGroup{}
			wg.Add(1)
			isInputClosed := false
			go func() {
				ticker := time.NewTicker(time.Second)
				select {
				case <-x:
					isInputClosed = true
				case <-ticker.C:
				}
				wg.Done()
			}()
			go RunStreamLogger(*cb, &inputStream, &outputStream)
			wg.Wait()
			So(isInputClosed, ShouldBeTrue)
		})

		PatchConvey("Test the logger callback on output stream", func() {
			x := make(chan struct{})
			MockGeneric((*schema.StreamReader[any]).Recv).Return(nil, io.EOF).Build()
			MockGeneric((*schema.StreamReader[any]).Close).To(func() {
				x <- struct{}{}
			}).Build()
			wg := sync.WaitGroup{}
			wg.Add(1)
			isOutputClosed := false
			go func() {
				ticker := time.NewTicker(time.Second)
				select {
				case <-x:
					isOutputClosed = true
				case <-ticker.C:
				}
				wg.Done()
			}()
			go RunStreamLogger(*cb, &inputStream, &outputStream)
			wg.Wait()
			So(isOutputClosed, ShouldBeTrue)
		})

		PatchConvey("Test the logger callback on output stream with error", func() {
			x := make(chan struct{})
			MockGeneric((*schema.StreamReader[any]).Recv).Return(nil, fmt.Errorf("")).Build()
			MockGeneric((*schema.StreamReader[any]).Close).To(func() {
				x <- struct{}{}
			}).Build()
			wg := sync.WaitGroup{}
			wg.Add(1)
			isOutputClosed := false
			go func() {
				ticker := time.NewTicker(time.Second)
				select {
				case <-x:
					isOutputClosed = true
				case <-ticker.C:
				}
				wg.Done()
			}()
			go RunStreamLogger(*cb, &inputStream, &outputStream)
			wg.Wait()
			So(isOutputClosed, ShouldBeTrue)
		})
	})
}
