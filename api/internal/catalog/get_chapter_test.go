package catalog_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/catalog"
)

type stubChapterRepository struct {
	chapter catalog.Chapter
	err     error
}

func (s stubChapterRepository) FindChapter(ctx context.Context, bookID string, number int) (catalog.Chapter, error) {
	return s.chapter, s.err
}

func TestGetChapterExecute(t *testing.T) {
	want := catalog.Chapter{
		BookID:   "gn",
		BookName: "Gênesis",
		Number:   9,
		Verses: []catalog.Verse{
			{BookID: "gn", Chapter: 9, Number: 1, Part: 1, Text: "primeiro versículo"},
			{BookID: "gn", Chapter: 9, Number: 9, Part: 1, Text: "nono versículo, primeira parte"},
			{BookID: "gn", Chapter: 9, Number: 9, Part: 2, Text: "nono versículo, segunda parte"},
			{BookID: "gn", Chapter: 9, Number: 10, Part: 1, Text: "décimo versículo"},
		},
	}

	uc := catalog.NewGetChapter(stubChapterRepository{chapter: want})

	got, err := uc.Execute(context.Background(), "gn", 9)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %+v, got %+v", want, got)
	}
}

func TestGetChapterExecuteNotFound(t *testing.T) {
	uc := catalog.NewGetChapter(stubChapterRepository{err: catalog.ErrChapterNotFound})

	_, err := uc.Execute(context.Background(), "xx", 1)
	if !errors.Is(err, catalog.ErrChapterNotFound) {
		t.Fatalf("expected %v, got %v", catalog.ErrChapterNotFound, err)
	}
}
