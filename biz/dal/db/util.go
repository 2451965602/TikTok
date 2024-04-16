package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"work4/pkg/constants"
	"work4/pkg/errmsg"
)

func IsVideoExist(ctx context.Context, videoid int64) (bool, error) {
	var video Video

	err := DB.
		WithContext(ctx).
		Table(constants.VideoTable).
		Where("video_id=?", videoid).
		First(&video).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, errmsg.DatabaseError.WithMessage(err.Error())
	}

	return true, nil
}

func IsUserExist(ctx context.Context, userid int64) (bool, error) {
	var user User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("user_id=?", userid).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, errmsg.DatabaseError.WithMessage(err.Error())
	}

	return true, nil
}

func IsUserNameExist(ctx context.Context, username string) (bool, error) {
	var user User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("username=?", username).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, errmsg.DatabaseError.WithMessage(err.Error())
	}

	return true, nil
}

func IsCommentExist(ctx context.Context, commentid int64) (bool, error) {
	var comment Comment

	err := DB.
		WithContext(ctx).
		Table(constants.CommentTable).
		Where("comment_id=?", commentid).
		First(&comment).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, errmsg.DatabaseError.WithMessage(err.Error())
	}

	return true, nil
}
