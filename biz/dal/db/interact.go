package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"work4/bootstrap/env"
)

func CreateLike(ctx context.Context, userid, id int64, actiontype, sort string) (err error) {

	if DB == nil {
		return errors.New("DB object is nil")
	}

	var LikeResp *Like

	if sort == "video" {

		LikeResp = &Like{
			UserId:  userid,
			VideoId: id,
			RootId:  0,
		}

	} else {
		LikeResp = &Like{
			UserId:  userid,
			VideoId: 0,
			RootId:  id,
		}

	}

	if actiontype == "1" {

		err = DB.
			WithContext(ctx).
			Table(env.LikeTable).
			Where("root_id=?", id).
			Or("video_id=?", id).
			Or("user_id=?", userid).
			First(&LikeResp).
			Error

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("当前已点赞，请勿重复操作")
		}

		err = DB.
			WithContext(ctx).
			Table(env.LikeTable).
			Create(&LikeResp).
			Error

		if err != nil {
			return err
		}

	} else {

		err = DB.
			WithContext(ctx).
			Table(env.LikeTable).
			Where("root_id=?", id).
			Or("video_id=?", id).
			Or("user_id=?", userid).
			First(&LikeResp).
			Error

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("当前未点赞，请勿重复操作")
		}

		err = DB.
			WithContext(ctx).
			Table(env.LikeTable).
			Where("root_id=?", id).
			Or("video_id=?", id).
			Or("user_id=?", userid).
			Delete(&LikeResp).
			Error

		if err != nil {
			return err
		}
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

func CreatComment(ctx context.Context, userid, id, content, sort string) error {

	if DB == nil {
		return errors.New("DB object is nil")
	}

	var CommentResp *Comment

	if sort == "video" {

		CommentResp = &Comment{
			UserId:  userid,
			VideoId: id,
			Content: content,
		}

	} else {

		CommentResp = &Comment{
			UserId:  userid,
			RootId:  id,
			Content: content,
		}

	}

	err := DB.
		WithContext(ctx).
		Table(env.CommentTable).
		Create(&CommentResp).
		Error

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

func DeleteComment(ctx context.Context, userid string, commentid int64) (videoid int64, err error) {

	if DB == nil {
		return -1, errors.New("DB object is nil")
	}

	var commentInfo Comment

	err = DB.
		WithContext(ctx).
		Table(env.CommentTable).
		Where("comment_id=?", commentid).
		Select("video_id").
		First(&commentInfo).
		Error

	if err != nil {
		return -1, err
	}

	err = DB.
		WithContext(ctx).
		Table(env.CommentTable).
		Where("comment_id = ?", commentid).
		Delete(&Comment{
			CommentId: commentid,
			UserId:    userid,
		}).
		Error

	if err != nil {
		return -1, err
	}

	return -1, nil
}
