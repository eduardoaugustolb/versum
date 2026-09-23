package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/cache"
)

type RateLimit struct {
	generateID func(r *http.Request) string
	window     time.Duration
	limit      int
	cache      cache.Cache
	canSkip    bool
}

func NewRateLimit(generateID func(r *http.Request) string, window time.Duration, limit int, cache cache.Cache, canSkip bool) *RateLimit {
	if cache == nil {
		slog.Warn("rate limit cache not configured, limited requests will fail unless canSkip is set")
	}
	return &RateLimit{
		generateID: generateID,
		window:     window,
		limit:      limit,
		cache:      cache,
		canSkip:    canSkip,
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func (r *RateLimit) Allow(ctx context.Context, key string) (bool, error) {
	// Without cache the limit cannot be evaluated; fail closed so callers
	// honor canSkip (global limiter skips, auth limiter returns 503).
	if r == nil || r.cache == nil {
		return false, fmt.Errorf("rate limit cache not configured")
	}
	batch := r.cache.Batch()
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
		if err := r.cache.Expire(ctx, key, r.window); err != nil {
			return false, fmt.Errorf("expiring rate limit: %w", err)
		}
	}

	return int(count) <= r.limit, nil
}

func (r *RateLimit) Execute(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Nil receiver is a wiring bug (NewRateLimit never returns nil),
		// so canSkip is unknowable — always fail closed.
		if r == nil {
			slog.ErrorContext(req.Context(), "rate limit middleware not configured")
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		key := r.generateID(req) + ":ratelimit:" + clientIP(req)
		allowed, err := r.Allow(req.Context(), key)
		if err != nil && r.canSkip {
			handler.ServeHTTP(w, req)
			return
		}

		if err != nil {
			slog.ErrorContext(req.Context(), "failed to allow rate limit", "error", err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		if !allowed {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		handler.ServeHTTP(w, req)
	})
}
