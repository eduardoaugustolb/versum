package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestCache(t *testing.T) (*RedisCache, *miniredis.Miniredis) {
	t.Helper()

	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	return NewRedisCache(client), server
}

func TestRedisCacheSetGetDelete(t *testing.T) {
	cache, _ := newTestCache(t)
	ctx := t.Context()

	if err := cache.Set(ctx, "session:1", []byte("value"), time.Hour); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	value, err := cache.Get(ctx, "session:1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(value) != "value" {
		t.Fatalf("Get() = %q, want %q", value, "value")
	}

	if err := cache.Delete(ctx, "session:1"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := cache.Get(ctx, "session:1"); !errors.Is(err, redis.Nil) {
		t.Fatalf("Get() after Delete() error = %v, want redis.Nil", err)
	}
}

func TestRedisCacheDeleteAcceptsMultipleKeys(t *testing.T) {
	cache, _ := newTestCache(t)
	ctx := t.Context()

	for _, key := range []string{"key:1", "key:2"} {
		if err := cache.Set(ctx, key, []byte("value"), 0); err != nil {
			t.Fatalf("Set(%q) error = %v", key, err)
		}
	}

	if err := cache.Delete(ctx, "key:1", "key:2"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	for _, key := range []string{"key:1", "key:2"} {
		if _, err := cache.Get(ctx, key); !errors.Is(err, redis.Nil) {
			t.Fatalf("Get(%q) error = %v, want redis.Nil", key, err)
		}
	}
}

func TestRedisCacheExpiresValues(t *testing.T) {
	cache, server := newTestCache(t)
	ctx := t.Context()

	if err := cache.Set(ctx, "temporary", []byte("value"), time.Minute); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	server.FastForward(time.Minute)

	if _, err := cache.Get(ctx, "temporary"); !errors.Is(err, redis.Nil) {
		t.Fatalf("Get() after expiration error = %v, want redis.Nil", err)
	}
}

func TestRedisCachePropagatesCanceledContext(t *testing.T) {
	cache, _ := newTestCache(t)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if err := cache.Set(ctx, "key", []byte("value"), time.Hour); !errors.Is(err, context.Canceled) {
		t.Fatalf("Set() error = %v, want context.Canceled", err)
	}
	if _, err := cache.Get(ctx, "key"); !errors.Is(err, context.Canceled) {
		t.Fatalf("Get() error = %v, want context.Canceled", err)
	}
	if err := cache.Delete(ctx, "key"); !errors.Is(err, context.Canceled) {
		t.Fatalf("Delete() error = %v, want context.Canceled", err)
	}
}
