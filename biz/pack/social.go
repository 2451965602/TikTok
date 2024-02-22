package pack

import (
	"work4/biz/dal/db"
	"work4/biz/model/model"
)

func Social(data *db.Social) *model.Social {

	return &model.Social{
		UserID:   data.UserId,
		ToUserID: data.ToUserId,
		Status:   data.Status,
	}
}
