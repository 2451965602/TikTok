package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

var VideocountKey = "VideoRank"

func CreatVideoId(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	var videoResp redis.Z

	videoResp = redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := RDB.ZAdd(ctx, VideocountKey, videoResp).Result()

	if err != nil {
		return err
	}

	return nil

}
