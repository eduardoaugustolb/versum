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

	t.Run("ReplaceBook publishes a new book and its verses", func(t *testing.T) {
		book := catalog.Book{ID: "ex", Order: 2, Name: "Êxodo", Testament: catalog.TestamentOld, ChapterCount: 40}
		verses := []catalog.Verse{
			{BookID: "ex", Chapter: 1, Number: 1, Part: 1, Text: "primeira versão do versículo 1"},
			{BookID: "ex", Chapter: 1, Number: 2, Part: 1, Text: "primeira versão do versículo 2"},
		}

		if err := repo.ReplaceBook(ctx, book, verses); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		chapter, err := repo.FindChapter(ctx, "ex", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(chapter.Verses) != 2 {
			t.Fatalf("expected 2 verses, got %d: %+v", len(chapter.Verses), chapter.Verses)
		}
	})

	t.Run("ReplaceBook replaces verses and book fields instead of accumulating them", func(t *testing.T) {
		// Publica uma primeira vez.
		book := catalog.Book{ID: "lv", Order: 3, Name: "Levítico", Testament: catalog.TestamentOld, ChapterCount: 27}
		firstVerses := []catalog.Verse{
			{BookID: "lv", Chapter: 1, Number: 1, Part: 1, Text: "texto original do versículo 1"},
			{BookID: "lv", Chapter: 1, Number: 2, Part: 1, Text: "texto original do versículo 2"},
		}
		if err := repo.ReplaceBook(ctx, book, firstVerses); err != nil {
			t.Fatalf("unexpected error on first publish: %v", err)
		}

		// Republica simulando uma revisão do corpus: um versículo removido,
		// outro com texto reescrito, e chapter_count diferente.
		updatedBook := catalog.Book{ID: "lv", Order: 3, Name: "Levítico", Testament: catalog.TestamentOld, ChapterCount: 26}
		secondVerses := []catalog.Verse{
			{BookID: "lv", Chapter: 1, Number: 1, Part: 1, Text: "texto revisado do versículo 1"},
		}
		if err := repo.ReplaceBook(ctx, updatedBook, secondVerses); err != nil {
			t.Fatalf("unexpected error on second publish: %v", err)
		}

		chapter, err := repo.FindChapter(ctx, "lv", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(chapter.Verses) != 1 {
			t.Fatalf("expected verse 2 to be gone after replace, got %d verses: %+v", len(chapter.Verses), chapter.Verses)
		}
		if chapter.Verses[0].Text != "texto revisado do versículo 1" {
			t.Errorf("expected updated verse text, got %q", chapter.Verses[0].Text)
		}

		books, err := repo.ListBooks(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var republished catalog.Book
		found := false
		for _, b := range books {
			if b.ID == "lv" {
				republished = b
				found = true
			}
		}
		if !found {
			t.Fatalf("expected book lv in ListBooks, got %+v", books)
		}
		if republished.ChapterCount != 26 {
			t.Errorf("expected book upserted with chapter_count 26, got %d", republished.ChapterCount)
		}
	})
}
