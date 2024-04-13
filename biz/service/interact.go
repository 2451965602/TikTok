package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/dal/redis"
	"work4/biz/model/interact"
	"work4/pkg/errmsg"
)

type InteractService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewInteractService(ctx context.Context, c *app.RequestContext) *InteractService {
	return &InteractService{ctx: ctx, c: c}
}

func (s *InteractService) Like(req *interact.LikeRequest) error {

	var (
		err error
	)

	if req.VideoID != nil && req.CommentID == nil {

		VideoID, err := strconv.ParseInt(*req.VideoID, 10, 64)
		if err != nil {
			return errmsg.ParseError
		}

		err = db.CreateLike(s.ctx, GetUidFormContext(s.c), VideoID, req.ActionType, "video")
		if err != nil {
			return err
		}

	} else if req.VideoID == nil && req.CommentID != nil {

		CommentID, err := strconv.ParseInt(*req.CommentID, 10, 64)
		if err != nil {
			return errmsg.ParseError
		}

		err = db.CreateLike(s.ctx, GetUidFormContext(s.c), CommentID, req.ActionType, "comment")
		if err != nil {
			return err
		}

	} else {
		return errmsg.DuplicationError.WithMessage("No liking of videos and comments at the same time")
	}

	if req.VideoID != nil {
		if req.ActionType == "1" {
			err = redis.AddLikeCount(s.ctx, *req.VideoID)
			if err != nil {
				return err
			}
		} else if req.ActionType == "2" {
			err = redis.ReduceLikeCount(s.ctx, *req.VideoID)
			if err != nil {
				return err
			}
		} else {
			return errmsg.ServiceError.WithMessage("parametric error")
		}

		err = redis.UpdateIdInRank(s.ctx, *req.VideoID)
		if err != nil {
			return err
		}
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
		count, err := redis.GetCounts(s.ctx, strconv.FormatInt(v.VideoId, 10))
		if err != nil {
			return nil, -1, err
		}
		v.VisitCount = count.VisitCount
		v.LikeCount = count.LikeCount
		v.CommentCount = count.CommentCount
		resp = append(resp, v)
	}
	return resp, num, nil
}

func (s *InteractService) Comment(req *interact.CommentRequest) error {

	if req.CommentID == nil && req.VideoID != nil {

		err := db.CreatComment(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), *req.VideoID, req.Content, "video")
		if err != nil {
			return err
		}

		err = redis.UpdateIdInRank(s.ctx, *req.VideoID)
		if err != nil {
			return err
		}

		err = redis.AddCommentCount(s.ctx, *req.VideoID)
		if err != nil {
			return err
		}

	} else if req.CommentID != nil && req.VideoID == nil {

		err := db.CreatComment(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), *req.CommentID, req.Content, "comment")
		if err != nil {
			return err
		}

	} else {
		return errmsg.DuplicationError.WithMessage("No liking of videos and comments at the same time")
	}

	return nil
}

func (s *InteractService) CommentList(req *interact.CommentListRequest) ([]*db.Comment, int64, error) {
	return db.CommentList(s.ctx, req.VideoID, req.PageNum, req.PageSize)
}

func (s *InteractService) DeleteComment(req *interact.DeleteCommentRequest) error {

	commentid, err := strconv.ParseInt(req.CommentID, 10, 64)
	if err != nil {
		return errmsg.ParseError
	}

	videoid, err := db.DeleteComment(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), commentid)
	if err != nil {
		return errmsg.ParseError
	}

	err = redis.ReduceCommentCount(s.ctx, videoid)
	if err != nil {
		return err
	}

	return nil
}
