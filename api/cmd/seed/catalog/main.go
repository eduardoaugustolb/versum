package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	fmt.Println("Iniciando...")
	fmt.Println()

	rootDir := flag.String("root", "..", "project root path")
	corpusPath := flag.String("corpus", DEFAULT_CORPUS_PATH, "bible corpus path")
	flag.Parse()

	corpusDir := filepath.Join(*rootDir, *corpusPath)
	bibleCorpusDir := filepath.Join(corpusDir, "bible.json")

	fmt.Printf("Corpus: %s\n", corpusDir)

	fmt.Println()

	fmt.Println("Lendo o corpus...")
	bibleCorpusJson, err := os.ReadFile(bibleCorpusDir)
	if err != nil {
		fmt.Printf("Lendo o corpus: Error: %v\n", err)
		return
	}

	fmt.Println("Corpus lido com sucesso!")

	fmt.Println()

	fmt.Println("Unmarshalling...")

	var bibleCorpus Corpus
	if err := json.Unmarshal(bibleCorpusJson, &bibleCorpus); err != nil {
		fmt.Printf("Unmarshalling: Error: %v\n", err)
		return
	}

	fmt.Println("Unmarshalling concluído!")
	fmt.Println()
	fmt.Println("SchemaVersion:", bibleCorpus.SchemaVersion)
	fmt.Println("Books:", len(bibleCorpus.Books))

	fmt.Println()

	godotenv.Load()
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		fmt.Printf("Configuração: Error: %v\n", err)
		return
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Printf("Configuração: Error: %v\n", err)
		return
	}
	defer pool.Close()

	publishBook := commands.NewPublishBook(catalogpostgres.NewTransactionManager(pool))

	processBooks(ctx, bibleCorpus.Books, publishBook)
}

func processBooks(ctx context.Context, books []Book, publishBook commands.PublishBook) {
	for _, b := range books {
		fmt.Println("Book:", b.Name)
		fmt.Println("Order:", b.Order)
		fmt.Println("Testament:", b.Testament)
		fmt.Println("Chapters:", len(b.Chapters))

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
			fmt.Printf("Publicando livro: Error: %v\n", err)
		}

		fmt.Println("Livro publicado com sucesso!")
		fmt.Println()
	}
}
