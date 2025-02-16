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
	"github.com/west2-online/DomTok/pkg/constants"
)

// Auth 负责校验用户身份，会提取 token 并做处理，Next 时会携带 token 类型
func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.GetHeader(constants.AuthHeader))
		_, err := CheckToken(token)
		if err != nil {
			pack.RespError(c, err)
			c.Abort()
			return
		}

		access, refresh, err := CreateAllToken()
		if err != nil {
			pack.RespError(c, err)
			c.Abort()
			return
		}

		c.Header(constants.AccessTokenHeader, access)
		c.Header(constants.RefreshTokenHeader, refresh)
		c.Next(ctx)
	}
}
