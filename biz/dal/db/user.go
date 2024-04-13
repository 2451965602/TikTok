package db

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"work4/biz/model/user"
	"work4/pkg/constants"
	"work4/pkg/errmsg"
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

	var userreq *User
	var userResp *UserInfoDetail

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("username = ?", req.Username).
		First(&userreq).
		Error

	if err != nil {
		return nil, errmsg.AuthError.WithMessage("Incorrect account number or password")
	}

	if !crypt.PasswordVerify(req.Password, userreq.Password) {
		return nil, errmsg.AuthError.WithMessage("Incorrect account number or password")
	}

	if userreq.OptSecret != "" {
		if req.Code == nil {
			return nil, errmsg.MfaOptCodeError.WithMessage("OTP code not provided")
		}
		value := *req.Code
		if !totp.Validate(value, userreq.OptSecret) {
			return nil, errmsg.MfaOptCodeError
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

	err = DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("user_id = ?", id).
		First(&userResp).
		Error

	if err != nil {
		return nil, errmsg.DatabaseError
	}

	return userResp, nil
}

func MFAGet(ctx context.Context, id string) (*MFA, error) {

	var userreq *User

	err := DB.
		WithContext(ctx).
		Table(constants.UserTable).
		Where("user_id = ?", id).
		First(&userreq).
		Error

	if err != nil {
		return nil, errmsg.DatabaseError
	}

	MFAResp, err := OptSecret(userreq)
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
			Error

		if err != nil {
			return errmsg.DatabaseError
		}

		return nil
	}

	return errmsg.MfaOptCodeError
}
