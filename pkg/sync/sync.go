package sync

import (
	"context"
	"log"
	"tiktok/biz/dal/db"
	"tiktok/biz/dal/redis"
	"time"
)

func StartSync() {
	time.Sleep(10 * time.Minute)
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		CountSync()
	}
}

func CountSync() {
	ctx := context.Background()
	videoIds, err := redis.GetAllVideoIds(ctx)
	if err != nil {
		log.Println("Error getting video ids:", err)

		return
	}
	for _, videoId := range videoIds {
		counts, err := redis.GetCounts(ctx, videoId)
		if err != nil {
			log.Println("Error getting counts for video id", videoId, ":", err)

			continue
		}
		err = db.UpdataVideoCounts(ctx, videoId, db.Counts(counts))
		if err != nil {
			log.Println("Error updating counts for video id", videoId, ":", err)

			continue
		}
	}
}
