package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

var videoKey = "VideoRank"

func AddRank(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	var videoResp redis.Z

	videoResp = redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := RDB.ZAdd(ctx, videoKey, videoResp).Result()

	if err != nil {
		return err
	}

	err = CreatVideoId(ctx, videoid)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRank(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	score, err := RDB.ZScore(ctx, videoKey, videoid).Result()

	if err != nil {
		return err
	}

	videoResp := redis.Z{
		Score:  score + 1,
		Member: videoid,
	}

	_, err = RDB.ZAdd(ctx, videoKey, videoResp).Result()

	if err != nil {
		return err
	}

	err = CreatVideoId(ctx, videoid)
	if err != nil {
		return err
	}

	return nil

}

func GetVisitCount(ctx context.Context, videoid string) (int64, error) {

	if RDB == nil {
		return -1, errors.New("RDB object is nil")
	}

	count, err := RDB.ZScore(ctx, videoKey, videoid).Result()

	if err != nil {
		return -1, err
	}

	return int64(count), nil

}

func RankList(ctx context.Context) ([]string, error) {

	if RDB == nil {
		return nil, errors.New("RDB object is nil")
	}

	rank, err := RDB.ZRevRange(ctx, videoKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return rank, nil

}
