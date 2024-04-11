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
		ID:        strconv.FormatInt(data.UserId, 10),
		Username:  data.Username,
		Password:  &data.Password,
		AvatarURL: data.AvatarUrl,
		OptSecret: &data.OptSecret,
		CreatedAt: &create,
		UpdatedAt: &update,
	}
}

func UserInfoDetail(data *db.UserInfoDetail) *model.UserInfo {
	createat := data.CreatedAt.Format("2006-01-02 15:04:05")

	updateat := ""
	if !data.UpdatedAt.IsZero() {
		updateat = data.UpdatedAt.Format("2006-01-02 15:04:05")
	} else {
		updateat = "1970-01-01 08:00:00"
	}

	deleteat := ""
	if !data.DeletedAt.Time.IsZero() {
		deleteat = data.DeletedAt.Time.Format("2006-01-02 15:04:05")
	} else {
		deleteat = "1970-01-01 08:00:00"
	}

	return &model.UserInfo{
		ID:        strconv.FormatInt(data.UserId, 10),
		Username:  data.Username,
		AvatarURL: data.AvatarUrl,
		CreatedAt: &createat,
		UpdatedAt: &updateat,
		DeletedAt: &deleteat,
	}
}

func UserInfo(data *db.UserInfoDetail) *model.UserInfo {
	createat := data.CreatedAt.Format("2006-01-02 15:04:05")

	updateat := ""
	if !data.UpdatedAt.IsZero() {
		updateat = data.UpdatedAt.Format("2006-01-02 15:04:05")
	} else {
		updateat = "1970-01-01 08:00:00"
	}

	return &model.UserInfo{
		ID:        strconv.FormatInt(data.UserId, 10),
		Username:  data.Username,
		AvatarURL: data.AvatarUrl,
		CreatedAt: &createat,
		UpdatedAt: &updateat,
	}
}

func UserInfoList(data []*db.UserInfoDetail, total int64) *model.UserInfoList {
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
