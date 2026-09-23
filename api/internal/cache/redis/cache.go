package redis

import (
	"context"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/cache"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

var _ cache.Cache = &RedisCache{}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) Batch() cache.Batch {
	pipe := r.client.Pipeline()
	return &RedisBatch{
		pipe: pipe,
	}
}

func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	err := r.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}
