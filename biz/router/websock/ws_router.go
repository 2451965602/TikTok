package websock

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"tiktok/biz/handler/websocket"
)

func register(h *server.Hertz) {
	h.GET(`/ws`, append(_homeMW(), websocket.Chat)...)
}
