package httpapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/catalog"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
)

type noopBookRepository struct{}

func (noopBookRepository) ListBooks(ctx context.Context) ([]catalog.Book, error) {
	return nil, nil
}

type noopChapterRepository struct{}

func (noopChapterRepository) FindChapter(ctx context.Context, bookID string, number int) (catalog.Chapter, error) {
	return catalog.Chapter{}, catalog.ErrChapterNotFound
}

func TestHealthEndpoint(t *testing.T) {
	router := httpapi.NewRouter(httpapi.Dependencies{
		Health: health.CheckHealth{},
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  catalog.NewListBooks(noopBookRepository{}),
			GetChapter: catalog.NewGetChapter(noopChapterRepository{}),
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", got)
	}

	want := `{"status":"ok"}`
	if got := rec.Body.String(); got != want {
		t.Errorf("expected body %q, got %q", want, got)
	}
}
