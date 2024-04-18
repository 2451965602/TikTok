package main

import (
	"bytes"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"os"
	"testing"
	"tiktok/biz/middleware/jwt"
	"tiktok/pkg/cfg"
)

var token string

func hInit() *server.Hertz {
	err := cfg.Init()
	if err != nil {
		hlog.Info("TEST 配置读取失败")
		os.Exit(1)
		return nil
	}
	jwt.Init()
	h := server.Default()
	register(h)
	return h
}

func TestUserRegister(t *testing.T) {
	h := hInit()
	req := `{
		"username":"test",
		"password":"123456789"
	}
	`
	resp := ut.PerformRequest(h.Engine, "POST", "/user/register", &ut.Body{Body: bytes.NewBufferString(req), Len: len(req)},
		ut.Header{Key: `Content-Type`, Value: `application/json`},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestUserLoginNoMfa(t *testing.T) {
	h := hInit()
	req := `{
		"username":"test",
		"password":"123456789"
	}
	`
	resp := ut.PerformRequest(h.Engine, "POST", "/user/login", &ut.Body{Body: bytes.NewBufferString(req), Len: len(req)},
		ut.Header{Key: `Content-Type`, Value: `application/json`},
	)
	token = resp.Header().Get("Access-Token")

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestUserLoginWithMfa(t *testing.T) {
	h := hInit()
	req := `{
		"username":"testmfa",
		"password":"123456",
		"code":"513144"
	}
	`
	resp := ut.PerformRequest(h.Engine, "POST", "/user/login", &ut.Body{Body: bytes.NewBufferString(req), Len: len(req)},
		ut.Header{Key: `Content-Type`, Value: `application/json`},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestUserInfo(t *testing.T) {
	h := hInit()
	const (
		userid = `10001`
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/user/info?user_id="+userid, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestGetMfa(t *testing.T) {
	h := hInit()
	resp := ut.PerformRequest(h.Engine, "GET", "/auth/mfa/qrcode", nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestGetFeed(t *testing.T) {
	h := hInit()
	const (
		pagesize = "10"
		pagenum  = "1"
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/video/feed?"+"&page_size="+pagesize+"&page_num="+pagenum, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestVideolist(t *testing.T) {
	h := hInit()
	const (
		userid   = "10001"
		pagesize = "10001"
		pagenum  = "1"
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/video/list?user_id="+userid+"&page_num="+pagenum+"&page_size="+pagesize, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestVideoPopular(t *testing.T) {
	h := hInit()
	const (
		pagesize = "10001"
		pagenum  = "1"
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/video/popular?page_size="+pagesize+"&page_num="+pagenum, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestVideoSearch(t *testing.T) {
	h := hInit()
	req := `{
		"keywords":"apple",
		"page_num":1,
		"page_size":10
	}
	`
	resp := ut.PerformRequest(h.Engine, "POST", "/video/search", &ut.Body{Body: bytes.NewBufferString(req), Len: len(req)},
		ut.Header{Key: "Access-Token", Value: token},
		ut.Header{Key: `Content-Type`, Value: `application/json`},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestLike(t *testing.T) {
	h := hInit()
	req := `{
		"video_id":"10001",
		"action_type":"1"
	}
	`
	resp := ut.PerformRequest(h.Engine, "POST", "/like/action", &ut.Body{Body: bytes.NewBufferString(req), Len: len(req)},
		ut.Header{Key: "Access-Token", Value: token},
		ut.Header{Key: `Content-Type`, Value: `application/json`},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestLikeList(t *testing.T) {
	h := hInit()
	const (
		userid   = "10001"
		pagesize = "10001"
		pagenum  = "1"
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/like/list?user_id="+userid+"&page_size="+pagesize+"&page_num="+pagenum, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestComment(t *testing.T) {
	h := hInit()
	req := `{
		"video_id":"10001",
		"content":"test"
	}
	`
	resp := ut.PerformRequest(h.Engine, "POST", "/comment/publish", &ut.Body{Body: bytes.NewBufferString(req), Len: len(req)},
		ut.Header{Key: "Access-Token", Value: token},
		ut.Header{Key: `Content-Type`, Value: `application/json`},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestCommentList(t *testing.T) {
	h := hInit()
	const (
		videoid  = "10001"
		pagesize = "10001"
		pagenum  = "1"
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/comment/list?video_id="+videoid+"&page_size="+pagesize+"&page_num="+pagenum, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestStar(t *testing.T) {
	h := hInit()
	req := `{
		"to_user_id":"10001",
		"action_type":0
	}
	`
	resp := ut.PerformRequest(h.Engine, "POST", "/relation/action", &ut.Body{Body: bytes.NewBufferString(req), Len: len(req)},
		ut.Header{Key: "Access-Token", Value: token},
		ut.Header{Key: `Content-Type`, Value: `application/json`},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestStarList(t *testing.T) {
	h := hInit()
	const (
		userid   = "10001"
		pagesize = "10001"
		pagenum  = "1"
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/following/list?user_id="+userid+"&page_size="+pagesize+"&page_num="+pagenum, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestFanList(t *testing.T) {
	h := hInit()
	const (
		userid   = "10001"
		pagesize = "10001"
		pagenum  = "1"
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/follower/list?user_id="+userid+"&page_size="+pagesize+"&page_num="+pagenum, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}

func TestFriendList(t *testing.T) {
	h := hInit()
	const (
		userid   = "10001"
		pagesize = "10"
		pagenum  = "1"
	)
	resp := ut.PerformRequest(h.Engine, "GET", "/friends/list?user_id="+userid+"&page_size="+pagesize+"&page_num="+pagenum, nil,
		ut.Header{Key: "Access-Token", Value: token},
	)

	assert.DeepEqual(t, 200, resp.Result().StatusCode())
}
