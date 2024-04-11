package dal

import (
	"work4/biz/dal/db"
	"work4/biz/dal/redis"
)

func MysqlInit() {
	db.Init()
}

func RedisInit() {
	redis.Init()
}
