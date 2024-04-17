package dal

import (
	"tiktok/biz/dal/db"
	"tiktok/biz/dal/redis"
)

func MysqlInit() {
	db.Init()
}

func RedisInit() {
	redis.Init()
}
