package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tiktok/biz/dal/db"
	"tiktok/pkg/errmsg"
	"time"
)

var VideoIdKey = "VideoId"
var VideoKey = "Video"

func AddIdToRank(ctx context.Context, videoid string) error {

	videoResp := redis.Z{
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
	if err != nil && !errors.Is(err, redis.Nil) {
		return errmsg.RedisError.WithMessage(err.Error())
	} else if errors.Is(err, redis.Nil) {
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
	if err != nil && !errors.Is(err, redis.Nil) {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil

}

func GetAllVideoIds(ctx context.Context) ([]string, error) {
	videoIds, err := redisDBVideoId.ZRange(ctx, VideoIdKey, 0, -1).Result()
	if err != nil {
		return nil, errmsg.RedisError.WithMessage(err.Error())
	}

	return videoIds, nil
}

func IdRankList(ctx context.Context) ([]string, error) {

	rank, err := redisDBVideoId.ZRevRange(ctx, VideoIdKey, 0, 99).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
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
