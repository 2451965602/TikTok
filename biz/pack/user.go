package pack

import (
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/model/model"
)

func User(data *db.User) *model.User {
	create := strconv.FormatInt(data.CreatedAt.Unix(), 10)
	update := strconv.FormatInt(data.UpdatedAt.Unix(), 10)
	return &model.User{
		ID:        data.UserId,
		Username:  data.Username,
		Password:  &data.Password,
		AvatarURL: data.AvatarUrl,
		OptSecret: &data.OptSecret,
		CreatedAt: &create,
		UpdatedAt: &update,
	}
}

func UserInfo(data *db.UserInfo) *model.UserInfo {
	create := strconv.FormatInt(data.CreatedAt.Unix(), 10)
	update := strconv.FormatInt(data.UpdatedAt.Unix(), 10)
	return &model.UserInfo{
		ID:        data.UserId,
		Username:  data.Username,
		AvatarURL: data.AvatarUrl,
		CreatedAt: &create,
		UpdatedAt: &update,
	}
}

func UserList(data []*db.User, total int64) *model.UserList {
	resp := make([]*model.User, 0, len(data))

	for _, v := range data {
		resp = append(resp, User(v))
	}

	return &model.UserList{
		Items: resp,
		Total: total,
	}
}

func UserInfoList(data []*db.UserInfo, total int64) *model.UserInfoList {
	resp := make([]*model.UserInfo, 0, len(data))

	for _, v := range data {
		resp = append(resp, UserInfo(v))
	}

	return &model.UserInfoList{
		Items: resp,
		Total: total,
	}
}

func MFA(data *db.MFA) *model.MFA {

	return &model.MFA{
		Secret: data.Secret,
		Qrcode: data.Qrcode,
	}
}
