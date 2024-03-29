package pack

import (
	"strconv"
	"work4/biz/dal/db"
	"work4/biz/model/model"
)

func Comment(data *db.Comment) *model.Comment {
	return &model.Comment{
		ID:        strconv.FormatInt(data.CommentId, 10),
		UserID:    data.UserId,
		VideoID:   data.VideoId,
		Content:   data.Content,
		CreatedAt: strconv.FormatInt(data.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(data.UpdatedAt.Unix(), 10),
	}
}

func Like(data *db.Like) *model.Like {
	return &model.Like{
		UserID:  data.UserId,
		VideoID: data.VideoId,
	}
}

func CommentList(data []*db.Comment, total int64) *model.CommentList {
	resp := make([]*model.Comment, 0, len(data))

	for _, v := range data {
		resp = append(resp, Comment(v))
	}

	return &model.CommentList{
		Items: resp,
		Total: total,
	}
}

func LikeList(data []*db.Video, total int64) *model.LikeList {
	resp := make([]*model.Video, 0, len(data))

	for _, v := range data {
		resp = append(resp, Video(v))
	}

	return &model.LikeList{
		Items: resp,
		Total: total,
	}
}
