package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"
	"github.com/west2-online/DomTok/app/assistant/service"
)

var upgrader = websocket.HertzUpgrader{}

func Entrypoint(ctx context.Context, c *app.RequestContext) {
	// upgrade the protocol to websocket
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		s := service.NewService()
		for {
			err := s.Accept(conn)
			if err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				return
			}
		}
	})

	// handle the error
	if err != nil {
		c.JSON(consts.StatusInternalServerError, err)
		return
	}
}
