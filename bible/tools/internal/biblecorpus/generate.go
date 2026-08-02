package biblecorpus

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	canonicalSchemaVersion = 1
	expectedBookCount      = 73
)

type rawCatalogEntry struct {
	ID           string `json:"livro"`
	ChapterCount int    `json:"quantidadeCap"`
}

type Corpus struct {
	SchemaVersion int    `json:"schemaVersion"`
	Books         []Book `json:"books"`
}

type Manifest struct {
	SchemaVersion            int            `json:"schemaVersion"`
	BookCount                int            `json:"bookCount"`
	ChapterCount             int            `json:"chapterCount"`
	VerseCount               int            `json:"verseCount"`
	NumberingGapCount        int            `json:"numberingGapCount"`
	RepeatedVerseNumberCount int            `json:"repeatedVerseNumberCount"`
	BibleSHA256              string         `json:"bibleSha256"`
	Books                    []ManifestBook `json:"books"`
}

type ManifestBook struct {
	ID       string `json:"id"`
	Order    int    `json:"order"`
	Chapters int    `json:"chapters"`
	Verses   int    `json:"verses"`
	SHA256   string `json:"sha256"`
}

func Generate(rawRoot, outputRoot string) (Manifest, error) {
	catalog, err := readCatalog(filepath.Join(rawRoot, "listalivros.json"))
	if err != nil {
		return Manifest{}, err
	}
	if len(catalog) != expectedBookCount {
		return Manifest{}, fmt.Errorf("raw catalog has %d books, want %d", len(catalog), expectedBookCount)
	}

	books, err := readBooks(rawRoot, catalog)
	if err != nil {
		return Manifest{}, err
	}
	return writeCorpus(outputRoot, books)
}

func readCatalog(path string) ([]rawCatalogEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read raw catalog: %w", err)
	}

	var catalog []rawCatalogEntry
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("decode raw catalog: %w", err)
	}

	seen := make(map[string]struct{}, len(catalog))
	for index, entry := range catalog {
		entry.ID = strings.ToLower(strings.TrimSpace(entry.ID))
		if entry.ID == "" || entry.ChapterCount < 1 {
			return nil, fmt.Errorf("raw catalog entry %d is invalid", index+1)
		}
		if _, exists := seen[entry.ID]; exists {
			return nil, fmt.Errorf("raw catalog has duplicate book id %q", entry.ID)
		}
		seen[entry.ID] = struct{}{}
		catalog[index] = entry
	}

	return catalog, nil
}

func readBooks(rawRoot string, catalog []rawCatalogEntry) ([]Book, error) {
	books := make([]Book, 0, len(catalog))
	for index, entry := range catalog {
		path, testament, err := rawBookPath(rawRoot, entry.ID)
		if err != nil {
			return nil, err
		}

		rawBook, err := readRawBook(path)
		if err != nil {
			return nil, fmt.Errorf("read book %q: %w", entry.ID, err)
		}
		book, err := NormalizeBook(rawBook, BookMetadata{
			ID:        entry.ID,
			Order:     index + 1,
			Testament: testament,
		})
		if err != nil {
			return nil, err
		}
		if len(book.Chapters) != entry.ChapterCount {
			return nil, fmt.Errorf("book %s has %d chapters, want %d", book.ID, len(book.Chapters), entry.ChapterCount)
		}

		books = append(books, book)
	}

	return books, nil
}

func rawBookPath(rawRoot, id string) (string, string, error) {
	oldPath := filepath.Join(rawRoot, "antigotestamento", id+".json")
	newPath := filepath.Join(rawRoot, "novotestamento", id+".json")

	oldExists := fileExists(oldPath)
	newExists := fileExists(newPath)
	switch {
	case oldExists && newExists:
		return "", "", fmt.Errorf("book %q is present in both testaments", id)
	case oldExists:
		return oldPath, "old", nil
	case newExists:
		return newPath, "new", nil
	default:
		return "", "", fmt.Errorf("raw book %q was not found", id)
	}
}

func readRawBook(path string) (RawBook, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return RawBook{}, err
	}

	var book RawBook
	if err := json.Unmarshal(data, &book); err != nil {
		return RawBook{}, err
	}
	return book, nil
}

func writeCorpus(outputRoot string, books []Book) (Manifest, error) {
	parent := filepath.Dir(outputRoot)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return Manifest{}, fmt.Errorf("create corpus parent directory: %w", err)
	}
	if fileExists(outputRoot) {
		return Manifest{}, fmt.Errorf("corpus output already exists: %s", outputRoot)
	}

	staging, err := os.MkdirTemp(parent, ".corpus-")
	if err != nil {
		return Manifest{}, fmt.Errorf("create corpus staging directory: %w", err)
	}
	defer os.RemoveAll(staging)

	if err := os.Mkdir(filepath.Join(staging, "books"), 0o755); err != nil {
		return Manifest{}, fmt.Errorf("create books directory: %w", err)
	}

	manifest := Manifest{SchemaVersion: canonicalSchemaVersion, BookCount: len(books), Books: make([]ManifestBook, 0, len(books))}
	for _, book := range books {
		data, err := marshalJSON(book)
		if err != nil {
			return Manifest{}, err
		}
		path := filepath.Join(staging, "books", fmt.Sprintf("%02d-%s.json", book.Order, book.ID))
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return Manifest{}, fmt.Errorf("write book %s: %w", book.ID, err)
		}

		chapters, verses := counts(book)
		manifest.ChapterCount += chapters
		manifest.VerseCount += verses
		manifest.NumberingGapCount += numberingGapCount(book)
		manifest.RepeatedVerseNumberCount += repeatedVerseNumberCount(book)
		manifest.Books = append(manifest.Books, ManifestBook{ID: book.ID, Order: book.Order, Chapters: chapters, Verses: verses, SHA256: sha256Hex(data)})
	}

	bibleData, err := marshalJSON(Corpus{SchemaVersion: canonicalSchemaVersion, Books: books})
	if err != nil {
		return Manifest{}, err
	}
	if err := os.WriteFile(filepath.Join(staging, "bible.json"), bibleData, 0o644); err != nil {
		return Manifest{}, fmt.Errorf("write complete corpus: %w", err)
	}
	manifest.BibleSHA256 = sha256Hex(bibleData)

	manifestData, err := marshalJSON(manifest)
	if err != nil {
		return Manifest{}, err
	}
	if err := os.WriteFile(filepath.Join(staging, "manifest.json"), manifestData, 0o644); err != nil {
		return Manifest{}, fmt.Errorf("write manifest: %w", err)
	}

	if err := os.Rename(staging, outputRoot); err != nil {
		return Manifest{}, fmt.Errorf("publish corpus: %w", err)
	}
	return manifest, nil
}

func counts(book Book) (int, int) {
	verses := 0
	for _, chapter := range book.Chapters {
		verses += len(chapter.Verses)
	}
	return len(book.Chapters), verses
}

func numberingGapCount(book Book) int {
	missing := 0
	for _, chapter := range book.Chapters {
		previous := 0
		for _, verse := range chapter.Verses {
			if verse.Number > previous {
				missing += verse.Number - previous - 1
			}
			previous = verse.Number
		}
	}
	return missing
}

func repeatedVerseNumberCount(book Book) int {
	repeated := 0
	for _, chapter := range book.Chapters {
		previous := 0
		for index, verse := range chapter.Verses {
			if index > 0 && verse.Number == previous {
				repeated++
			}
			previous = verse.Number
		}
	}
	return repeated
}

func marshalJSON(value any) ([]byte, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode canonical JSON: %w", err)
	}
	return append(data, '\n'), nil
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
