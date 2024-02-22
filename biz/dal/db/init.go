package db

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
	"work4/biz/dal/redis"
	"work4/pkg/env"
)

var DB *gorm.DB

func Init() {
	var err error

	DB, err = gorm.Open(mysql.Open(env.MySQLDSN), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}

func Sync() {

	var ctx context.Context

	rank, err := redis.RDB.ZRevRange(ctx, redis.VideocountKey, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	for _, videoID := range rank {
		CommentCount, err := redis.GetCommentCount(ctx, videoID)
		LikeCount, err := redis.GetLikeCount(ctx, videoID)
		VisitCount, err := redis.GetVisitCount(ctx, videoID)
		if err != nil {
			panic(err)
		}

		err = UpdateLikeCount(ctx, strconv.FormatInt(LikeCount, 10), videoID)
		if err != nil {
			panic(err)
		}

		err = UpdateCommentCount(ctx, strconv.FormatInt(CommentCount, 10), videoID)
		if err != nil {
			panic(err)
		}

		err = UpdateVisitCount(ctx, strconv.FormatInt(VisitCount, 10), videoID)
		if err != nil {
			panic(err)
		}
	}
}
