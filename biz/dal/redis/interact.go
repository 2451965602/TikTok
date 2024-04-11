package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

var LikeToIdKey = "VideoLike"
var CommentToIdKey = "CommentToId"

type Counts struct {
	VisitCount   int64
	LikeCount    int64
	CommentCount int64
}

func CreatLikeCount(ctx context.Context, videoid string) error {

	if redisDBVideoId == nil {
		return errors.New("redisDBVideoId object is nil")
	}

	videoResp := redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := redisDBVideoId.ZAdd(ctx, LikeToIdKey, videoResp).Result()

	if err != nil {
		return err
	}

	return nil

}

func AddLikeCount(ctx context.Context, videoid string) error {

	if redisDBVideoId == nil {
		return errors.New("redisDBVideoId object is nil")
	}

	score, err := redisDBVideoId.ZScore(ctx, LikeToIdKey, videoid).Result()

	if err != nil {

		if err == redis.Nil {
			err = CreatLikeCount(ctx, videoid)
			if err != nil && err != redis.Nil {
				return err
			}
		} else {
			return err
		}

	}

	videoResp := redis.Z{
		Score:  score + 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, LikeToIdKey, videoResp).Result()
	if err != nil {
		return err
	}

	return nil

}

func ReduceLikeCount(ctx context.Context, videoid string) error {

	if redisDBVideoId == nil {
		return errors.New("redisDBVideoId object is nil")
	}

	score, err := redisDBVideoId.ZScore(ctx, LikeToIdKey, videoid).Result()

	if err != nil && err != redis.Nil {
		return err
	}

	videoResp := redis.Z{
		Score:  score - 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, LikeToIdKey, videoResp).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	return nil

}

func CreatCommentCount(ctx context.Context, videoid string) error {

	if redisDBVideoId == nil {
		return errors.New("redisDBVideoId object is nil")
	}

	videoResp := redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := redisDBVideoId.ZAdd(ctx, CommentToIdKey, videoResp).Result()

	if err != nil && err != redis.Nil {
		return err
	}

	return nil

}

func AddCommentCount(ctx context.Context, videoid string) error {

	if redisDBVideoId == nil {
		return errors.New("redisDBVideoId object is nil")
	}

	score, err := redisDBVideoId.ZScore(ctx, CommentToIdKey, videoid).Result()

	if err != nil {

		if err == redis.Nil {
			err = CreatCommentCount(ctx, videoid)
			if err != nil {
				return err
			}
		} else {
			return err
		}

	}

	videoResp := redis.Z{
		Score:  score + 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, CommentToIdKey, videoResp).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	return nil

}

func ReduceCommentCount(ctx context.Context, videoid string) error {

	if redisDBVideoId == nil {
		return errors.New("redisDBVideoId object is nil")
	}

	score, err := redisDBVideoId.ZScore(ctx, CommentToIdKey, videoid).Result()

	if err != nil && err != redis.Nil {
		return err
	}

	videoResp := redis.Z{
		Score:  score - 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, CommentToIdKey, videoResp).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	return nil

}

func GetCounts(ctx context.Context, videoId string) (Counts, error) {
	if redisDBVideoId == nil {
		return Counts{}, errors.New("redisDBVideoId object is nil")
	}

	pipe := redisDBVideoId.Pipeline()

	visitCmd := pipe.ZScore(ctx, VideoIdKey, videoId)
	likeCmd := pipe.ZScore(ctx, LikeToIdKey, videoId)
	commentCmd := pipe.ZScore(ctx, CommentToIdKey, videoId)

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return Counts{}, err
	}

	visitCount, visitErr := visitCmd.Result()
	if visitErr != nil && visitErr != redis.Nil {
		return Counts{}, visitErr
	}

	likeCount, likeErr := likeCmd.Result()
	if likeErr != nil && likeErr != redis.Nil {
		return Counts{}, likeErr
	}

	commentCount, commentErr := commentCmd.Result()
	if commentErr != nil && commentErr != redis.Nil {
		return Counts{}, commentErr
	}

	return Counts{
		VisitCount:   int64(visitCount),
		LikeCount:    int64(likeCount),
		CommentCount: int64(commentCount),
	}, nil
}
