package db

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserId    int64
	Username  string
	Password  string
	AvatarUrl string
	OptSecret string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserInfo struct {
	UserId    int64
	Username  string
	AvatarUrl string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Video struct {
	VideoId      int64
	UserId       string
	VideoUrl     string
	CoverUrl     string
	Title        string
	Description  string
	VisitCount   int64
	LikeCount    int64
	CommentCount int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Comment struct {
	CommentId int64
	UserId    string
	VideoId   string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MFA struct {
	Secret string
	Qrcode string
}

type Social struct {
	UserId   string
	ToUserId string
	Status   int64
}

type Like struct {
	UserId  string
	VideoId string
}
