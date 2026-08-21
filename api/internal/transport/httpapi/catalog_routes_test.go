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

type booksReader struct{ books []domain.Book }

func (r booksReader) ListBooks(context.Context) ([]domain.Book, error) { return r.books, nil }

type chaptersReader struct {
	chapter domain.Chapter
	found   bool
}

func (r chaptersReader) FindChapter(context.Context, string, int) (domain.Chapter, error) {
	if !r.found {
		return domain.Chapter{}, domain.ErrChapterNotFound
	}
	return r.chapter, nil
}
func router(br booksReader, cr chaptersReader) http.Handler {
	return httpapi.NewRouter(httpapi.Dependencies{Health: health.CheckHealth{}, Catalog: httpapi.CatalogDependencies{ListBooks: queries.NewListBooks(br), GetChapter: queries.NewGetChapter(cr)}})
}
func TestListBooksEndpoint(t *testing.T) {
	book, _ := domain.NewBook(domain.NewBookParams{ID: "gn", Order: 1, Name: "Gênesis", Testament: domain.TestamentOld, ChapterCount: 50})
	rec := httptest.NewRecorder()
	router(booksReader{[]domain.Book{book}}, chaptersReader{}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/books", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != `[{"id":"gn","order":1,"name":"Gênesis","testament":"old","chapter_count":50}]` {
		t.Fatalf("unexpected response: %d %s", rec.Code, rec.Body)
	}
}
func TestGetChapterEndpoint(t *testing.T) {
	chapter := domain.NewChapter("gn", "Gênesis", 9, nil)
	rec := httptest.NewRecorder()
	router(booksReader{}, chaptersReader{chapter: chapter, found: true}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/books/gn/chapters/9", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != `{"book_id":"gn","book_name":"Gênesis","number":9,"verses":[]}` {
		t.Fatalf("unexpected response: %d %s", rec.Code, rec.Body)
	}
}
func TestGetChapterEndpointNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	router(booksReader{}, chaptersReader{}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/books/xx/chapters/1", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
func TestGetChapterEndpointInvalidNumber(t *testing.T) {
	rec := httptest.NewRecorder()
	router(booksReader{}, chaptersReader{}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/books/gn/chapters/0", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
