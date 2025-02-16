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

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/schema"

	"github.com/west2-online/DomTok/pkg/logger"
)

type LoggerCallback struct {
	callbacks.HandlerBuilder
}

func (cb *LoggerCallback) OnStart(ctx context.Context, _ *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	logger.Debugf("[AI-Agent] input: %#v", input)
	return ctx
}

func (cb *LoggerCallback) OnEnd(ctx context.Context, _ *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	logger.Debugf("[AI-Agent] output: %#v", output)
	return ctx
}

func (cb *LoggerCallback) OnError(ctx context.Context, _ *callbacks.RunInfo, err error) context.Context {
	logger.Debugf("[AI-Agent Stream] error: %v", err)
	return ctx
}

func (cb *LoggerCallback) OnEndWithStreamOutput(ctx context.Context, _ *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput],
) context.Context {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Debugf("internal error: %v", err)
			}
		}()

		defer output.Close() // remember to close the stream in defer

		for {
			_, err := output.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				logger.Debugf("[AI-Agent Stream] error: %v", err)
				return
			}
		}
	}()
	return ctx
}

func (cb *LoggerCallback) OnStartWithStreamInput(ctx context.Context, _ *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput],
) context.Context {
	defer input.Close()
	return ctx
}
