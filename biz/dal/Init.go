package dal

import (
	"tiktok/biz/dal/db"
	"tiktok/biz/dal/redis"
)

func MysqlInit() error {
	err := db.Init()
	if err != nil {
		return err
	}

	return nil
}

func RedisInit() error {
	err := redis.Init()
	if err != nil {
		return err
	}

	return nil
}
