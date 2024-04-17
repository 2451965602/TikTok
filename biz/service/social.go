package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"tiktok/biz/dal/db"
	"tiktok/biz/model/social"
	"tiktok/pkg/errmsg"
)

type SocialService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewSocialService(ctx context.Context, c *app.RequestContext) *SocialService {
	return &SocialService{ctx: ctx, c: c}
}

func (s *SocialService) Star(req *social.StarRequest) error {

	var (
		userid   int64
		touserid int64
	)

	userid = GetUidFormContext(s.c)
	touserid, err := strconv.ParseInt(req.ToUserID, 10, 64)
	if err != nil {
		return errmsg.ParseError
	}

	if userid > touserid {
		return db.StarUser(s.ctx, userid, touserid, req.ActionType, 1)
	} else {
		return db.StarUser(s.ctx, touserid, userid, req.ActionType, 0)
	}

}

func (s *SocialService) StarList(req *social.StarListRequest) ([]*db.UserInfoDetail, int64, error) {

	var (
		resp []*db.UserInfoDetail
	)
	data, count, err := db.StarUserList(s.ctx, GetUidFormContext(s.c), req.PageNum, req.PageSize)
	if err != nil {
		return nil, -1, err
	}

	for _, v := range data {
		UserInfo, err := db.GetInfo(s.ctx, strconv.FormatInt(v.UserId, 10))
		if err != nil {
			return nil, -1, err
		}

		resp = append(resp, UserInfo)
	}

	return resp, count, nil
}

func (s *SocialService) FanList(req *social.FanListRequest) ([]*db.UserInfoDetail, int64, error) {

	var (
		resp []*db.UserInfoDetail
	)
	data, count, err := db.FanUserList(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.PageNum, req.PageSize)
	if err != nil {
		return nil, -1, err
	}

	for _, v := range data {
		UserInfo, err := db.GetInfo(s.ctx, strconv.FormatInt(v.UserId, 10))
		if err != nil {
			return nil, -1, err
		}

		resp = append(resp, UserInfo)
	}

	return resp, count, nil
}

func (s *SocialService) FriendList(req *social.FriendListRequest) ([]*db.UserInfoDetail, int64, error) {

	var (
		resp []*db.UserInfoDetail
	)
	data, count, err := db.FriendUser(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.PageNum, req.PageSize)
	if err != nil {
		return nil, -1, err
	}

	for _, v := range data {
		UserInfo, err := db.GetInfo(s.ctx, strconv.FormatInt(v.UserId, 10))
		if err != nil {
			return nil, -1, err
		}

		resp = append(resp, UserInfo)
	}

	return resp, count, nil
}
