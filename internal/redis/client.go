package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
    client *redis.Client
}

type RedisClient interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
	HSet(ctx context.Context, key string, value ...interface{}) error
}

func RedisClientStart() RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

    return &redisClient{
        client: rdb,
	}
}

func (r *redisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}

func (r *redisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values).Err()
}