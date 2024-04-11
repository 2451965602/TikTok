package db

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"work4/biz/model/user"
	"work4/pkg/constants"
	"work4/pkg/util"

	"work4/pkg/crypt"
)

func OptSecret(user *User) (*MFA, error) {
	var MFAResp = &MFA{}
	var buf bytes.Buffer

	if user.OptSecret == "" {
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "westonline",
			AccountName: user.Username,
		})
		if err != nil {
			return nil, err
		}

		user.OptSecret = key.String()

	}

	key, err := otp.NewKeyFromURL(user.OptSecret)
	if err != nil {
		return nil, err
	}

	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}

	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
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

func CreateUser(ctx context.Context, username, password string) (*User, error) {
	if DB == nil {
		return nil, errors.New("DB object is nil")
	}

	var userResp *User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("username = ?", username).
		First(&userResp).
		Error

	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	pwd, err := crypt.PasswordHash(password)

	if err != nil {
		return nil, err
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
		return nil, err
	}

	userResp.Password = password

	return userResp, nil
}

func LoginCheck(ctx context.Context, req *user.LoginRequest) (*UserInfoDetail, error) {

	if DB == nil {
		return nil, errors.New("DB object is nil")
	}

	var userreq *User
	var userResp *UserInfoDetail

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("username = ?", req.Username).
		First(&userreq).
		Error

	if err != nil {
		return nil, errors.New("账号或密码错误")
	}

	if !crypt.PasswordVerify(req.Password, userreq.Password) {
		return nil, errors.New("账号或密码错误")
	}

	if userreq.OptSecret != "" {
		if req.Code == nil {
			// OTP 已启用，但未提供 OTP 代码
			return nil, errors.New("未提供 OTP 代码")
		}
		value := *req.Code
		if !totp.Validate(value, userreq.OptSecret) {
			// OTP 验证失败
			return nil, errors.New("OTP 一次性密码错误")
		}
	}

	userResp = &UserInfoDetail{
		UserId:    userreq.UserId,
		Username:  userreq.Username,
		AvatarUrl: userreq.AvatarUrl,
		CreatedAt: userreq.CreatedAt,
	}

	return userResp, nil
}

func GetInfo(ctx context.Context, id string) (*UserInfoDetail, error) {

	if DB == nil {
		return nil, errors.New("DB object is nil")
	}

	var userResp *UserInfoDetail

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Select("user_id,username,avatar_url,created_at,updated_at,deleted_at").
		Where("user_id = ?", id).
		First(&userResp).
		Error

	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return userResp, nil
}

func UploadAvatar(ctx context.Context, id, url string) (*User, error) {

	if DB == nil {
		return nil, errors.New("DB object is nil")
	}

	var userResp *User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("user_id = ?", id).
		Update("avatar_url", url).
		First(&userResp).
		Error

	if err != nil {
		return nil, err
	}

	return userResp, nil
}

func MFAGet(ctx context.Context, id string) (*MFA, error) {

	if DB == nil {
		return nil, errors.New("DB object is nil")
	}

	var userreq *User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("user_id = ?", id).
		First(&userreq).
		Error

	if err != nil {
		return nil, err
	}

	MFAResp, err := OptSecret(userreq)

	if err != nil {
		return nil, err
	}

	return MFAResp, err
}

func MFABind(ctx context.Context, id, secret, code string) error {

	if DB == nil {
		return errors.New("DB object is nil")
	}

	if totp.Validate(code, secret) {
		err := DB.
			WithContext(ctx).
			Table(constants.UserTable).
			Where("user_id = ?", id).
			Update("opt_secret", secret).
			Error

		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("OTP 一次性密码错误")
}
