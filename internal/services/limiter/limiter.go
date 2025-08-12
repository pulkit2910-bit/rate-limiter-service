package service

import (
	"context"

	"github.com/pulkit2910-bit/rate-limiter-service/internal/redis"
)

type Service struct {
	rdb redis.RedisClient
}

type LimiterService interface {
    RunLuaScript(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error)
}

func NewLimiterService(rdb redis.RedisClient) LimiterService {
	return &Service{rdb: rdb}
}

func (s *Service) RunLuaScript(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error) {
    return s.rdb.Eval(ctx, script, keys, args...)
}