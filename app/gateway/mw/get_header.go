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

package mw

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/DomTok/app/gateway/pack"
	"github.com/west2-online/DomTok/kitex_gen/model"
	metainfoContext "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/errno"
)

// GetHeaderParams 获取请求头的信息，处理 id  并附加到 Context 中
func GetHeaderParams() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		id := string(c.GetHeader("UserId"))

		if id == "" {
			pack.RespError(c, errno.Errorf(errno.ParamMissingHeaderCode, "header is missing the ID field"))
			c.Abort()
			return
		}

		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			pack.RespError(c, errno.Errorf(errno.ParamInvalidHeaderCode, " header ID field is invalid"))
			c.Abort()
			return
		}

		// 实现规范化服务透传，不需要中间进行编解码
		ctx = metainfoContext.WithLoginData(ctx, &model.LoginData{
			UserId: userId,
		})
		c.Next(ctx)
	}
}
