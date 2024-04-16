package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"work4/pkg/errmsg"
)

var LikeToIdKey = "VideoLike"
var CommentToIdKey = "CommentToId"

type Counts struct {
	VisitCount   int64
	LikeCount    int64
	CommentCount int64
}

func CreatLikeCount(ctx context.Context, videoid string) error {

	videoResp := redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := redisDBVideoId.ZAdd(ctx, LikeToIdKey, videoResp).Result()

	if err != nil {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil

}

func AddLikeCount(ctx context.Context, videoid string) error {

	score, err := redisDBVideoId.ZScore(ctx, LikeToIdKey, videoid).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = CreatLikeCount(ctx, videoid)
			if err != nil && !errors.Is(err, redis.Nil) {
				return errmsg.RedisError.WithMessage(err.Error())
			}
		} else {
			return errmsg.RedisError.WithMessage(err.Error())
		}
	}

	videoResp := redis.Z{
		Score:  score + 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, LikeToIdKey, videoResp).Result()
	if err != nil {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil

}

func ReduceLikeCount(ctx context.Context, videoid string) error {

	score, err := redisDBVideoId.ZScore(ctx, LikeToIdKey, videoid).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	videoResp := redis.Z{
		Score:  score - 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, LikeToIdKey, videoResp).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil

}

func CreatCommentCount(ctx context.Context, videoid string) error {

	videoResp := redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := redisDBVideoId.ZAdd(ctx, CommentToIdKey, videoResp).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil

}

func AddCommentCount(ctx context.Context, videoid string) error {

	score, err := redisDBVideoId.ZScore(ctx, CommentToIdKey, videoid).Result()

	if err != nil {

		if errors.Is(err, redis.Nil) {
			err = CreatCommentCount(ctx, videoid)
			if err != nil {
				return errmsg.RedisError.WithMessage(err.Error())
			}
		} else {
			return errmsg.RedisError.WithMessage(err.Error())
		}

	}

	videoResp := redis.Z{
		Score:  score + 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, CommentToIdKey, videoResp).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil

}

func ReduceCommentCount(ctx context.Context, videoid string) error {

	score, err := redisDBVideoId.ZScore(ctx, CommentToIdKey, videoid).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	videoResp := redis.Z{
		Score:  score - 1,
		Member: videoid,
	}

	_, err = redisDBVideoId.ZAdd(ctx, CommentToIdKey, videoResp).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errmsg.RedisError.WithMessage(err.Error())
	}

	return nil

}

func GetCounts(ctx context.Context, videoId string) (Counts, error) {

	pipe := redisDBVideoId.Pipeline()

	visitCmd := pipe.ZScore(ctx, VideoIdKey, videoId)
	likeCmd := pipe.ZScore(ctx, LikeToIdKey, videoId)
	commentCmd := pipe.ZScore(ctx, CommentToIdKey, videoId)

	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return Counts{}, errmsg.RedisError.WithMessage(err.Error())
	}

	visitCount, visitErr := visitCmd.Result()
	if visitErr != nil && !errors.Is(visitErr, redis.Nil) {
		return Counts{}, errmsg.RedisError.WithMessage(visitErr.Error())
	}

	likeCount, likeErr := likeCmd.Result()
	if likeErr != nil && !errors.Is(likeErr, redis.Nil) {
		return Counts{}, errmsg.RedisError.WithMessage(likeErr.Error())
	}

	commentCount, commentErr := commentCmd.Result()
	if commentErr != nil && !errors.Is(commentErr, redis.Nil) {
		return Counts{}, errmsg.RedisError.WithMessage(commentErr.Error())
	}

	return Counts{
		VisitCount:   int64(visitCount),
		LikeCount:    int64(likeCount),
		CommentCount: int64(commentCount),
	}, nil
}
