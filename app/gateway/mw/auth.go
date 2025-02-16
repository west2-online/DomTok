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

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/DomTok/app/gateway/pack"
	"github.com/west2-online/DomTok/kitex_gen/model"
	metainfoContext "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

// Auth 负责校验用户身份，会提取 token 并做处理，Next 时会携带 token 类型
func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.GetHeader(constants.AuthHeader))
		_, uid, err := utils.CheckToken(token)
		if err != nil {
			pack.RespError(c, err)
			c.Abort()
			return
		}

		access, refresh, err := utils.CreateAllToken(uid)
		if err != nil {
			pack.RespError(c, err)
			c.Abort()
			return
		}

		// 实现规范化服务透传，不需要中间进行编解码
		ctx = metainfoContext.WithLoginData(ctx, &model.LoginData{
			UserId: uid,
		})

		c.Header(constants.AccessTokenHeader, access)
		c.Header(constants.RefreshTokenHeader, refresh)
		c.Next(ctx)
	}
}
