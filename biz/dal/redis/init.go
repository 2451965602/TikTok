package redis

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
	"work4/pkg/constants"
)

var (
	redisDBVideo   *redis.Client
	redisDBVideoId *redis.Client
)

func Init() {

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

	hlog.Info("Redis连接成功")
}
