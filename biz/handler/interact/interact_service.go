// Code generated by hertz generator.

package interact

import (
	"context"
	"tiktok/biz/pack"
	"tiktok/biz/service"
	"tiktok/pkg/errmsg"

	"github.com/cloudwego/hertz/pkg/app"
	"tiktok/biz/model/interact"
)

// Like .
// @router /like/action [POST]
func Like(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interact.LikeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp := new(interact.LikeResponse)

	err = service.NewInteractService(ctx, c).Like(&req)

	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp.Base = pack.BuildBaseResp(errmsg.NoError)

	pack.SendResponse(c, resp)
}

// LikeList .
// @router /like/list [GET]
func LikeList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interact.LikeListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp := new(interact.LikeListResponse)

	interactResp, count, err := service.NewInteractService(ctx, c).LikeList(&req)

	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp.Base = pack.BuildBaseResp(errmsg.NoError)
	resp.Data = pack.LikeList(interactResp, count)

	pack.SendResponse(c, resp)
}

// Comment .
// @router /comment/publish [POST]
func Comment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interact.CommentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp := new(interact.CommentResponse)

	err = service.NewInteractService(ctx, c).Comment(&req)

	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp.Base = pack.BuildBaseResp(errmsg.NoError)

	pack.SendResponse(c, resp)
}

// CommentList .
// @router /comment/list [GET]
func CommentList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interact.CommentListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp := new(interact.CommentListResponse)

	interactResp, count, err := service.NewInteractService(ctx, c).CommentList(&req)

	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp.Base = pack.BuildBaseResp(errmsg.NoError)
	resp.Data = pack.CommentList(interactResp, count)

	pack.SendResponse(c, resp)
}

// DeleteComment .
// @router /comment/delete [DELETE]
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interact.DeleteCommentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp := new(interact.DeleteCommentResponse)

	err = service.NewInteractService(ctx, c).DeleteComment(&req)

	if err != nil {
		pack.BuildFailResponse(c, err)
		return
	}

	resp.Base = pack.BuildBaseResp(errmsg.NoError)

	pack.SendResponse(c, resp)
}
