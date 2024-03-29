package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"work4/biz/model/video"
	"work4/bootstrap/env"
)

func Feed(ctx context.Context, req *video.FeedRequest) ([]*Video, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

	var videoResp []*Video
	var err error
	var count int64

	if req.LatestTime != nil && *req.LatestTime != "" {
		toTime, err := strconv.ParseInt(*req.LatestTime, 10, 64)
		if err != nil {
			return nil, -1, err
		}

		err = DB.
			WithContext(ctx).
			Table(env.VideoTable).
			Where("created_at > ? ", time.Unix(toTime, 0)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, err
		}
	} else {
		err = DB.
			WithContext(ctx).
			Table(env.VideoTable).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, err
		}
	}

	return videoResp, count, nil
}

func UploadVideo(ctx context.Context, userid, videourl, coverurl, title, description string) (int64, error) {

	if DB == nil {
		return -1, errors.New("DB object is nil")
	}

	var err error
	var videoInfo *Video

	videoInfo = &Video{
		UserId:      userid,
		VideoUrl:    videourl,
		CoverUrl:    coverurl,
		Title:       title,
		Description: description,
	}

	err = DB.
		WithContext(ctx).
		Table(env.VideoTable).
		Where("user_id=?", userid).
		Create(&videoInfo).
		Error

	if err != nil {
		return -1, err
	}

	return videoInfo.VideoId, nil
}

func UploadList(ctx context.Context, pagenum, pagesize int64, userid string) ([]*Video, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

	var videoResp []*Video
	var err error
	var count int64

	err = DB.
		WithContext(ctx).
		Table(env.VideoTable).
		Where("user_id=?", userid).
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&videoResp).
		Error

	if err != nil {
		return nil, -1, err
	}

	return videoResp, count, nil
}

func Rank(ctx context.Context, pagenum, pagesize int64, rank []string) ([]*Video, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

	var videoResp []*Video
	var err error
	var count int64

	err = DB.
		WithContext(ctx).
		Table("video").
		Where("video_id IN (?)", rank).
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&videoResp).
		Error

	if err != nil {
		return nil, -1, err
	}

	return videoResp, count, nil
}

func Query(ctx context.Context, req *video.QueryRequest) ([]*Video, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

	var videoResp []*Video
	var err error
	var count int64
	var userinfo *UserInfo

	if req.Username != nil && req.FromDate != nil && req.ToDate != nil {

		err := DB.
			WithContext(ctx).
			Table(env.UserTable).
			Select("user_id,username,avatar_url,created_at,updated_at,deleted_at").
			Where("username = ?", req.Username).
			First(&userinfo).
			Error

		if err != nil {
			return nil, -1, errors.New("用户不存在")
		}

		err = DB.
			WithContext(ctx).
			Table(env.VideoTable).
			Where("id=?", userinfo.UserId).
			Where("created_at > ? and created_at < ?", time.Unix(*req.FromDate, 0), time.Unix(*req.ToDate, 0)).
			Where("title LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).Or("description LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).
			Limit(int(req.PageSize)).
			Offset(int((req.PageNum - 1) * req.PageSize)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, err
		}
	}

	if req.Username != nil && req.FromDate == nil && req.ToDate == nil {

		err := DB.
			WithContext(ctx).
			Table(env.UserTable).
			Select("user_id,username,avatar_url,created_at,updated_at,deleted_at").
			Where("username = ?", req.Username).
			First(&userinfo).
			Error

		if err != nil {
			return nil, -1, errors.New("用户不存在")
		}

		err = DB.
			WithContext(ctx).
			Table(env.VideoTable).
			Where("id=?", userinfo.UserId).
			Where("title LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).Or("description LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).
			Limit(int(req.PageSize)).
			Offset(int((req.PageNum - 1) * req.PageSize)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, err
		}
	}

	if req.Username == nil && req.FromDate == nil && req.ToDate == nil {

		err = DB.
			WithContext(ctx).
			Table(env.VideoTable).
			Where("title LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).Or("description LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).
			Limit(int(req.PageSize)).
			Offset(int((req.PageNum - 1) * req.PageSize)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, err
		}
	}

	return videoResp, count, nil
}
