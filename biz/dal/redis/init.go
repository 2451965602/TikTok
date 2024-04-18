package redis

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
	"tiktok/pkg/constants"
	"tiktok/pkg/errmsg"
)

var (
	redisDBVideo   *redis.Client
	redisDBVideoId *redis.Client
)

func Init() error {
	ctx := context.Background()

	redisDBVideoId = redis.NewClient(&redis.Options{
		Addr:     constants.RedisHost + ":" + constants.RedisPort,
		Username: constants.RedisUserName,
		Password: constants.RedisPassWord,
		DB:       0,
	})

	redisDBVideo = redis.NewClient(&redis.Options{
		Addr:     constants.RedisHost + ":" + constants.RedisPort,
		Username: constants.RedisUserName,
		Password: constants.RedisPassWord,
		DB:       1,
	})

	if _, err := redisDBVideo.Ping(ctx).Result(); err != nil {
		return errmsg.RedisInitError
	}

	if _, err := redisDBVideoId.Ping(ctx).Result(); err != nil {
		return errmsg.RedisInitError
	}

	hlog.Info("Redis连接成功")

	return nil
}
