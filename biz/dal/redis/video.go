package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
	"work4/biz/dal/db"
)

var VideoIdKey = "VideoId"
var VideoKey = "Video"

func AddIdToRank(ctx context.Context, videoid string) error {

	if redisDBVideoId == nil {
		return errors.New("redisDBVideoId object is nil")
	}

	var videoResp redis.Z

	videoResp = redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := redisDBVideoId.ZAdd(ctx, VideoIdKey, videoResp).Result()
	if err != nil {
		return err
	}

	return nil
}

func UpdateIdInRank(ctx context.Context, videoid string) error {

	if redisDBVideoId == nil {
		return errors.New("redisDBVideoId object is nil")
	}

	score, err := redisDBVideoId.ZScore(ctx, VideoIdKey, videoid).Result()
	if err != nil && err != redis.Nil {
		return err
	} else if err == redis.Nil {
		err := AddIdToRank(ctx, videoid)
		if err != nil {
			return err
		}
	}

	videoResp := redis.Z{
		Score:  score + 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, VideoIdKey, videoResp).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	return nil

}

func IdRankList(ctx context.Context) ([]string, error) {

	if redisDBVideoId == nil {
		return nil, errors.New("redisDBVideoId object is nil")
	}

	rank, err := redisDBVideoId.ZRevRange(ctx, VideoIdKey, 0, 99).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	return rank, nil

}

func AddToRank(ctx context.Context, videolist []*db.Video) error {
	if redisDBVideo == nil {
		return errors.New("redisDBVideo 对象为 nil")
	}

	pipe := redisDBVideo.Pipeline()
	for _, video := range videolist {

		videoJSON, err := json.Marshal(video)
		if err != nil {
			return err
		}

		pipe.ZAdd(ctx, VideoKey, redis.Z{
			Score:  float64(video.VisitCount),
			Member: videoJSON,
		})
	}

	pipe.Expire(ctx, VideoKey, 10*time.Minute)

	_, err := pipe.Exec(ctx)
	return err
}

func RankList(ctx context.Context) ([]*db.Video, error) {
	if redisDBVideo == nil {
		return nil, errors.New("redisDBVideo 对象为 nil")
	}

	var exist = false

	memberJSONStrings, err := redisDBVideo.ZRevRange(ctx, VideoKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	videos := make([]*db.Video, len(memberJSONStrings))
	for i, memberJSON := range memberJSONStrings {
		var video db.Video
		err := json.Unmarshal([]byte(memberJSON), &video)
		if err != nil {
			return nil, err
		}
		videos[i] = &video
		exist = true
	}
	if !exist {
		return nil, nil
	}
	return videos, nil
}
