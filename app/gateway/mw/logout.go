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
	"fmt"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/DomTok/app/gateway/pack"
	"github.com/west2-online/DomTok/app/gateway/service"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
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

		if GateWayService.Bf.Test([]byte(fmt.Sprintf("%d", uid))) {
			pack.RespError(c, errno.Errorf(errno.UserBaned, "bf:user:%d baned, but get request", uid))
			c.Abort()
			return
		}
		ban := GateWayService.Re.IsUserBanned(ctx, uid)
		if ban {
			GateWayService.Bf.Add([]byte(fmt.Sprintf("%d", uid)))
			pack.RespError(c, errno.Errorf(errno.UserBaned, "redis: user:%d baned, but get request", uid))
			c.Abort()
			return
		}
		logout := GateWayService.Re.IsUserLogout(ctx, uid)
		if !logout {
			pack.RespError(c, errno.Errorf(errno.UserLogOut, "user:%d logout, but get request", uid))
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
