package pack

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"work4/biz/model/model"
	"work4/pkg/errmsg"
)

type Base struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type Response struct {
	Base Base `json:"base"`
}

func SendResponse(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, data)
}

func SendFailResponse(c *app.RequestContext, data *model.BaseResp) {
	c.JSON(consts.StatusOK, utils.H{
		"base": data,
	})
}

func BuildBaseResp(err errmsg.ErrorMessage) *model.BaseResp {
	return &model.BaseResp{
		Code: err.ErrorCode,
		Msg:  err.ErrorMsg,
	}
}

func BuildFailResponse(c *app.RequestContext, err error) {
	if err == nil {
		SendFailResponse(c, BuildBaseResp(errmsg.NoError))
		return
	}

	e := errmsg.ErrorMessage{}
	if errors.As(err, &e) {
		SendFailResponse(c, BuildBaseResp(e))
		return
	}

	e = errmsg.ServiceError.WithMessage(err.Error())
	SendFailResponse(c, BuildBaseResp(e))
}
