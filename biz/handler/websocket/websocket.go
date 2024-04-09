package websocket

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"
	"strconv"
	"work4/biz/service"
)

var upgrader = websocket.HertzUpgrader{}

// Chat .
// @router / [GET]
func Chat(ctx context.Context, c *app.RequestContext) {
	var err error
	hlog.Infof("success")
	err = upgrader.Upgrade(c, func(conn *websocket.Conn) {
		uid := strconv.FormatInt(service.GetUidFormContext(c), 10)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("BadConnection"))
			return
		}
		conn.WriteMessage(websocket.TextMessage, []byte(`Welcome, `+uid))

		s := service.NewChatService(ctx, c, conn)

		if err := s.Login(); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("BadConnection"))
			return
		}
		defer s.Logout()

		if err := s.ReadOfflineMessage(); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("BadConnection"))
			return
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("BadConnection"))
				return
			}

			if err := s.SendMessage(message); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("BadConnection"))
				return
			}
		}
	})

	if err != nil {
		c.JSON(consts.StatusOK, `error`)
		return
	}
}
