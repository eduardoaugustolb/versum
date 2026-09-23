package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/eduardoaugustolb/versum/api/internal/cache"
	redisclient "github.com/redis/go-redis/v9"
)

func newTestBatch(t *testing.T) *RedisBatch {
	t.Helper()

	server := miniredis.RunT(t)
	client := redisclient.NewClient(&redisclient.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	return NewRedisBatch(client.Pipeline())
}

func TestRedisBatchResolvesResultsAfterExec(t *testing.T) {
	batch := newTestBatch(t)
	ctx := t.Context()

	set := batch.Set(ctx, "rate-limit:1", []byte("3"), time.Hour)
	get := batch.Get(ctx, "rate-limit:1")
	deleted := batch.Delete(ctx, "rate-limit:1")

	if value, err := get.Value(); err != nil || value != nil {
		t.Fatalf("Get result before Exec() = %q, %v; want nil, nil", value, err)
	}

	if err := batch.Exec(ctx); err != nil {
		t.Fatalf("Exec() error = %v", err)
	}
	if _, err := set.Value(); err != nil {
		t.Fatalf("Set result error = %v", err)
	}
	if value, err := get.Value(); err != nil || string(value) != "3" {
		t.Fatalf("Get result = %q, %v; want %q, nil", value, err, "3")
	}
	if count, err := deleted.Value(); err != nil || count != 1 {
		t.Fatalf("Delete result = %d, %v; want 1, nil", count, err)
	}
}

func TestRedisBatchResolvesIncrementAndTTL(t *testing.T) {
	batch := newTestBatch(t)
	ctx := t.Context()

	set := batch.Set(ctx, "counter", []byte("0"), time.Hour)
	incremented := batch.Increment(ctx, "counter")
	ttl := batch.TTL(ctx, "counter")

	if err := batch.Exec(ctx); err != nil {
		t.Fatalf("Exec() error = %v", err)
	}
	if _, err := set.Value(); err != nil {
		t.Fatalf("Set result error = %v", err)
	}
	if value, err := incremented.Value(); err != nil || value != 1 {
		t.Fatalf("Increment result = %d, %v; want 1, nil", value, err)
	}
	if value, err := ttl.Value(); err != nil || value <= 0 || value > time.Hour {
		t.Fatalf("TTL result = %v, %v; want duration within (0, %v]", value, err, time.Hour)
	}
}

func TestRedisBatchPreservesCommandErrorInResult(t *testing.T) {
	batch := newTestBatch(t)
	missing := batch.Get(t.Context(), "missing")

	err := batch.Exec(t.Context())
	if !errors.Is(err, redisclient.Nil) {
		t.Fatalf("Exec() error = %v, want redis.Nil", err)
	}
	if value, resultErr := missing.Value(); value != nil || !errors.Is(resultErr, redisclient.Nil) {
		t.Fatalf("Get result = %q, %v; want nil, redis.Nil", value, resultErr)
	}
}

func TestRedisBatchPropagatesCanceledContext(t *testing.T) {
	batch := newTestBatch(t)
	ctx, cancel := context.WithCancel(t.Context())
	batch.Set(ctx, "key", []byte("value"), time.Hour)
	cancel()

	err := batch.Exec(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Exec() error = %v, want context.Canceled", err)
	}
}

func TestRedisCacheCreatesIndependentConcurrentBatches(t *testing.T) {
	redisCache, _ := newTestCache(t)
	const workers = 32

	errs := make(chan error, workers)
	var group sync.WaitGroup
	for i := range workers {
		group.Add(1)
		go func() {
			defer group.Done()

			key := fmt.Sprintf("key:%d", i)
			value := fmt.Appendf(nil, "value:%d", i)
			batch := redisCache.Batch()
			set := batch.Set(t.Context(), key, value, time.Hour)
			get := batch.Get(t.Context(), key)

			if err := batch.Exec(t.Context()); err != nil {
				errs <- err
				return
			}
			if _, err := set.Value(); err != nil {
				errs <- err
				return
			}
			if actual, err := get.Value(); err != nil || string(actual) != string(value) {
				errs <- fmt.Errorf("Get(%q) = %q, %v; want %q, nil", key, actual, err, value)
			}
		}()
	}

	group.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

var _ cache.Batch = (*RedisBatch)(nil)
