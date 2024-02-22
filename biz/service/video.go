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
	return db.Feed(s.ctx, req)
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

	videourl, err := UploadAndGetUrl(videodata, userid, "video")
	if err != nil {
		return err
	}

	coverurl, err := UploadAndGetUrl(coverdata, userid, "cover")
	if err != nil {
		return err
	}

	return db.UploadVideo(s.ctx, userid, videourl, coverurl, req.Title, req.Description)
}

func (s *VideoService) UploadList(req *video.UploadListRequest) ([]*db.Video, int64, error) {
	return db.UploadList(s.ctx, req.PageNum, req.PageSize, req.UserID)
}

func (s *VideoService) Rank(req *video.RankRequest) ([]*db.Video, int64, error) {

	rank, err := redis.RankList(s.ctx)
	if err != nil {
		return nil, -1, err
	}
	return db.Rank(s.ctx, req.PageNum, req.PageSize, rank)
}

func (s *VideoService) Query(req *video.QueryRequest) ([]*db.Video, int64, error) {
	return db.Query(s.ctx, req)
}
