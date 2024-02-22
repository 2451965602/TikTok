package redis

import (
	"github.com/redis/go-redis/v9"
	"work4/pkg/env"
)

var RDB *redis.Client

func Init() {
	var err error

	opt, err := redis.ParseURL(env.RedisDSN)
	if err != nil {
		panic(err)
	}

	RDB = redis.NewClient(opt)

}
