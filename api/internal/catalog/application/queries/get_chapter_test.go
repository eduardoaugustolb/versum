package queries_test

import (
	"context"
	"errors"
	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"testing"
)

func TestGetChapter(t *testing.T) {
	want := ports.ChapterDTO{BookID: "gn", BookName: "Gênesis", Number: 9}
	got, err := queries.NewGetChapter(catalogReader{chapter: want}).Execute(context.Background(), "gn", 9)
	if err != nil || got.BookID != want.BookID {
		t.Fatalf("unexpected result: %+v, %v", got, err)
	}
}
func TestGetChapterNotFound(t *testing.T) {
	_, err := queries.NewGetChapter(catalogReader{err: domain.ErrChapterNotFound}).Execute(context.Background(), "xx", 1)
	if !errors.Is(err, domain.ErrChapterNotFound) {
		t.Fatal(err)
	}
}

func TestGetChapterRejectsInvalidInputWithoutHTTP(t *testing.T) {
	for _, test := range []struct {
		book   string
		number int
		want   error
	}{
		{" ", 1, domain.ErrInvalidBookID}, {"gn", 0, domain.ErrInvalidChapter}, {"gn", -1, domain.ErrInvalidChapter},
	} {
		// A nil reader ensures invalid input cannot reach persistence.
		_, err := queries.NewGetChapter(nil).Execute(t.Context(), test.book, test.number)
		if !errors.Is(err, test.want) {
			t.Fatalf("want %v, got %v", test.want, err)
		}
	}
}
