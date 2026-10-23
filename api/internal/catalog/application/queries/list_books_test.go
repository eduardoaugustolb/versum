package queries_test

import (
	"context"
	"errors"
	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"testing"
)

type catalogReader struct {
	books   []ports.BookDTO
	chapter ports.ChapterDTO
	err     error
}

func (r catalogReader) ListBooks(context.Context) ([]ports.BookDTO, error) { return r.books, r.err }
func (r catalogReader) FindChapter(context.Context, string, int) (ports.ChapterDTO, error) {
	return r.chapter, r.err
}
func TestListBooks(t *testing.T) {
	book := ports.BookDTO{ID: "gn", Order: 1, Name: "Gênesis", Testament: "old", ChapterCount: 50}
	got, err := queries.NewListBooks(catalogReader{books: []ports.BookDTO{book}}).Execute(context.Background())
	if err != nil || len(got) != 1 || got[0].ID != "gn" {
		t.Fatalf("unexpected result: %+v, %v", got, err)
	}
}
func TestListBooksPropagatesError(t *testing.T) {
	want := errors.New("repository failure")
	_, err := queries.NewListBooks(catalogReader{err: want}).Execute(context.Background())
	if !errors.Is(err, want) {
		t.Fatal(err)
	}
}
