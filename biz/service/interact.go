package service

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/dal/redis"
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

	var err error

	err = db.LikeVideo(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.VideoID, req.ActionType)
	if err != nil {
		return err
	}

	if req.ActionType == "1" {
		err = redis.AddLikeCount(s.ctx, req.VideoID)
		if err != nil {
			return err
		}
	} else if req.ActionType == "2" {
		err = redis.ReduceLikeCount(s.ctx, req.VideoID)
		if err != nil {
			return err
		}
	} else {
		return errors.New("非法操作")
	}

	err = redis.UpdateRank(s.ctx, req.VideoID)
	if err != nil {
		return err
	}

	return nil
}

func (s *InteractService) LikeList(req *interact.LikeListRequest) ([]*db.Video, int64, error) {

	var resp []*db.Video

	temp, num, err := db.LikeList(s.ctx, req.UserID, req.PageNum, req.PageSize)

	if err != nil {
		return nil, -1, err
	}
	for _, v := range temp {
		v.VisitCount, err = redis.GetVisitCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.LikeCount, err = redis.GetLikeCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.CommentCount, err = redis.GetCommentCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		resp = append(resp, v)
	}
	return resp, num, nil
}

func (s *InteractService) Comment(req *interact.CommentRequest) error {

	err := db.CreatComment(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.VideoID, req.Content)
	if err != nil {
		return err
	}

	err = redis.UpdateRank(s.ctx, req.VideoID)
	if err != nil {
		return err
	}

	err = redis.AddCommentCount(s.ctx, req.VideoID)
	if err != nil {
		return err
	}

	return nil
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

	err = redis.UpdateRank(s.ctx, req.VideoID)
	if err != nil {
		return err
	}

	err = redis.ReduceCommentCount(s.ctx, req.VideoID)
	if err != nil {
		return err
	}

	return db.DeleteComment(s.ctx, GetUidFormContext(s.c), VideoID, CommentID)
}
