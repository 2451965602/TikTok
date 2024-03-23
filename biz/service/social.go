package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/model/social"
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
		temp     int
		touserid int64
	)

	userid = GetUidFormContext(s.c)
	temp, _ = strconv.Atoi(req.ToUserID)
	touserid = int64(temp)

	if userid > touserid {
		return db.StarUser(s.ctx, strconv.FormatInt(userid, 10), strconv.FormatInt(touserid, 10), req.ActionType, 1)
	} else {
		return db.StarUser(s.ctx, strconv.FormatInt(touserid, 10), strconv.FormatInt(userid, 10), req.ActionType, 0)
	}

}

func (s *SocialService) StarList(req *social.StarListRequest) ([]*db.UserInfoDetail, int64, error) {

	var (
		resp []*db.UserInfoDetail
	)
	data, count, err := db.StarUserList(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.PageNum, req.PageSize)
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
