package httpapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/catalog"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
)

type stubBookRepository struct {
	books []catalog.Book
}

func (s stubBookRepository) ListBooks(ctx context.Context) ([]catalog.Book, error) {
	return s.books, nil
}

type stubChapterRepository struct {
	chapters map[string]catalog.Chapter
}

func (s stubChapterRepository) FindChapter(ctx context.Context, bookID string, number int) (catalog.Chapter, error) {
	chapter, ok := s.chapters[fmt.Sprintf("%s:%d", bookID, number)]
	if !ok {
		return catalog.Chapter{}, catalog.ErrChapterNotFound
	}
	return chapter, nil
}

func newTestRouter(books []catalog.Book, chapters map[string]catalog.Chapter) http.Handler {
	return httpapi.NewRouter(httpapi.Dependencies{
		Health: health.CheckHealth{},
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  catalog.NewListBooks(stubBookRepository{books: books}),
			GetChapter: catalog.NewGetChapter(stubChapterRepository{chapters: chapters}),
		},
	})
}

func TestListBooksEndpoint(t *testing.T) {
	books := []catalog.Book{
		{ID: "gn", Order: 1, Name: "Gênesis", Testament: catalog.TestamentOld, ChapterCount: 50},
		{ID: "ap", Order: 73, Name: "Apocalipse", Testament: catalog.TestamentNew, ChapterCount: 22},
	}
	router := newTestRouter(books, nil)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []catalog.Book
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if !reflect.DeepEqual(got, books) {
		t.Errorf("expected %+v, got %+v", books, got)
	}
}

func TestGetChapterEndpoint(t *testing.T) {
	chapter := catalog.Chapter{
		BookID:   "gn",
		BookName: "Gênesis",
		Number:   9,
		Verses: []catalog.Verse{
			{BookID: "gn", Chapter: 9, Number: 1, Part: 1, Text: "primeiro versículo"},
			{BookID: "gn", Chapter: 9, Number: 9, Part: 1, Text: "nono versículo, primeira parte"},
			{BookID: "gn", Chapter: 9, Number: 9, Part: 2, Text: "nono versículo, segunda parte"},
		},
	}
	router := newTestRouter(nil, map[string]catalog.Chapter{"gn:9": chapter})

	req := httptest.NewRequest(http.MethodGet, "/books/gn/chapters/9", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got catalog.Chapter
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if !reflect.DeepEqual(got, chapter) {
		t.Errorf("expected %+v, got %+v", chapter, got)
	}
}

func TestGetChapterEndpointBookNotFound(t *testing.T) {
	router := newTestRouter(nil, map[string]catalog.Chapter{})

	req := httptest.NewRequest(http.MethodGet, "/books/xx/chapters/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetChapterEndpointChapterNotFound(t *testing.T) {
	router := newTestRouter(nil, map[string]catalog.Chapter{"gn:1": {}})

	req := httptest.NewRequest(http.MethodGet, "/books/gn/chapters/9999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetChapterEndpointInvalidNumber(t *testing.T) {
	tests := []string{"0", "-1", "abc"}

	for _, number := range tests {
		t.Run(number, func(t *testing.T) {
			router := newTestRouter(nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/books/gn/chapters/"+number, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}
