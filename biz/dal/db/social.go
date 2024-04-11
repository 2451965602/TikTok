package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"work4/pkg/constants"
)

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

// createOrUpdateFollowRecord 创建或更新关注记录
func createOrUpdateFollowRecord(ctx context.Context, bigid, smallid string, to int64) error {
	var (
		social Social
		status int64
	)
	err := DB.
		WithContext(ctx).
		Table(constants.SocialTable).
		Where("user_id = ?", bigid).
		First(&social).
		Error

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
			return err
		}
		return err
	}

	// 更新已有的关注记录
	err = DB.
		WithContext(ctx).
		Table(constants.SocialTable).
		Where("user_id = ?", bigid).
		Update("Status", 0).
		Error
	return err
}

// StarUser 关注或取消关注用户
func StarUser(ctx context.Context, bigid, smallid string, actiontype, to int64) (err error) {
	// 检查数据库对象是否为空
	if DB == nil {
		return errors.New("DB 对象为空")
	}
	// 关注操作
	if actiontype == 0 {
		err = createOrUpdateFollowRecord(ctx, bigid, smallid, to)
	} else { // 取消关注操作
		var social Social
		err = DB.
			WithContext(ctx).
			Table(constants.SocialTable).
			Where("user_id = ?", bigid).
			First(&social).
			Error

		if err != nil {
			return err
		}

		// 如果是互相关注状态
		if social.Status == 0 {
			if to == 1 {
				// 更新状态为取消关注
				err = DB.
					WithContext(ctx).
					Table(constants.SocialTable).
					Where("user_id = ?", bigid).
					Update("Status", 2).
					Error
			} else {
				// 删除关注记录
				err = DB.
					WithContext(ctx).
					Table(constants.SocialTable).
					Where("user_id = ?", bigid).
					Update("Status", 1).
					Error
			}
		} else {
			err = DB.
				WithContext(ctx).
				Table(constants.SocialTable).
				Where("user_id = ?", bigid).
				Delete(&social).
				Error
		}
	}

	return err
}

func StarUserList(ctx context.Context, userid string, pagenum, pagesize int64) ([]*UserInfo, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

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
		return nil, -1, err
	}

	return StarResp, count, nil
}

func FanUserList(ctx context.Context, userid string, pagenum, pagesize int64) ([]*UserInfo, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

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
		return nil, -1, err
	}
	return StarResp, count, nil
}

func FriendUser(ctx context.Context, userid string, pagenum, pagesize int64) ([]*UserInfo, int64, error) {

	if DB == nil {
		return nil, -1, errors.New("DB object is nil")
	}

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
		return nil, -1, err
	}

	return StarResp, count, nil
}
