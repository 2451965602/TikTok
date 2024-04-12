package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
	"work4/biz/dal/db"
	"work4/pkg/errmsg"
)

var VideoIdKey = "VideoId"
var VideoKey = "Video"

func AddIdToRank(ctx context.Context, videoid string) error {

	var videoResp redis.Z

	videoResp = redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := redisDBVideoId.ZAdd(ctx, VideoIdKey, videoResp).Result()
	if err != nil {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil
}

func UpdateIdInRank(ctx context.Context, videoid string) error {

	score, err := redisDBVideoId.ZScore(ctx, VideoIdKey, videoid).Result()
	if err != nil && err != redis.Nil {
		return errmsg.RedisError.WithMessage(err.Error())
	} else if err == redis.Nil {
		err := AddIdToRank(ctx, videoid)
		if err != nil {
			return errmsg.RedisError.WithMessage(err.Error())
		}
	}

	videoResp := redis.Z{
		Score:  score + 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, VideoIdKey, videoResp).Result()
	if err != nil && err != redis.Nil {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil

}

func IdRankList(ctx context.Context) ([]string, error) {

	rank, err := redisDBVideoId.ZRevRange(ctx, VideoIdKey, 0, 99).Result()
	if err != nil && err != redis.Nil {
		return nil, errmsg.RedisError.WithMessage(err.Error())
	}

	return rank, nil

}

func AddToRank(ctx context.Context, videolist []*db.Video) error {

	pipe := redisDBVideo.Pipeline()
	for _, video := range videolist {

		videoJSON, err := json.Marshal(video)
		if err != nil {
			return errmsg.RedisError.WithMessage(err.Error())
		}

		pipe.ZAdd(ctx, VideoKey, redis.Z{
			Score:  float64(video.VisitCount),
			Member: videoJSON,
		})
	}

	pipe.Expire(ctx, VideoKey, 10*time.Minute)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errmsg.RedisError.WithMessage(err.Error())
	}
	return nil
}

func RankList(ctx context.Context) ([]*db.Video, error) {

	memberJSONStrings, err := redisDBVideo.ZRevRange(ctx, VideoKey, 0, -1).Result()
	if err != nil {
		return nil, errmsg.RedisError.WithMessage(err.Error())
	}

	videos := make([]*db.Video, len(memberJSONStrings))
	for i, memberJSON := range memberJSONStrings {
		var video db.Video
		err := json.Unmarshal([]byte(memberJSON), &video)
		if err != nil {
			return nil, errmsg.RedisError.WithMessage(err.Error())
		}
		videos[i] = &video
	}
	return videos, nil
}
