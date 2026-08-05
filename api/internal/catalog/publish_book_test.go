package catalog_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/catalog"
)

type stubCatalogWriter struct {
	gotBook   catalog.Book
	gotVerses []catalog.Verse
	err       error
}

func (s *stubCatalogWriter) ReplaceBook(ctx context.Context, book catalog.Book, verses []catalog.Verse) error {
	s.gotBook = book
	s.gotVerses = verses
	return s.err
}

func TestPublishBookExecute(t *testing.T) {
	book := catalog.Book{ID: "gn", Order: 1, Name: "Gênesis", Testament: catalog.TestamentOld, ChapterCount: 50}
	verses := []catalog.Verse{
		{BookID: "gn", Chapter: 1, Number: 1, Part: 1, Text: "No princípio criou Deus o céu e a terra."},
		{BookID: "gn", Chapter: 1, Number: 2, Part: 1, Text: "A terra, porém, estava informe e vazia."},
	}

	writer := &stubCatalogWriter{}
	uc := catalog.NewPublishBook(writer)

	err := uc.Execute(context.Background(), book, verses)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(writer.gotBook, book) {
		t.Errorf("expected book %+v, got %+v", book, writer.gotBook)
	}
	if !reflect.DeepEqual(writer.gotVerses, verses) {
		t.Errorf("expected verses %+v, got %+v", verses, writer.gotVerses)
	}
}

func TestPublishBookExecutePropagatesWriterError(t *testing.T) {
	writerErr := errors.New("writer failure")
	writer := &stubCatalogWriter{err: writerErr}
	uc := catalog.NewPublishBook(writer)

	err := uc.Execute(context.Background(), catalog.Book{}, nil)
	if !errors.Is(err, writerErr) {
		t.Fatalf("expected error %v, got %v", writerErr, err)
	}
}
