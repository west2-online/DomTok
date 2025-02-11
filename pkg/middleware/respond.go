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

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/utils/kitexutil"

	"github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

type response interface {
	GetResult() interface{}
	IsSetSuccess() bool
}

type baser interface {
	IsSetBase() bool
	GetBase() *model.BaseResp
	SetBase(base *model.BaseResp)
}

// Respond 会对所有的响应和 error 进行拦截
// 主要用于为 response 中的 BaseResponse 加上 code 和 msg (仅当 response 未设置且 error 为 errno)
// 让业务层不需要处理相关的操作, 只需要返回 error 即可.
//
// Respond 流程如下:
//  1. 尝试从 resp 中获取到真正的在 idl 中定义的 response (Ref: https://www.cloudwego.io/zh/docs/kitex/tutorials/framework-exten/middleware/)
//  2. 获取到 response 后尝试从 response 中提取出来 model.BaseResp, 然后对其进行赋值
//  3. 尝试判断 err 是否为 errno, 如果是的话说明这是一个可控的 error, 我们对外部返回 nil 即可 (// TODO explain)
func Respond() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) error {
			err := next(ctx, req, resp)

			packResp, ok := resp.(response)
			if !ok || !packResp.IsSetSuccess() {
				return err
			}

			var res baser
			res, ok = packResp.GetResult().(baser)
			if !ok {
				server, _ := kitexutil.GetCaller(ctx)
				method, _ := kitexutil.GetMethod(ctx)
				logger.Errorf("the response struct of the %s method of the %s service does not contain model.BaseResp", method, server)
				return err
			}

			if !res.IsSetBase() {
				res.SetBase(&model.BaseResp{})
			}

			base := res.GetBase()
			var e errno.ErrNo
			if err == nil && base.Code == 0 {
				base.Code = errno.SuccessCode // rpc 接口在没有错误的时候, 就返回 errno.SuccessCode
				base.Msg = errno.SuccessMsg
				return nil
			} else if errors.As(err, &e) {
				if base.Code == 0 {
					base.Code = e.ErrorCode
					base.Msg = e.ErrorMsg
				}
				return nil
			}

			return err
		}
	}
}
