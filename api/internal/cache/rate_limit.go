package cache

import (
	"context"
	"fmt"
	"time"
)

func AllowWithinWindow(ctx context.Context, store Cache, key string, window time.Duration, limit int) (bool, error) {
	// Callers choose how to handle an unavailable cache; no quota is granted here.
	if store == nil {
		return false, fmt.Errorf("rate limit cache not configured")
	}
	batch := store.Batch()
	countResult := batch.Increment(ctx, key)
	ttlResult := batch.TTL(ctx, key)

	if err := batch.Exec(ctx); err != nil {
		return false, fmt.Errorf("incrementing rate limit count: %w", err)
	}

	count, err := countResult.Value()
	if err != nil {
		return false, fmt.Errorf("getting rate limit count: %w", err)
	}
	ttl, err := ttlResult.Value()
	if err != nil {
		return false, fmt.Errorf("getting rate limit ttl: %w", err)
	}

	if ttl < 0 {
		if err := store.Expire(ctx, key, window); err != nil {
			return false, fmt.Errorf("expiring rate limit: %w", err)
		}
	}

	return int(count) <= limit, nil
}
