package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"work4/pkg/constants"
	"work4/pkg/errmsg"
)

func QuerySocialStatus(ctx context.Context, userID int64) (Social, error) {

	var (
		resp Social
		err  error
	)
	err = DB.
		WithContext(ctx).
		Table(constants.SocialTable).
		Where("user_id = ?", userID).
		First(&resp).
		Error

	return resp, err
}

func UpdateSocialStatus(ctx context.Context, userID int64, status int) error {
	return DB.
		WithContext(ctx).
		Table(constants.SocialTable).
		Where("user_id = ?", userID).
		Update("Status", status).
		Error
}

func DeleteSocialRecord(ctx context.Context, userID int64, social *Social) error {
	return DB.WithContext(ctx).
		Table(constants.SocialTable).
		Where("user_id = ?", userID).
		Delete(social).
		Error
}

/*
UserId > ToUserId

actiontype
0-关注
1-取关

status
0-互相关注
1- smallid（touserid被关注
2- bigid（userid）被关注

to
1-smallid（touserid）被(取消)关注
0-bigid（userid）被(取消)关注
*/

// CreateOrUpdateFollowRecord 创建或更新关注记录
func CreateOrUpdateFollowRecord(ctx context.Context, bigid, smallid, to int64) error {

	var (
		status int64
	)
	_, err := QuerySocialStatus(ctx, bigid)

	// 如果关注记录不存在
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新的关注记录
			status = 1
			if to == 0 {
				status = 2
			}

			err = DB.
				WithContext(ctx).
				Table(constants.SocialTable).
				Create(&Social{
					UserId:   bigid,
					ToUserId: smallid,
					Status:   status,
				}).
				Error
			if err != nil {
				return errmsg.DatabaseError
			}

			return nil
		}

		return errmsg.DatabaseError
	}

	// 更新已有的关注记录
	err = UpdateSocialStatus(ctx, bigid, 0)
	if err != nil {
		return errmsg.DatabaseError
	}

	return nil
}

// CancleStarUser 取消关注
func CancleStarUser(ctx context.Context, bigid, to int64) error {
	social, err := QuerySocialStatus(ctx, bigid)
	if err != nil {
		return errmsg.DatabaseError
	}

	// 如果是互相关注状态
	if social.Status == 0 {
		if to == 1 {
			// 更新状态为取消关注
			err = UpdateSocialStatus(ctx, bigid, 2)
		} else {
			// 更新状态为取消关注
			err = UpdateSocialStatus(ctx, bigid, 1)
		}

		if err != nil {
			return errmsg.DatabaseError
		}

	} else {
		// 删除关注记录
		err = DeleteSocialRecord(ctx, bigid, &social)
		if err != nil {
			return errmsg.DatabaseError
		}

	}

	return nil
}

// StarUser 关注或取消关注用户
func StarUser(ctx context.Context, bigid, smallid, actiontype, to int64) (err error) {

	// 关注操作
	if actiontype == 0 {
		err = CreateOrUpdateFollowRecord(ctx, bigid, smallid, to)
		if err != nil {
			return err
		}

	} else if actiontype == 1 { // 取消关注操作
		err = CancleStarUser(ctx, bigid, to)
		if err != nil {
			return err
		}
	} else {
		return errmsg.IllegalParamError
	}

	return nil
}

func StarUserList(ctx context.Context, userid, pagenum, pagesize int64) ([]*UserInfo, int64, error) {

	var StarResp []*UserInfo
	var err error
	var count int64

	err = DB.
		WithContext(ctx).
		Table(constants.SocialTable).
		Where(&Social{UserId: userid, Status: 1}).
		Or(&Social{ToUserId: userid, Status: 2}).
		Or(&Social{UserId: userid, Status: 0}).
		Or(&Social{ToUserId: userid, Status: 0}).
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&StarResp).
		Error

	if err != nil {
		return nil, -1, errmsg.DatabaseError
	}

	return StarResp, count, nil
}

func FanUserList(ctx context.Context, userid string, pagenum, pagesize int64) ([]*UserInfo, int64, error) {

	var StarResp []*UserInfo
	var err error
	var count int64

	err = DB.
		WithContext(ctx).
		Table(constants.SocialTable).
		Where("to_user_id=?", userid).
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&StarResp).
		Error

	if err != nil {
		return nil, -1, errmsg.DatabaseError
	}

	return StarResp, count, nil
}

func FriendUser(ctx context.Context, userid string, pagenum, pagesize int64) ([]*UserInfo, int64, error) {

	var StarResp []*UserInfo
	var userId []*string
	var err error
	var count int64

	err = DB.
		WithContext(ctx).
		Table(constants.SocialTable).
		Where("user_id = ?", userid).Or("to_user_id = ?", userid).
		Where("status = ?", 0).
		Limit(int(pagesize)).
		Offset(int((pagenum - 1) * pagesize)).
		Count(&count).
		Find(&userId).
		Error

	if err != nil {
		return nil, -1, errmsg.DatabaseError
	}

	return StarResp, count, nil
}
