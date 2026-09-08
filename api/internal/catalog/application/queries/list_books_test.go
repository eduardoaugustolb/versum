package queries_test

import (
	"context"
	"errors"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"testing"
)

type catalogRepository struct {
	books   []domain.Book
	chapter domain.Chapter
	err     error
}

func (r catalogRepository) ReplaceBook(context.Context, domain.Book, []domain.Verse) error {
	return r.err
}
func (r catalogRepository) ListBooks(context.Context) ([]domain.Book, error) { return r.books, r.err }
func (r catalogRepository) FindChapter(context.Context, string, int) (domain.Chapter, error) {
	return r.chapter, r.err
}
func TestListBooks(t *testing.T) {
	book, _ := domain.NewBook(domain.NewBookParams{ID: "gn", Order: 1, Name: "Gênesis", Testament: domain.TestamentOld, ChapterCount: 50})
	got, err := queries.NewListBooks(catalogRepository{books: []domain.Book{book}}).Execute(context.Background())
	if err != nil || len(got) != 1 || got[0].ID() != "gn" {
		t.Fatalf("unexpected result: %+v, %v", got, err)
	}
}
func TestListBooksPropagatesError(t *testing.T) {
	want := errors.New("repository failure")
	_, err := queries.NewListBooks(catalogRepository{err: want}).Execute(context.Background())
	if !errors.Is(err, want) {
		t.Fatal(err)
	}
}
