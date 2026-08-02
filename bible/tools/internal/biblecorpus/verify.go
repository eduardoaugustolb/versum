package biblecorpus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func Verify(root string) error {
	manifest, err := readManifest(filepath.Join(root, "manifest.json"))
	if err != nil {
		return err
	}
	if err := validateManifest(manifest); err != nil {
		return err
	}

	bibleData, err := os.ReadFile(filepath.Join(root, "bible.json"))
	if err != nil {
		return fmt.Errorf("read complete corpus: %w", err)
	}
	if actual := sha256Hex(bibleData); actual != manifest.BibleSHA256 {
		return fmt.Errorf("complete corpus hash mismatch: got %s", actual)
	}

	var corpus Corpus
	if err := json.Unmarshal(bibleData, &corpus); err != nil {
		return fmt.Errorf("decode complete corpus: %w", err)
	}
	if corpus.SchemaVersion != canonicalSchemaVersion {
		return fmt.Errorf("complete corpus schema version is %d, want %d", corpus.SchemaVersion, canonicalSchemaVersion)
	}
	if len(corpus.Books) != expectedBookCount {
		return fmt.Errorf("complete corpus has %d books, want %d", len(corpus.Books), expectedBookCount)
	}

	entries, err := os.ReadDir(filepath.Join(root, "books"))
	if err != nil {
		return fmt.Errorf("read books directory: %w", err)
	}
	if len(entries) != expectedBookCount {
		return fmt.Errorf("books directory has %d entries, want %d", len(entries), expectedBookCount)
	}

	totalChapters := 0
	totalVerses := 0
	gapCount := 0
	repeatedCount := 0
	for index, details := range manifest.Books {
		book := corpus.Books[index]
		if err := verifyBook(root, book, details, index+1); err != nil {
			return err
		}

		chapters, verses := counts(book)
		totalChapters += chapters
		totalVerses += verses
		gapCount += numberingGapCount(book)
		repeatedCount += repeatedVerseNumberCount(book)
	}

	if manifest.ChapterCount != totalChapters || manifest.VerseCount != totalVerses {
		return fmt.Errorf("manifest totals do not match corpus")
	}
	if manifest.NumberingGapCount != gapCount || manifest.RepeatedVerseNumberCount != repeatedCount {
		return fmt.Errorf("manifest numbering audit does not match corpus")
	}

	return nil
}

func readManifest(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	return manifest, nil
}

func validateManifest(manifest Manifest) error {
	if manifest.SchemaVersion != canonicalSchemaVersion {
		return fmt.Errorf("manifest schema version is %d, want %d", manifest.SchemaVersion, canonicalSchemaVersion)
	}
	if manifest.BookCount != expectedBookCount || len(manifest.Books) != expectedBookCount {
		return fmt.Errorf("manifest has %d books and %d entries, want %d", manifest.BookCount, len(manifest.Books), expectedBookCount)
	}
	if len(manifest.BibleSHA256) != 64 {
		return fmt.Errorf("manifest complete corpus hash is invalid")
	}
	return nil
}

func verifyBook(root string, book Book, details ManifestBook, expectedOrder int) error {
	if book.Order != expectedOrder || details.Order != expectedOrder {
		return fmt.Errorf("book position %d has inconsistent order", expectedOrder)
	}
	if book.ID != details.ID {
		return fmt.Errorf("book position %d has inconsistent id", expectedOrder)
	}
	if err := validateBook(book); err != nil {
		return err
	}

	path := filepath.Join(root, "books", fmt.Sprintf("%02d-%s.json", details.Order, details.ID))
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read book %s: %w", details.ID, err)
	}
	if actual := sha256Hex(data); actual != details.SHA256 {
		return fmt.Errorf("book %s hash mismatch: got %s", details.ID, actual)
	}

	var fromFile Book
	if err := json.Unmarshal(data, &fromFile); err != nil {
		return fmt.Errorf("decode book %s: %w", details.ID, err)
	}
	if !reflect.DeepEqual(fromFile, book) {
		return fmt.Errorf("book %s does not match complete corpus", details.ID)
	}

	chapters, verses := counts(book)
	if details.Chapters != chapters || details.Verses != verses {
		return fmt.Errorf("book %s totals do not match manifest", details.ID)
	}
	return nil
}

func validateBook(book Book) error {
	if strings.TrimSpace(book.ID) == "" || strings.TrimSpace(book.Name) == "" {
		return fmt.Errorf("book %d has required metadata missing", book.Order)
	}
	if book.Testament != "old" && book.Testament != "new" {
		return fmt.Errorf("book %s has invalid testament", book.ID)
	}
	if len(book.Chapters) == 0 {
		return fmt.Errorf("book %s has no chapters", book.ID)
	}

	for chapterIndex, chapter := range book.Chapters {
		if chapter.Number != chapterIndex+1 {
			return fmt.Errorf("book %s has invalid chapter order", book.ID)
		}
		if len(chapter.Verses) == 0 {
			return fmt.Errorf("book %s chapter %d has no verses", book.ID, chapter.Number)
		}

		previousNumber := 0
		expectedPart := 0
		for verseIndex, verse := range chapter.Verses {
			if verse.Number < 1 || (verseIndex > 0 && verse.Number < previousNumber) {
				return fmt.Errorf("book %s chapter %d has invalid verse order", book.ID, chapter.Number)
			}
			if verse.Number == previousNumber {
				expectedPart++
			} else {
				expectedPart = 1
			}
			if verse.Part != expectedPart {
				return fmt.Errorf("book %s chapter %d verse %d has invalid part", book.ID, chapter.Number, verse.Number)
			}
			if verse.Text == "" || verse.Text != strings.TrimSpace(verse.Text) || leadingVerseMarker.MatchString(verse.Text) {
				return fmt.Errorf("book %s chapter %d verse %d has invalid text", book.ID, chapter.Number, verse.Number)
			}
			previousNumber = verse.Number
		}
	}

	return nil
}
