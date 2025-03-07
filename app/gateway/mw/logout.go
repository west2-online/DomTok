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
	"errors"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/DomTok/app/gateway/pack"
	"github.com/west2-online/DomTok/app/gateway/service"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

var (
	GateWayService *service.GateWayService
	once           sync.Once
)

func UserLoginStatus() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		once.Do(func() {
			GateWayService = service.NewGateWayService()
		})
		token := string(c.GetHeader(constants.AuthHeader))
		_, uid, err := utils.CheckToken(token)
		if err != nil {
			pack.RespError(c, err)
			c.Abort()
			return
		}

		// 判断是否是白名单
		if !GateWayService.Bf.Test(uid) {
			// 是的话检查是否退出登录
			if GateWayService.Re.IsUserLogout(ctx, uid) {
				pack.RespError(c, errors.New("please login again"))
				c.Abort()
			}
			c.Next(ctx)
			return
		}

		// 检查是否是是否真的是黑名单中的人
		if GateWayService.Re.IsUserBanned(ctx, uid) {
			pack.RespError(c, errors.New("you are banned"))
			c.Abort()
			return
		}

		// 如果不是的话就尝试继续判断是否登录
		if GateWayService.Re.IsUserLogout(ctx, uid) {
			pack.RespError(c, errors.New("please login again"))
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
