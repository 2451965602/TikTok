package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	"strconv"
	"time"
	"work4/pkg/constants"
	"work4/pkg/errmsg"

	"work4/biz/dal/db"
)

type ChatService struct {
	ctx  context.Context
	c    *app.RequestContext
	conn *websocket.Conn
}

type _user struct {
	username string
	conn     *websocket.Conn
}

var userMap = make(map[string]*_user)

func NewChatService(ctx context.Context, c *app.RequestContext, conn *websocket.Conn) *ChatService {
	return &ChatService{ctx: ctx, c: c, conn: conn}
}

func (s ChatService) Login() error {
	uid := strconv.FormatInt(GetUidFormContext(s.c), 10)

	user, err := db.GetInfo(s.ctx, uid)
	if err != nil {
		return err
	}
	userMap[uid] = &_user{conn: s.conn, username: user.Username}

	return nil
}

func (s ChatService) Logout() {
	uid := strconv.FormatInt(GetUidFormContext(s.c), 10)
	userMap[uid] = nil
}

func (s ChatService) SendMessage(content []byte) error {
	from := strconv.FormatInt(GetUidFormContext(s.c), 10)
	to := s.c.Query(`to_user_id`)

	uid, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		return errmsg.ParseError
	}

	exist, err := db.IsUserExist(s.ctx, uid)
	if err != nil {
		return err
	}
	if !exist {
		return errmsg.UserNotExistError
	}

	toConn := userMap[to]
	switch toConn {
	case nil: // 离线
		{
			if err := db.CreateMessage(from, to, string(userinfoAppend(content, from))); err != nil {
				return errmsg.WebsockChatWriteError
			}
		}
	default: // 在线
		{
			err = toConn.conn.WriteMessage(websocket.TextMessage, content)
			if err != nil {
				return errmsg.WebsockChatWriteError
			}
		}
	}
	return nil
}

func (s ChatService) ReadOfflineMessage() error {
	uid := strconv.FormatInt(GetUidFormContext(s.c), 10)

	list, err := db.GetMessage(uid)
	if err != nil {
		return errmsg.WeBsockChatReadError
	}

	for _, item := range *list {
		ciphertext := userinfoAppend([]byte(item.Content), item.FromUserId)
		if err != nil {
			return errmsg.WebsockChatParseError
		}

		err = s.conn.WriteMessage(websocket.TextMessage, ciphertext)
		if err != nil {
			return errmsg.WeBsockChatReadError
		}
	}
	return nil
}

func userinfoAppend(rawText []byte, from string) []byte {
	return []byte(time.Now().Format(constants.TimeFormat) + ` [` + from + `]: ` + string(rawText))
}
