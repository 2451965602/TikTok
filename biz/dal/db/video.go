package db

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"work4/biz/model/video"
	"work4/pkg/constants"
	"work4/pkg/errmsg"
)

func Feed(ctx context.Context, req *video.FeedRequest) ([]*Video, int64, error) {

	var videoResp []*Video
	var err error
	var count int64

	if req.LatestTime != nil && *req.LatestTime != "" {
		toTime, err := strconv.ParseInt(*req.LatestTime, 10, 64)
		if err != nil {
			return nil, -1, errmsg.ParseError
		}

		err = DB.
			WithContext(ctx).
			Table(constants.VideoTable).
			Where("created_at > ? ", time.Unix(toTime, 0)).
			Limit(int(req.PageSize)).
			Offset(int((req.PageNum - 1) * req.PageSize)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, errmsg.DatabaseError
		}
	} else {
		err = DB.
			WithContext(ctx).
			Table(constants.VideoTable).
			Limit(int(req.PageSize)).
			Offset(int((req.PageNum - 1) * req.PageSize)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, errmsg.DatabaseError
		}
	}

	return videoResp, count, nil
}

func UploadVideo(ctx context.Context, userid, videourl, coverurl, title, description string) (int64, error) {

	videoInfo := &Video{
		UserId:      userid,
		VideoUrl:    videourl,
		CoverUrl:    coverurl,
		Title:       title,
		Description: description,
	}

	err := DB.
		WithContext(ctx).
		Table(constants.VideoTable).
		Where("user_id=?", userid).
		Create(&videoInfo).
		Error

	if err != nil {
		return -1, errmsg.DatabaseError
	}

	return videoInfo.VideoId, nil
}

func UploadList(ctx context.Context, pagenum, pagesize, userid int64) ([]*Video, int64, error) {

	var videoResp []*Video
	var err error
	var count int64

	exist, err := IsUserExist(ctx, userid)
	if err != nil {
		return nil, -1, err
	}
	if !exist {
		return nil, -1, errmsg.UserNotExistError
	}

	err = DB.
		WithContext(ctx).
		Table(constants.VideoTable).
		Where("user_id=?", userid).
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&videoResp).
		Error

	if err != nil {
		return nil, -1, errmsg.DatabaseError
	}

	return videoResp, count, nil
}

func Rank(ctx context.Context, rank []string) ([]*Video, error) {

	var videoResp []*Video
	var err error
	var count int64

	err = DB.
		WithContext(ctx).
		Table("video").
		Where("video_id IN (?)", rank).
		Count(&count).
		Find(&videoResp).
		Error

	if err != nil {
		return nil, errmsg.DatabaseError
	}

	return videoResp, nil
}

func Query(ctx context.Context, req *video.QueryRequest) ([]*Video, int64, error) {

	var videoResp []*Video
	var err error
	var count int64
	var userinfo *UserInfo

	if req.Username != nil && req.FromDate != nil && req.ToDate != nil {

		err := DB.
			WithContext(ctx).
			Table(constants.UserTable).
			Select("user_id,username,avatar_url,created_at,updated_at,deleted_at").
			Where("username = ?", req.Username).
			First(&userinfo).
			Error

		if err != nil {
			return nil, -1, errmsg.UserNotExistError
		}

		err = DB.
			WithContext(ctx).
			Table(constants.VideoTable).
			Where("id=?", userinfo.UserId).
			Where("created_at > ? and created_at < ?", time.Unix(*req.FromDate, 0), time.Unix(*req.ToDate, 0)).
			Where("title LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).Or("description LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).
			Limit(int(req.PageSize)).
			Offset(int((req.PageNum - 1) * req.PageSize)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, errmsg.DatabaseError
		}
	}

	if req.Username != nil && req.FromDate == nil && req.ToDate == nil {

		err := DB.
			WithContext(ctx).
			Table(constants.UserTable).
			Select("user_id,username,avatar_url,created_at,updated_at,deleted_at").
			Where("username = ?", req.Username).
			First(&userinfo).
			Error

		if err != nil {
			return nil, -1, errmsg.UserNotExistError
		}

		err = DB.
			WithContext(ctx).
			Table(constants.VideoTable).
			Where("id=?", userinfo.UserId).
			Where("title LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).Or("description LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).
			Limit(int(req.PageSize)).
			Offset(int((req.PageNum - 1) * req.PageSize)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, errmsg.DatabaseError
		}
	}

	if req.Username == nil && req.FromDate == nil && req.ToDate == nil {

		err = DB.
			WithContext(ctx).
			Table(constants.VideoTable).
			Where("title LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).Or("description LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords)).
			Limit(int(req.PageSize)).
			Offset(int((req.PageNum - 1) * req.PageSize)).
			Count(&count).
			Find(&videoResp).
			Error

		if err != nil {
			return nil, -1, errmsg.DatabaseError
		}
	}

	return videoResp, count, nil
}
