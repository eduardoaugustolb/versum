package catalog_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/eduardoaugustolb/versum/api/internal/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/catalog"
)

func TestRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, "DELETE FROM verses"); err != nil {
		t.Fatalf("failed to clean verses: %v", err)
	}
	if _, err := pool.Exec(ctx, "DELETE FROM books"); err != nil {
		t.Fatalf("failed to clean books: %v", err)
	}

	if _, err := pool.Exec(ctx, `INSERT INTO books (id, "order", name, testament, chapter_count) VALUES ('gn', 1, 'Gênesis', 'old', 50)`); err != nil {
		t.Fatalf("failed to seed book: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO verses (book_id, chapter, number, part, text) VALUES
		('gn', 9, 1, 1, 'primeiro versículo'),
		('gn', 9, 9, 2, 'nono versículo, segunda parte'),
		('gn', 9, 9, 1, 'nono versículo, primeira parte')
	`); err != nil {
		t.Fatalf("failed to seed verses: %v", err)
	}

	repo := catalog.NewRepository(postgres.NewPgxExecutor(pool))

	t.Run("ListBooks", func(t *testing.T) {
		books, err := repo.ListBooks(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(books) != 1 || books[0].ID != "gn" {
			t.Fatalf("expected 1 book gn, got %+v", books)
		}
	})

	t.Run("FindChapter populates BookName and orders repeated verse numbers", func(t *testing.T) {
		chapter, err := repo.FindChapter(ctx, "gn", 9)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if chapter.BookName != "Gênesis" {
			t.Errorf("expected BookName %q, got %q", "Gênesis", chapter.BookName)
		}
		if len(chapter.Verses) != 3 {
			t.Fatalf("expected 3 verses, got %d", len(chapter.Verses))
		}
		if chapter.Verses[1].Number != 9 || chapter.Verses[1].Part != 1 {
			t.Errorf("expected verse 9.1 second, got %+v", chapter.Verses[1])
		}
		if chapter.Verses[2].Number != 9 || chapter.Verses[2].Part != 2 {
			t.Errorf("expected verse 9.2 third, got %+v", chapter.Verses[2])
		}
	})

	t.Run("FindChapter not found", func(t *testing.T) {
		_, err := repo.FindChapter(ctx, "gn", 9999)
		if !errors.Is(err, catalog.ErrChapterNotFound) {
			t.Fatalf("expected ErrChapterNotFound, got %v", err)
		}
	})
}
