package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eduardoaugustolb/versum/api/internal/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	catalogpostgres "github.com/eduardoaugustolb/versum/api/internal/catalog/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	DEFAULT_CORPUS_PATH = `bible/corpus/v1`
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("seed failed: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Iniciando...")
	fmt.Println()

	rootDir := flag.String("root", "..", "project root path")
	corpusPath := flag.String("corpus", DEFAULT_CORPUS_PATH, "bible corpus path")
	flag.Parse()

	corpusDir := filepath.Join(*rootDir, *corpusPath)
	corpusManifestDir := filepath.Join(corpusDir, "manifest.json")
	bibleCorpusDir := filepath.Join(corpusDir, "bible.json")

	fmt.Printf("Corpus: %s\n", corpusDir)

	fmt.Println()

	fmt.Println("reading corpus manifest")
	corpusManifestJson, err := os.ReadFile(corpusManifestDir)
	if err != nil {
		return fmt.Errorf("reading corpus manifest: %w", err)
	}
	var corpusManifest CorpusManifest
	if err := json.Unmarshal(corpusManifestJson, &corpusManifest); err != nil {
		return fmt.Errorf("unmarshalling corpus manifest: %w", err)
	}

	fmt.Println("success on read corpus manifest!")
	fmt.Println("corpus manifest sha256:", corpusManifest.BibleSha256)

	fmt.Println()

	fmt.Println("reading corpus...")
	bibleCorpusJson, err := os.ReadFile(bibleCorpusDir)
	if err != nil {
		return fmt.Errorf("reading corpus bible: %w", err)
	}

	fmt.Println("success on read corpus bible!")

	fmt.Println()

	fmt.Println("unmarshalling corpus bible...")

	var bibleCorpus CorpusBible
	if err := json.Unmarshal(bibleCorpusJson, &bibleCorpus); err != nil {
		return fmt.Errorf("unmarshalling corpus bible: %w", err)
	}

	fmt.Println("corpus bilbe unmarshalling finished!")
	fmt.Println()
	fmt.Println("schema version:", bibleCorpus.SchemaVersion)
	fmt.Println("books:", len(bibleCorpus.Books))

	fmt.Println()

	godotenv.Load()
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("creating database pool: %w", err)
	}
	defer pool.Close()

	catalogVersionWriter := catalogpostgres.NewCatalogVersionRepository(postgres.NewPgxExecutor(pool))
	publishCatalogVersion := commands.NewPublishCatalogVersion(catalogVersionWriter)
	publishBook := commands.NewPublishBook(catalogpostgres.NewTransactionManager(pool))

	if err := processBooks(ctx, bibleCorpus.Books, publishBook); err != nil {
		return err
	}

	if err := publishCatalogVersion.Execute(ctx, commands.PublishCatalogVersionInput{
		CorpusSHA256: corpusManifest.BibleSha256,
	}); err != nil {
		return fmt.Errorf("publishing catalog version: %w", err)
	}

	return nil
}

func processBooks(ctx context.Context, books []Book, publishBook commands.PublishBook) error {
	for _, b := range books {
		fmt.Println("book:", b.Name)
		fmt.Println("order:", b.Order)
		fmt.Println("testament:", b.Testament)
		fmt.Println("chapters:", len(b.Chapters))

		verses := make([]commands.VerseInput, 0, len(b.Chapters)*10)
		for _, c := range b.Chapters {
			for _, v := range c.Verses {
				verses = append(verses, commands.VerseInput{
					Number:  v.Number,
					Text:    v.Text,
					Part:    v.Part,
					Chapter: c.Number,
					BookID:  b.Id,
				})
			}
		}

		if err := publishBook.Execute(ctx, commands.PublishBookInput{
			ID: b.Id, Order: b.Order, Name: b.Name, Testament: domain.Testament(b.Testament),
			ChapterCount: len(b.Chapters), Verses: verses,
		}); err != nil {
			return fmt.Errorf("publishing book %q: %w", b.Id, err)
		}

		fmt.Println("book published!")
		fmt.Println()
	}

	return nil
}
