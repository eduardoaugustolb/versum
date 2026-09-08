package commands_test

import (
	"context"
	"errors"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"testing"
)

type catalogRepository struct {
	book   domain.Book
	verses []domain.Verse
	err    error
}

func (w *catalogRepository) ReplaceBook(_ context.Context, book domain.Book, verses []domain.Verse) error {
	w.book = book
	w.verses = verses
	return w.err
}

func (w *catalogRepository) ListBooks(context.Context) ([]domain.Book, error) { return nil, w.err }

func (w *catalogRepository) FindChapter(context.Context, string, int) (domain.Chapter, error) {
	return domain.Chapter{}, w.err
}

type txManager struct{ repository ports.CatalogRepository }

func (m txManager) WithinTransaction(ctx context.Context, fn func(context.Context, ports.CatalogRepository) error) error {
	return fn(ctx, m.repository)
}
func input() commands.PublishBookInput {
	return commands.PublishBookInput{ID: "gn", Order: 1, Name: "Gênesis", Testament: domain.TestamentOld, ChapterCount: 50, Verses: []commands.VerseInput{{BookID: "gn", Chapter: 1, Number: 1, Text: "No princípio", Part: 1}}}
}
func TestPublishBookBuildsDomainAndWrites(t *testing.T) {
	w := &catalogRepository{}
	if err := commands.NewPublishBook(txManager{repository: w}).Execute(context.Background(), input()); err != nil {
		t.Fatal(err)
	}
	if w.book.ID() != "gn" || len(w.verses) != 1 {
		t.Fatalf("unexpected write: %+v %+v", w.book, w.verses)
	}
}
func TestPublishBookValidatesInput(t *testing.T) {
	w := &catalogRepository{}
	in := input()
	in.Name = ""
	if err := commands.NewPublishBook(txManager{repository: w}).Execute(context.Background(), in); !errors.Is(err, domain.ErrInvalidBookName) {
		t.Fatal(err)
	}
}
func TestPublishBookPropagatesRepositoryError(t *testing.T) {
	want := errors.New("repository failure")
	if err := commands.NewPublishBook(txManager{repository: &catalogRepository{err: want}}).Execute(context.Background(), input()); !errors.Is(err, want) {
		t.Fatal(err)
	}
}
