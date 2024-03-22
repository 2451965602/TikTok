package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"mime/multipart"
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/dal/redis"
	"work4/biz/model/video"
	"work4/pkg/upload"
)

type VideoService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewVideoService(ctx context.Context, c *app.RequestContext) *VideoService {
	return &VideoService{ctx: ctx, c: c}
}

func (s *VideoService) Feed(req *video.FeedRequest) ([]*db.Video, int64, error) {

	var resp []*db.Video

	temp, num, err := db.Feed(s.ctx, req)

	if err != nil {
		return nil, -1, err
	}
	for _, v := range temp {
		v.VisitCount, err = redis.GetVisitCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.LikeCount, err = redis.GetLikeCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.CommentCount, err = redis.GetCommentCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		resp = append(resp, v)
	}

	return resp, num, err
}

func (s *VideoService) UploadVideo(videodata *multipart.FileHeader, coverdata *multipart.FileHeader, req *video.UploadRequest) error {

	userid := strconv.FormatInt(GetUidFormContext(s.c), 10)

	err := upload.IsVideo(videodata)
	if err != nil {
		return err
	}

	err = upload.IsImage(coverdata)
	if err != nil {
		return err
	}

	videoUrl, err := UploadAndGetUrl(videodata, userid, "video")
	if err != nil {
		return err
	}

	coverUrl, err := UploadAndGetUrl(coverdata, userid, "cover")
	if err != nil {
		return err
	}

	videoId, err := db.UploadVideo(s.ctx, userid, videoUrl, coverUrl, req.Title, req.Description)

	err = redis.AddRank(s.ctx, strconv.FormatInt(videoId, 10))
	if err != nil {
		return err
	}

	return nil
}

func (s *VideoService) UploadList(req *video.UploadListRequest) ([]*db.Video, int64, error) {

	var resp []*db.Video

	temp, num, err := db.UploadList(s.ctx, req.PageNum, req.PageSize, req.UserID)

	if err != nil {
		return nil, -1, err
	}
	for _, v := range temp {
		v.VisitCount, err = redis.GetVisitCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.LikeCount, err = redis.GetLikeCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.CommentCount, err = redis.GetCommentCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		resp = append(resp, v)
	}

	return resp, num, err
}

func (s *VideoService) Rank(req *video.RankRequest) ([]*db.Video, int64, error) {

	var resp []*db.Video

	rank, err := redis.RankList(s.ctx)
	if err != nil {
		return nil, -1, err
	}

	temp, num, err := db.Rank(s.ctx, req.PageNum, req.PageSize, rank)
	if err != nil {
		return nil, -1, err
	}

	for _, v := range temp {
		v.VisitCount, err = redis.GetVisitCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.LikeCount, err = redis.GetLikeCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.CommentCount, err = redis.GetCommentCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		resp = append(resp, v)
	}

	return resp, num, err

}

func (s *VideoService) Query(req *video.QueryRequest) ([]*db.Video, int64, error) {

	var resp []*db.Video

	temp, num, err := db.Query(s.ctx, req)

	if err != nil {
		return nil, -1, err
	}
	for _, v := range temp {
		v.VisitCount, err = redis.GetVisitCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.LikeCount, err = redis.GetLikeCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		v.CommentCount, err = redis.GetCommentCount(s.ctx, strconv.FormatInt(v.VideoId, 10))
		resp = append(resp, v)
	}

	return resp, num, err
}
