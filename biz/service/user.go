package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"mime/multipart"
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/model/user"
	"work4/pkg/upload"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

func (s *UserService) Register(req *user.RegisterRequest) (*db.User, error) {
	return db.CreateUser(s.ctx, req.Username, req.Password)
}

func (s *UserService) Login(req *user.LoginRequest) (*db.UserInfoDetail, error) {

	return db.LoginCheck(s.ctx, req)
}

func (s *UserService) GetInfo(req *user.InfoRequest) (*db.UserInfoDetail, error) {
	return db.GetInfo(s.ctx, req.UserID)
}

func (s *UserService) UploadAvatar(avatar *multipart.FileHeader) (*db.User, error) {

	userid := strconv.FormatInt(GetUidFormContext(s.c), 10)

	err := upload.IsImage(avatar)
	if err != nil {
		return nil, err
	}

	imgurl, err := UploadAndGetUrl(avatar, userid, "avatar")

	if err != nil {
		return nil, err
	}

	return db.UploadAvatar(s.ctx, userid, imgurl)
}

func (s *UserService) MFAGet() (*db.MFA, error) {
	return db.MFAGet(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10))
}

func (s *UserService) MFABind(req *user.MFABindRequest) error {
	return db.MFABind(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.Secret, req.Code)
}
