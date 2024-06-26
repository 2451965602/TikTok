package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"mime/multipart"
	"strconv"
	"tiktok/biz/dal/db"
	"tiktok/biz/dal/redis"
	"tiktok/biz/model/video"
	"tiktok/pkg/errmsg"
	"tiktok/pkg/oss"
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

	resp, num, err := db.Feed(s.ctx, req)
	if err != nil {
		return nil, -1, err
	}

	return resp, num, err
}

func (s *VideoService) UploadVideo(videodata *multipart.FileHeader, coverdata *multipart.FileHeader, req *video.UploadRequest) error {

	userid := strconv.FormatInt(GetUidFormContext(s.c), 10)

	err := oss.IsVideo(videodata)
	if err != nil {
		return err
	}

	err = oss.IsImage(coverdata)
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
	if err != nil {
		return err
	}

	err = redis.AddIdToRank(s.ctx, strconv.FormatInt(videoId, 10))
	if err != nil {
		return err
	}

	return nil
}

func (s *VideoService) UploadList(req *video.UploadListRequest) ([]*db.Video, int64, error) {

	var resp []*db.Video

	userid, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		return nil, -1, errmsg.ParseError
	}

	resp, num, err := db.UploadList(s.ctx, req.PageNum, req.PageSize, userid)
	if err != nil {
		return nil, -1, err
	}

	return resp, num, err
}

func (s *VideoService) Rank(req *video.RankRequest) ([]*db.Video, error) {

	resp, err := redis.RankList(s.ctx)
	if err != nil {
		return nil, err
	}

	if resp == nil {

		rank, err := redis.IdRankList(s.ctx)
		if err != nil {
			return nil, err
		}

		resp, err := db.Rank(s.ctx, rank)
		if err != nil {
			return nil, err
		}

		err = redis.AddToRank(s.ctx, resp)
		if err != nil {
			return nil, err
		}
	}

	startIndex := (req.PageNum - 1) * req.PageSize
	endIndex := startIndex + req.PageSize

	if startIndex >= int64(len(resp)) {
		return []*db.Video{}, nil
	}

	if endIndex > int64(len(resp)) {
		endIndex = int64(len(resp))
	}

	return resp[startIndex:endIndex], nil

}

func (s *VideoService) Query(req *video.QueryRequest) ([]*db.Video, int64, error) {

	if req.Keywords == nil && req.Username == nil {

		return nil, -1, errmsg.IllegalParamError
	}

	var resp []*db.Video

	resp, num, err := db.Query(s.ctx, req)
	if err != nil {
		return nil, -1, err
	}

	return resp, num, err
}
