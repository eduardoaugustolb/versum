package catalog_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	catalogpostgres "github.com/eduardoaugustolb/versum/api/internal/catalog/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	_, _ = pool.Exec(ctx, "DELETE FROM verses")
	_, _ = pool.Exec(ctx, "DELETE FROM books")
	if _, err = pool.Exec(ctx, `INSERT INTO books (id,"order",name,testament,chapter_count) VALUES ('gn',1,'Gênesis','old',50)`); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO verses (book_id,chapter,number,part,text) VALUES ('gn',9,1,1,'primeiro versículo'),('gn',9,9,2,'segundo'),('gn',9,9,1,'primeiro')`); err != nil {
		t.Fatal(err)
	}
	repo := catalogpostgres.NewRepository(postgres.NewPgxExecutor(pool))
	books, err := repo.ListBooks(ctx)
	if err != nil || len(books) != 1 || books[0].ID() != "gn" {
		t.Fatalf("unexpected books: %+v %v", books, err)
	}
	chapter, err := repo.FindChapter(ctx, "gn", 9)
	if err != nil || chapter.BookName() != "Gênesis" || len(chapter.Verses()) != 3 {
		t.Fatalf("unexpected chapter: %+v %v", chapter, err)
	}
	if chapter.Verses()[1].Part() != 1 {
		t.Fatalf("expected verses ordered by part: %+v", chapter.Verses())
	}
	_, err = repo.FindChapter(ctx, "gn", 9999)
	if !errors.Is(err, domain.ErrChapterNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	book, _ := domain.NewBook(domain.NewBookParams{ID: "ex", Order: 2, Name: "Êxodo", Testament: domain.TestamentOld, ChapterCount: 40})
	verse, _ := domain.NewVerse(domain.NewVerseParams{BookID: "ex", Chapter: 1, Number: 1, Text: "primeira", Part: 1})
	if err := repo.ReplaceBook(ctx, book, []domain.Verse{verse}); err != nil {
		t.Fatal(err)
	}
}
