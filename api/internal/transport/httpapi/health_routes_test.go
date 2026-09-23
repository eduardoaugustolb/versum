package httpapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	redisadapter "github.com/eduardoaugustolb/versum/api/internal/cache/redis"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
	redisclient "github.com/redis/go-redis/v9"
)

type noopCatalogRepository struct{}

func (noopCatalogRepository) ReplaceBook(context.Context, domain.Book, []domain.Verse) error {
	return nil
}
func (noopCatalogRepository) ListBooks(context.Context) ([]domain.Book, error) { return nil, nil }
func (noopCatalogRepository) FindChapter(context.Context, string, int) (domain.Chapter, error) {
	return domain.Chapter{}, domain.ErrChapterNotFound
}

func newHealthHandler(t *testing.T) http.Handler {
	t.Helper()

	redisServer := miniredis.RunT(t)
	redisClient := redisclient.NewClient(&redisclient.Options{Addr: redisServer.Addr()})
	t.Cleanup(func() { _ = redisClient.Close() })

	return httpapi.NewHandler(httpapi.Dependencies{
		Health: health.CheckHealth{},
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  queries.NewListBooks(noopCatalogRepository{}),
			GetChapter: queries.NewGetChapter(noopCatalogRepository{}),
		},
		Cache: redisadapter.NewRedisCache(redisClient),
	})
}

func TestHealthEndpoint(t *testing.T) {
	h := newHealthHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != `{"status":"ok"}` {
		t.Fatalf("unexpected response: %d %s", rec.Code, rec.Body)
	}
}

func TestDefaultRateLimitBlocksTheFiftyFirstRequest(t *testing.T) {
	h := newHealthHandler(t)

	for requestNumber := 1; requestNumber <= 50; requestNumber++ {
		rec := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/health", nil)
		request.RemoteAddr = "198.51.100.2:54321"
		h.ServeHTTP(rec, request)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d: %s", requestNumber, rec.Code, rec.Body.String())
		}
	}

	rec := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.RemoteAddr = "198.51.100.2:54321"
	h.ServeHTTP(rec, request)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d: %s", rec.Code, rec.Body.String())
	}
}
