package dal

import (
	"time"
	"work4/biz/dal/db"
	"work4/biz/dal/redis"
)

func SyncInit(stopChan chan struct{}) {
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 执行同步操作
			db.Sync()
		case <-stopChan:
			// 收到停止信号，退出循环
			return
		}
	}

}

func Init() {
	db.Init()
	redis.Init()
}
