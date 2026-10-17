package catalog_test

import (
	"os"
	"testing"

	catalogpostgres "github.com/eduardoaugustolb/versum/api/internal/catalog/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"github.com/eduardoaugustolb/versum/api/internal/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCatalogPublicationCommitsContentAndVersionTogether(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := t.Context()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if _, err := pool.Exec(ctx, "TRUNCATE books, verses, catalog_version CASCADE"); err != nil {
		t.Fatal(err)
	}
	uc := commands.NewPublishCatalog(catalogpostgres.NewTransactionManager(postgres.NewPgxExecutor(pool)))
	book := commands.PublishBookInput{ID: "gn", Order: 1, Name: "Gênesis", Testament: domain.TestamentOld, ChapterCount: 1, Verses: []commands.VerseInput{{BookID: "gn", Chapter: 1, Number: 1, Text: "original", Part: 1}}}
	if err := uc.Execute(ctx, []commands.PublishBookInput{book}, "version-1"); err != nil {
		t.Fatal(err)
	}
	book.Verses[0].Text = "replacement"
	// Both books are valid in the domain; the second violates the unique order in storage.
	conflict := commands.PublishBookInput{ID: "ex", Order: 1, Name: "Êxodo", Testament: domain.TestamentOld, ChapterCount: 1}
	if err := uc.Execute(ctx, []commands.PublishBookInput{book, conflict}, "version-2"); err == nil {
		t.Fatal("expected storage conflict")
	}
	var text, version string
	if err := pool.QueryRow(ctx, "SELECT text FROM verses WHERE book_id='gn'").Scan(&text); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "SELECT corpus_sha256 FROM catalog_version").Scan(&version); err != nil {
		t.Fatal(err)
	}
	if text != "original" || version != "version-1" {
		t.Fatalf("partial publication: %q, %q", text, version)
	}
	// A failure while recording the version must also roll back earlier book writes.
	if _, err := pool.Exec(ctx, "ALTER TABLE catalog_version ADD CONSTRAINT reject_test_version CHECK (corpus_sha256 <> 'rejected')"); err != nil {
		t.Fatal(err)
	}
	defer pool.Exec(ctx, "ALTER TABLE catalog_version DROP CONSTRAINT reject_test_version")
	if err := uc.Execute(ctx, []commands.PublishBookInput{book}, "rejected"); err == nil {
		t.Fatal("expected version failure")
	}
	if err := pool.QueryRow(ctx, "SELECT text FROM verses WHERE book_id='gn'").Scan(&text); err != nil {
		t.Fatal(err)
	}
	if text != "original" {
		t.Fatal("version failure did not roll back content")
	}
	if err := uc.Execute(ctx, []commands.PublishBookInput{book}, "version-3"); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "SELECT corpus_sha256 FROM catalog_version").Scan(&version); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "SELECT text FROM verses WHERE book_id='gn'").Scan(&text); err != nil {
		t.Fatal(err)
	}
	if text != "replacement" || version != "version-3" {
		t.Fatalf("publication not committed: %q, %q", text, version)
	}
}
