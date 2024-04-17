package db

import (
	"context"
	"github.com/pquerna/otp/totp"
	"tiktok/biz/model/user"
	"tiktok/pkg/constants"
	"tiktok/pkg/crypt"
	"tiktok/pkg/errmsg"
)

func CreateUser(ctx context.Context, username, password string) (*User, error) {

	var userResp *User

	exist, err := IsUserNameExist(ctx, username)

	if exist {
		return nil, errmsg.UserExistError
	}
	if err != nil {
		return nil, err
	}

	pwd, err := crypt.PasswordHash(password)

	if err != nil {
		return nil, errmsg.CryptEncodeError
	}

	userResp = &User{
		Username: username,
		Password: pwd,
	}

	err = DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Create(&userResp).
		Error

	if err != nil {
		return nil, errmsg.DatabaseError
	}

	userResp.Password = password

	return userResp, nil
}

func LoginCheck(ctx context.Context, req *user.LoginRequest) (*UserInfoDetail, error) {

	var userReq *User
	var userResp *UserInfoDetail

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("username=?", req.Username).
		First(&userReq).
		Error

	if err != nil {
		return nil, errmsg.AuthError.WithMessage("Incorrect account number or password")
	}

	if !crypt.PasswordVerify(req.Password, userReq.Password) {
		return nil, errmsg.AuthError.WithMessage("Incorrect account number or password")
	}

	if userReq.OptSecret != "" {
		if req.Code == nil {
			return nil, errmsg.MfaOptCodeError.WithMessage("OTP code not provided")
		}
		value := *req.Code
		if !totp.Validate(value, userReq.OptSecret) {
			return nil, errmsg.MfaOptCodeError
		}
	}

	userResp = &UserInfoDetail{
		UserId:    userReq.UserId,
		Username:  userReq.Username,
		AvatarUrl: userReq.AvatarUrl,
		CreatedAt: userReq.CreatedAt,
	}

	return userResp, nil
}

func GetInfo(ctx context.Context, id string) (*UserInfoDetail, error) {

	var userResp *UserInfoDetail

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Select("user_id,username,avatar_url,created_at,updated_at,deleted_at").
		Where("user_id = ?", id).
		First(&userResp).
		Error

	if err != nil {
		return nil, errmsg.UserNotExistError
	}

	return userResp, nil
}

func UploadAvatar(ctx context.Context, id, url string) (*User, error) {

	var userResp *User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("user_id = ?", id).
		Update("avatar_url", url).
		Error

	if err != nil {
		return nil, errmsg.DatabaseError
	}

	userResp, err = GetUserInfo(ctx, id)
	if err != nil {
		return nil, errmsg.DatabaseError
	}

	return userResp, nil
}

func MFAGet(ctx context.Context, id string) (*MFA, error) {

	userReq, err := GetUserInfo(ctx, id)

	if err != nil {
		return nil, errmsg.DatabaseError
	}

	MFAResp, err := OptSecret(userReq)
	if err != nil {
		return nil, err
	}

	return MFAResp, nil
}

func MFABind(ctx context.Context, id, secret, code string) error {

	if totp.Validate(code, secret) {
		err := DB.
			WithContext(ctx).
			Table(constants.UserTable).
			Where("user_id = ?", id).
			Update("opt_secret", secret).
			Update("mfa_status", 1).
			Error

		if err != nil {
			return errmsg.DatabaseError
		}

		return nil
	}

	return errmsg.MfaOptCodeError
}

func MFAStatus(ctx context.Context, id, code, ActionType string) error {

	userInfo, err := GetUserInfo(ctx, id)
	if err != nil {
		return errmsg.DatabaseError
	}

	if ActionType == userInfo.MfaStatus {
		return errmsg.DuplicationError
	}

	if totp.Validate(code, userInfo.OptSecret) {
		err := DB.
			WithContext(ctx).
			Table(constants.UserTable).
			Where("user_id = ?", id).
			Update("mfa_status", ActionType).
			Error

		if err != nil {
			return errmsg.DatabaseError
		}

		return nil
	}

	return errmsg.MfaOptCodeError
}
