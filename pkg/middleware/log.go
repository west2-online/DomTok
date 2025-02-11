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

package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"go.uber.org/zap"

	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

// ErrorLog WARNING: 如果使用该中间件, 请务必注册于 Respond 之前
//
// # ErrorLog 会在 rpc 处理完请求之后对错误进行截取, 并进行日志的输出
//
// ErrorLog 流程如下
//  1. 如果错误为 nil 就直接返回
//  2. 尝试将 err 转换为 errno, 如果成功就再判断是否需要输出调用栈
//  3. 如果不是 errno, 说明这是一个不可控的 error, 但这里不对他进行处理, 打印返回即可
func ErrorLog() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) error {
			err := next(ctx, req, resp) // next 就是继续往深处去, 一直到处理我们的业务请求
			if err == nil {
				return nil
			}

			var e errno.ErrNo
			if errors.As(err, &e) {
				logFields := []zap.Field{
					zap.Int(constants.ErrorCodeKey, int(e.ErrorCode)),
				}
				if stackTrace := e.StackTrace(); stackTrace != nil {
					logFields = append(logFields, zap.String(constants.StackTraceKey, fmt.Sprintf("%+v", stackTrace)))
				}
				logger.Error(err.Error(), logFields...)
			} else {
				logger.Errorf("unexpected system error, err: %v", err)
			}

			return err
		}
	}
}
