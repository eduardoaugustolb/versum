package biblecorpus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCreatesCanonicalCorpusFromAllRawBooks(t *testing.T) {
	rawRoot := writeRawCorpus(t, 73)
	outputRoot := filepath.Join(t.TempDir(), "corpus", "v1")

	manifest, err := Generate(rawRoot, outputRoot)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if got, want := manifest.BookCount, 73; got != want {
		t.Fatalf("manifest bookCount = %d, want %d", got, want)
	}
	if got, want := manifest.VerseCount, 73; got != want {
		t.Fatalf("manifest verseCount = %d, want %d", got, want)
	}
	if got, want := manifest.NumberingGapCount, 0; got != want {
		t.Fatalf("manifest numberingGapCount = %d, want %d", got, want)
	}
	if got, want := manifest.RepeatedVerseNumberCount, 0; got != want {
		t.Fatalf("manifest repeatedVerseNumberCount = %d, want %d", got, want)
	}

	data, err := os.ReadFile(filepath.Join(outputRoot, "books", "01-b01.json"))
	if err != nil {
		t.Fatalf("read normalized book: %v", err)
	}

	var book Book
	if err := json.Unmarshal(data, &book); err != nil {
		t.Fatalf("decode normalized book: %v", err)
	}
	if got, want := book.Chapters[0].Verses[0].Text, "Texto 1."; got != want {
		t.Fatalf("normalized verse = %q, want %q", got, want)
	}

	if _, err := os.Stat(filepath.Join(outputRoot, "bible.json")); err != nil {
		t.Fatalf("complete canonical corpus: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputRoot, "manifest.json")); err != nil {
		t.Fatalf("manifest: %v", err)
	}
}

func TestGenerateRejectsAnIncompleteRawCatalog(t *testing.T) {
	_, err := Generate(writeRawCorpus(t, 72), filepath.Join(t.TempDir(), "corpus", "v1"))
	if err == nil {
		t.Fatal("Generate() error = nil, want error for a catalog with fewer than 73 books")
	}
}

func writeRawCorpus(t *testing.T, bookCount int) string {
	t.Helper()

	root := t.TempDir()
	booksDirectory := filepath.Join(root, "antigotestamento")
	if err := os.MkdirAll(booksDirectory, 0o755); err != nil {
		t.Fatalf("create raw book directory: %v", err)
	}

	catalog := make([]rawCatalogEntry, 0, bookCount)
	complete := make([]RawBook, 0, bookCount)
	for index := 1; index <= bookCount; index++ {
		id := fmt.Sprintf("b%02d", index)
		book := RawBook{
			Name: fmt.Sprintf("Livro %d", index),
			Chapters: []RawChapter{{
				Number: 1,
				Verses: []RawVerse{{Number: 1, Text: fmt.Sprintf("[%d] Texto %d.", 1, index)}},
			}},
		}

		catalog = append(catalog, rawCatalogEntry{ID: id, ChapterCount: 1})
		complete = append(complete, book)
		writeJSON(t, filepath.Join(booksDirectory, id+".json"), book)
	}

	writeJSON(t, filepath.Join(root, "listalivros.json"), catalog)
	complete[0].Chapters[0].Verses[0].Text = "Texto 1. (ver nota)"
	writeJSON(t, filepath.Join(root, "biblia.json"), complete)
	return root
}

func writeJSON(t *testing.T, path string, value any) {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal %s: %v", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
