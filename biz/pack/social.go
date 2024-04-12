package pack

import (
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/model/model"
)

func Social(data *db.Social) *model.Social {

	return &model.Social{
		UserID:   strconv.FormatInt(data.UserId, 10),
		ToUserID: strconv.FormatInt(data.ToUserId, 10),
		Status:   data.Status,
	}
}
