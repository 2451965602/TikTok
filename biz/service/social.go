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

func (s *SocialService) StarList(req *social.StarListRequest) ([]*db.UserInfo, int64, error) {
	return db.StarUserList(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.PageNum, req.PageSize)
}

func (s *SocialService) FanList(req *social.FanListRequest) ([]*db.UserInfo, int64, error) {
	return db.FanUserList(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.PageNum, req.PageSize)
}

func (s *SocialService) FriendList(req *social.FriendListRequest) ([]*db.UserInfo, int64, error) {
	return db.FriendUser(s.ctx, strconv.FormatInt(GetUidFormContext(s.c), 10), req.PageNum, req.PageSize)
}
