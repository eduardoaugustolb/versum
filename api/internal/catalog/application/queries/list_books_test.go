package queries_test

import (
	"context"
	"errors"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"testing"
)

type bookReader struct {
	books []domain.Book
	err   error
}

func (r bookReader) ListBooks(context.Context) ([]domain.Book, error) { return r.books, r.err }
func TestListBooks(t *testing.T) {
	book, _ := domain.NewBook(domain.NewBookParams{ID: "gn", Order: 1, Name: "Gênesis", Testament: domain.TestamentOld, ChapterCount: 50})
	got, err := queries.NewListBooks(bookReader{books: []domain.Book{book}}).Execute(context.Background())
	if err != nil || len(got) != 1 || got[0].ID() != "gn" {
		t.Fatalf("unexpected result: %+v, %v", got, err)
	}
}
func TestListBooksPropagatesError(t *testing.T) {
	want := errors.New("repository failure")
	_, err := queries.NewListBooks(bookReader{err: want}).Execute(context.Background())
	if !errors.Is(err, want) {
		t.Fatal(err)
	}
}
