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

package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"

	"github.com/west2-online/DomTok/app/assistant/pack"
	"github.com/west2-online/DomTok/app/assistant/service"
)

var upgrader = websocket.HertzUpgrader{}

func Entrypoint(ctx context.Context, c *app.RequestContext) {
	// upgrade the protocol to websocket
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		// assign id to ctx
		ctx = context.WithValue(ctx, service.CtxKeyID, pack.GenerateUUID())

		// although the service is like a non-stateful service, we still need to log in
		// in this case, we need to log in to check some args is properly set
		err := service.Service.Login(ctx)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
		// start to accept the message
		for {
			errOnAccept := service.Service.Accept(conn, ctx)
			if errOnAccept != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(errOnAccept.Error()))
				break
			}
		}
		// although the service is like a non-stateful service, we still need to log out
		// in this case, we need to log out to clean up the dialog
		// in order to avoid the over-accumulation of memory
		err = service.Service.Logout(ctx)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, err)
			return
		}
	})
	// handle the error of upgrading the protocol
	if err != nil {
		c.JSON(consts.StatusInternalServerError, err)
		return
	}
}
