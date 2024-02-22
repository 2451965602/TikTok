package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/model/interact"
)

type InteractService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewInteractService(ctx context.Context, c *app.RequestContext) *InteractService {
	return &InteractService{ctx: ctx, c: c}
}

func (s *InteractService) Like(req *interact.LikeRequest) error {
	return db.LikeVideo(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.VideoID, req.ActionType)
}

func (s *InteractService) LikeList(req *interact.LikeListRequest) ([]*db.Video, int64, error) {

	return db.LikeList(s.ctx, req.UserID, req.PageNum, req.PageSize)
}

func (s *InteractService) Comment(req *interact.CommentRequest) error {
	return db.CreatComment(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.VideoID, req.Content)
}

func (s *InteractService) CommentList(req *interact.CommentListRequest) ([]*db.Comment, int64, error) {
	return db.CommentList(s.ctx, req.VideoID, req.PageNum, req.PageSize)
}

func (s *InteractService) DeleteComment(req *interact.DeleteCommentRequest) error {

	VideoID, err := strconv.ParseInt(req.VideoID, 10, 64)
	CommentID, err := strconv.ParseInt(req.CommentID, 10, 64)

	if err != nil {
		return err
	}

	return db.DeleteComment(s.ctx, GetUidFormContext(s.c), VideoID, CommentID)
}
