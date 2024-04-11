package websock

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"work4/biz/handler/websocket"
)

func register(h *server.Hertz) {
	h.GET(`/`, append(_homeMW(), websocket.Chat)...)
}
