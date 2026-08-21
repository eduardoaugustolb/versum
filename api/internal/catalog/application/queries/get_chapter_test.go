package queries_test

import (
	"context"
	"errors"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"testing"
)

type chapterReader struct {
	chapter domain.Chapter
	err     error
}

func (r chapterReader) FindChapter(context.Context, string, int) (domain.Chapter, error) {
	return r.chapter, r.err
}
func TestGetChapter(t *testing.T) {
	want := domain.NewChapter("gn", "Gênesis", 9, nil)
	got, err := queries.NewGetChapter(chapterReader{chapter: want}).Execute(context.Background(), "gn", 9)
	if err != nil || got.BookID() != want.BookID() {
		t.Fatalf("unexpected result: %+v, %v", got, err)
	}
}
func TestGetChapterNotFound(t *testing.T) {
	_, err := queries.NewGetChapter(chapterReader{err: domain.ErrChapterNotFound}).Execute(context.Background(), "xx", 1)
	if !errors.Is(err, domain.ErrChapterNotFound) {
		t.Fatal(err)
	}
}
