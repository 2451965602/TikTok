package db

import (
	"context"
	"errors"
	"strconv"
	"work4/biz/dal/redis"
	"work4/pkg/env"
)

func LikeVideo(ctx context.Context, userid, videoid, actiontype string) error {

	if DB == nil {
		return errors.New("DB object is nil")
	}

	var LikeResp *Like

	LikeResp = &Like{
		UserId:  userid,
		VideoId: videoid,
	}

	if actiontype == "1" {
		err := DB.
			WithContext(ctx).
			Table(env.LikeTable).
			Create(&LikeResp).
			Error

		if err != nil {
			return err
		}

		err = redis.AddLikeCount(ctx, videoid)
		if err != nil {
			return err
		}
	} else {
		err := DB.
			WithContext(ctx).
			Table(env.LikeTable).
			Delete(&LikeResp).
			Error

		if err != nil {
			return err
		}

		err = redis.ReduceLikeCount(ctx, videoid)
		if err != nil {
			return err
		}
	}

	err := redis.UpdateRank(ctx, videoid)
	if err != nil {
		return err
	}

	return nil
}

func LikeList(ctx context.Context, userid string, pagenum, pagesize int64) ([]*Video, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

	var LikeResp []*Video
	var videoid []*int64
	var err error
	var count int64

	// 获取同一个user_id下所有的video_id
	err = DB.
		WithContext(ctx).
		Table(env.LikeTable).
		Where("user_id = ?", userid).
		Select("video_id").
		Find(&videoid).
		Error

	if err != nil {
		return nil, -1, err
	}

	// 查询video表中与likeIDs匹配的视频信息
	err = DB.
		Table(env.VideoTable).
		Where("`video_id` IN (?)", videoid).
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&LikeResp).
		Error

	if err != nil {
		return nil, -1, err
	}

	return LikeResp, count, nil
}

func CreatComment(ctx context.Context, userid, videoid, content string) error {

	if DB == nil {
		return errors.New("DB object is nil")
	}

	var CommentResp *Comment

	CommentResp = &Comment{
		UserId:  userid,
		VideoId: videoid,
		Content: content,
	}

	err := DB.
		WithContext(ctx).
		Table(env.CommentTable).
		Create(&CommentResp).
		Error

	if err != nil {
		return err
	}

	err = redis.UpdateRank(ctx, videoid)
	if err != nil {
		return err
	}

	err = redis.AddCommentCount(ctx, videoid)
	if err != nil {
		return err
	}

	return nil
}

func CommentList(ctx context.Context, videoid string, pagenum, pagesize int64) ([]*Comment, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

	var CommentResp []*Comment
	var err error
	var count int64

	err = DB.
		WithContext(ctx).
		Table(env.CommentTable).
		Where("video_id=?", videoid).
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&CommentResp).
		Error

	if err != nil {
		return nil, -1, err
	}

	return CommentResp, count, nil
}

func DeleteComment(ctx context.Context, userid, videoid, commentid int64) error {

	if DB == nil {
		return errors.New("DB object is nil")
	}

	err := DB.
		WithContext(ctx).
		Table(env.CommentTable).
		Where("user_id=?", userid).
		Where("video_id=?", videoid).
		Where("comment_id=?", commentid).
		Delete(&Comment{
			CommentId: commentid,
			UserId:    strconv.FormatInt(userid, 10),
			VideoId:   strconv.FormatInt(videoid, 10),
		}).
		Error

	if err != nil {
		return err
	}

	err = redis.UpdateRank(ctx, strconv.FormatInt(videoid, 10))
	if err != nil {
		return err
	}

	err = redis.ReduceCommentCount(ctx, strconv.FormatInt(videoid, 10))
	if err != nil {
		return err
	}

	return nil
}
