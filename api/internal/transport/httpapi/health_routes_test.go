package httpapi_test

import (
	"context"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
	"net/http"
	"net/http/httptest"
	"testing"
)

type noopReader struct{}

func (noopReader) ListBooks(context.Context) ([]domain.Book, error) { return nil, nil }
func (noopReader) FindChapter(context.Context, string, int) (domain.Chapter, error) {
	return domain.Chapter{}, domain.ErrChapterNotFound
}
func TestHealthEndpoint(t *testing.T) {
	h := httpapi.NewRouter(httpapi.Dependencies{Health: health.CheckHealth{}, Catalog: httpapi.CatalogDependencies{ListBooks: queries.NewListBooks(noopReader{}), GetChapter: queries.NewGetChapter(noopReader{})}})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != `{"status":"ok"}` {
		t.Fatalf("unexpected response: %d %s", rec.Code, rec.Body)
	}
}
