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
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"

	"github.com/west2-online/DomTok/app/assistant/model"
	"github.com/west2-online/DomTok/app/assistant/pack"
	"github.com/west2-online/DomTok/app/assistant/service"
	"github.com/west2-online/DomTok/pkg/constants"
)

var upgrader = websocket.HertzUpgrader{}

func Entrypoint(ctx context.Context, c *app.RequestContext) {
	token := string(c.GetHeader(constants.AccessTokenHeader))
	// upgrade the protocol to websocket
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		if token == "" {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("missing token in header"))
			return
		}

		// generate a uuid
		id := pack.GenerateUUID()

		// assign user info to ctx
		ctx = context.WithValue(ctx, service.CtxKeyID, id)
		ctx = context.WithValue(ctx, service.CtxKeyAccessToken, token)

		// although the service is like a non-stateful service, we still need to log in
		// in this case, we need to log in to check some args is properly set
		err := service.Service.Login(ctx)
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		// assign turn to ctx
		turn := int64(0)
		ctx = context.WithValue(ctx, service.CtxKeyTurn, turn)

		// test if the connection is valid to send the message
		err = conn.WriteMessage(websocket.TextMessage,
			pack.ResponseFactory.ConnectSuccess(model.NewConnectSuccess(id, time.Now().Local().String())))
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		// start to accept the message
		for {
			// accept the message
			errOnAccept := service.Service.Accept(conn, ctx)
			if errOnAccept != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(errOnAccept.Error()))
				break
			}

			// increase the turn
			turn += 1
			ctx = context.WithValue(ctx, service.CtxKeyTurn, turn)
		}
		// although the service is like a non-stateful service, we still need to log out
		// in this case, we need to log out to clean up the dialog
		// in order to avoid the over-accumulation of memory
		err = service.Service.Logout(ctx)
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
	})
	// handle the error of upgrading the protocol
	if err != nil {
		c.JSON(consts.StatusInternalServerError, err)
		return
	}
}
