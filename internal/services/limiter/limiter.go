package service

import (
	"context"
	"fmt"

	"github.com/pulkit2910-bit/rate-limiter-service/internal/redis"
)

type Service struct {
	rdb redis.RedisClient
}

type LimiterService interface {
    RunLuaScript(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
	RateLimitingConfig(ctx context.Context, userId string, capacity string, refillRate string) error
}

func NewLimiterService(rdb redis.RedisClient) LimiterService {
	return &Service{rdb: rdb}
}

func (s *Service) RunLuaScript(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
    return s.rdb.Eval(ctx, script, keys, args...)
}

func (s *Service) RateLimitingConfig(ctx context.Context, userId string, capacity string, refillRate string) error {
    key := fmt.Sprintf("config:%s", userId)

	if err := s.rdb.HSet(ctx, key, "capacity", capacity, "refill_rate", refillRate); err != nil {
		fmt.Printf("Error setting rate limiting config: %v\n", err)
		return err
	}
	return nil
}