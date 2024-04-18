package db

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"image/png"
	"tiktok/pkg/constants"
	"tiktok/pkg/errmsg"
	"tiktok/pkg/util"
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

func GetUserInfo(ctx context.Context, userid string) (*User, error) {
	var user *User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("user_id=?", userid).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserInfoByName(ctx context.Context, username string) (*User, error) {
	var user *User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func OptSecret(user *User) (*MFA, error) {
	var MFAResp = &MFA{}
	var buf bytes.Buffer

	if user.OptSecret == "" {
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "westonline",
			AccountName: user.Username,
		})
		if err != nil {
			return nil, errmsg.MfaGenareteError
		}

		user.OptSecret = key.String()
	}

	key, err := otp.NewKeyFromURL(user.OptSecret)
	if err != nil {
		return nil, errmsg.MfaGenareteError.WithMessage(err.Error())
	}

	img, err := key.Image(200, 200)
	if err != nil {
		return nil, errmsg.MfaGenareteError.WithMessage(err.Error())
	}

	err = png.Encode(&buf, img)
	if err != nil {
		return nil, errmsg.MfaGenareteError.WithMessage(err.Error())
	}

	qrcode := base64.StdEncoding.EncodeToString(buf.Bytes())

	secret, err := util.ExtractSecretFromTOTPURL(user.OptSecret)
	if err != nil {
		return nil, err
	}

	MFAResp.Secret = secret
	MFAResp.Qrcode = qrcode

	return MFAResp, nil
}
