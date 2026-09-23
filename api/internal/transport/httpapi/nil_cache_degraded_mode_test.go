package httpapi_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/cache"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
)

func TestNilCacheDegradedMode(t *testing.T) {
	handler := httpapi.NewHandler(httpapi.Dependencies{
		Health: health.CheckHealth{},
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  queries.NewListBooks(noopCatalogRepository{}),
			GetChapter: queries.NewGetChapter(noopCatalogRepository{}),
		},
		Cache: nil,
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("health: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/books", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("catalog: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":"ana@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(rec, request)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("magic-link: expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}

// flakyCache succeeds on the first n batches (global + IP limiters) and
// fails afterwards, simulating a transient Redis outage between the IP and
// email rate-limit checks.
type flakyCache struct {
	calls     atomic.Int64
	succeedOn int64
}

func (f *flakyCache) Batch() cache.Batch { return &flakyBatch{parent: f} }

func (f *flakyCache) Get(context.Context, string) ([]byte, error) { return nil, nil }

func (f *flakyCache) Set(context.Context, string, []byte, time.Duration) error {
	return nil
}

func (f *flakyCache) Delete(context.Context, ...string) error { return nil }

func (f *flakyCache) Expire(context.Context, string, time.Duration) error { return nil }

type flakyBatch struct {
	parent *flakyCache
}

func (b *flakyBatch) Get(context.Context, string) *cache.Result[[]byte] {
	return cache.NewResult[[]byte](nil, nil)
}

func (b *flakyBatch) Set(context.Context, string, []byte, time.Duration) *cache.Result[struct{}] {
	return cache.NewResult(struct{}{}, nil)
}

func (b *flakyBatch) Delete(context.Context, ...string) *cache.Result[int64] {
	return cache.NewResult[int64](0, nil)
}

func (b *flakyBatch) Increment(context.Context, string) *cache.Result[int64] {
	return cache.NewResult[int64](1, nil)
}

func (b *flakyBatch) TTL(context.Context, string) *cache.Result[time.Duration] {
	return cache.NewResult(time.Minute, nil)
}

func (b *flakyBatch) Exec(context.Context) error {
	if b.parent.calls.Add(1) > b.parent.succeedOn {
		return errors.New("cache unavailable")
	}
	return nil
}

func TestEmailRateLimitCacheFailureReturns503(t *testing.T) {
	flaky := &flakyCache{succeedOn: 2}
	handler := httpapi.NewHandler(httpapi.Dependencies{
		Health: health.CheckHealth{},
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  queries.NewListBooks(noopCatalogRepository{}),
			GetChapter: queries.NewGetChapter(noopCatalogRepository{}),
		},
		Cache: flaky,
	})

	rec := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/auth/magic-link", strings.NewReader(`{"email":"ana@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(rec, request)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("email limiter: expected 503, got %d: %s", rec.Code, rec.Body.String())
	}
}
