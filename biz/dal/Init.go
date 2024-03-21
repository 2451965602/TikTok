package dal

import (
	"work4/biz/dal/db"
	"work4/biz/dal/redis"
)

func Init() {
	db.Init()
	redis.Init()
}
