package redis

import (
	"context"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/cache"
	"github.com/redis/go-redis/v9"
)

type RedisBatch struct {
	pipe    redis.Pipeliner
	pending []func()
}

var _ cache.Batch = &RedisBatch{}

func NewRedisBatch(pipe redis.Pipeliner) *RedisBatch {
	return &RedisBatch{pipe: pipe}
}

func (b *RedisBatch) Exec(ctx context.Context) error {
	_, err := b.pipe.Exec(ctx)
	for _, resolve := range b.pending {
		resolve()
	}
	b.pending = nil
	return err
}

func (b *RedisBatch) Get(ctx context.Context, key string) *cache.Result[[]byte] {
	cmd := b.pipe.Get(ctx, key)
	result := cache.NewResult([]byte(nil), nil)
	b.pending = append(b.pending, func() {
		value, err := cmd.Bytes()
		result.Resolve(value, err)
	})
	return result
}

func (b *RedisBatch) Set(ctx context.Context, key string, value []byte, ttl time.Duration) *cache.Result[struct{}] {
	cmd := b.pipe.Set(ctx, key, value, ttl)
	result := cache.NewResult(struct{}{}, nil)
	b.pending = append(b.pending, func() {
		result.Resolve(struct{}{}, cmd.Err())
	})
	return result
}

func (b *RedisBatch) Delete(ctx context.Context, key ...string) *cache.Result[int64] {
	cmd := b.pipe.Del(ctx, key...)
	result := cache.NewResult(int64(0), nil)
	b.pending = append(b.pending, func() {
		value, err := cmd.Result()
		result.Resolve(value, err)
	})
	return result
}

func (b *RedisBatch) Increment(ctx context.Context, key string) *cache.Result[int64] {
	cmd := b.pipe.Incr(ctx, key)
	result := cache.NewResult(int64(0), nil)
	b.pending = append(b.pending, func() {
		value, err := cmd.Result()
		result.Resolve(value, err)
	})
	return result
}

func (b *RedisBatch) TTL(ctx context.Context, key string) *cache.Result[time.Duration] {
	cmd := b.pipe.TTL(ctx, key)
	result := cache.NewResult(time.Duration(0), nil)
	b.pending = append(b.pending, func() {
		value, err := cmd.Result()
		result.Resolve(value, err)
	})
	return result
}
