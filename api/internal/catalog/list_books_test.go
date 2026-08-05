package catalog_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/catalog"
)

type stubBookRepository struct {
	books []catalog.Book
	err   error
}

func (s stubBookRepository) ListBooks(ctx context.Context) ([]catalog.Book, error) {
	return s.books, s.err
}

func TestListBooksExecute(t *testing.T) {
	want := []catalog.Book{
		{ID: "gn", Order: 1, Name: "Gênesis", Testament: catalog.TestamentOld, ChapterCount: 50},
		{ID: "ap", Order: 73, Name: "Apocalipse", Testament: catalog.TestamentNew, ChapterCount: 22},
	}
	uc := catalog.NewListBooks(stubBookRepository{books: want})

	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %+v, got %+v", want, got)
	}
}

func TestListBooksExecutePropagatesRepositoryError(t *testing.T) {
	repoErr := errors.New("repository failure")
	uc := catalog.NewListBooks(stubBookRepository{err: repoErr})

	_, err := uc.Execute(context.Background())
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected error %v, got %v", repoErr, err)
	}
}
