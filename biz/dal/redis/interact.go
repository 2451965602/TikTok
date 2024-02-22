package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

var LikeKey = "LikeRank"

func CreatLikeCount(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	videoResp := redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := RDB.ZAdd(ctx, LikeKey, videoResp).Result()

	if err != nil {
		return err
	}

	return nil

}

func AddLikeCount(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	score, err := RDB.ZScore(ctx, LikeKey, videoid).Result()

	if err != nil {

		if err == redis.Nil {
			err = CreatLikeCount(ctx, videoid)
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

	_, err = RDB.ZAdd(ctx, LikeKey, videoResp).Result()

	if err != nil {
		return err
	}

	err = CreatVideoId(ctx, videoid)
	if err != nil {
		return err
	}

	return nil

}

func ReduceLikeCount(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	score, err := RDB.ZScore(ctx, LikeKey, videoid).Result()

	if err != nil {
		return err
	}

	videoResp := redis.Z{
		Score:  score - 1,
		Member: videoid,
	}

	_, err = RDB.ZAdd(ctx, LikeKey, videoResp).Result()

	if err != nil {
		return err
	}

	err = CreatVideoId(ctx, videoid)
	if err != nil {
		return err
	}

	return nil

}

func GetLikeCount(ctx context.Context, videoid string) (int64, error) {

	if RDB == nil {
		return -1, errors.New("RDB object is nil")
	}

	count, err := RDB.ZScore(ctx, LikeKey, videoid).Result()

	if err != nil {
		return -1, err
	}

	return int64(count), nil

}

var CommentKey = "LikeRank"

func CreatCommentCount(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	videoResp := redis.Z{
		Score:  0,
		Member: videoid,
	}

	_, err := RDB.ZAdd(ctx, CommentKey, videoResp).Result()

	if err != nil {
		return err
	}

	return nil

}

func AddCommentCount(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	score, err := RDB.ZScore(ctx, CommentKey, videoid).Result()

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

	_, err = RDB.ZAdd(ctx, CommentKey, videoResp).Result()

	if err != nil {
		return err
	}

	err = CreatVideoId(ctx, videoid)
	if err != nil {
		return err
	}

	return nil

}

func ReduceCommentCount(ctx context.Context, videoid string) error {

	if RDB == nil {
		return errors.New("RDB object is nil")
	}

	score, err := RDB.ZScore(ctx, CommentKey, videoid).Result()

	if err != nil {
		return err
	}

	videoResp := redis.Z{
		Score:  score - 1,
		Member: videoid,
	}

	_, err = RDB.ZAdd(ctx, CommentKey, videoResp).Result()

	if err != nil {
		return err
	}

	err = CreatVideoId(ctx, videoid)
	if err != nil {
		return err
	}

	return nil

}

func GetCommentCount(ctx context.Context, videoid string) (int64, error) {

	if RDB == nil {
		return -1, errors.New("RDB object is nil")
	}

	count, err := RDB.ZScore(ctx, CommentKey, videoid).Result()

	if err != nil {
		return -1, err
	}

	return int64(count), nil

}
