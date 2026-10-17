package domain_test

import (
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"testing"
)

func TestPublicationRejectsInconsistentVerses(t *testing.T) {
	book, err := domain.NewBook(domain.NewBookParams{ID: "gn", Order: 1, Name: "Gênesis", Testament: domain.TestamentOld, ChapterCount: 1})
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name, book string
		chapter    int
	}{{"foreign book", "ex", 1}, {"chapter outside book", "gn", 2}} {
		t.Run(test.name, func(t *testing.T) {
			verse, err := domain.NewVerse(domain.NewVerseParams{BookID: test.book, Chapter: test.chapter, Number: 1, Text: "text", Part: 1})
			if err != nil {
				t.Fatal(err)
			}
			if err := domain.ValidateBookVerses(book, []domain.Verse{verse}); err == nil {
				t.Fatal("accepted inconsistent composition")
			}
		})
	}
	verse, _ := domain.NewVerse(domain.NewVerseParams{BookID: "gn", Chapter: 1, Number: 1, Text: "text", Part: 1})
	if err := domain.ValidateBookVerses(book, []domain.Verse{verse, verse}); err == nil {
		t.Fatal("accepted duplicate verses")
	}
}
func TestChapterValidatesCompositionAndOwnsVerses(t *testing.T) {
	verse, err := domain.NewVerse(domain.NewVerseParams{BookID: "gn", Chapter: 1, Number: 1, Text: "text", Part: 1})
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		book, name string
		number     int
		verses     []domain.Verse
	}{
		{"", "Gênesis", 1, nil}, {"gn", " ", 1, nil}, {"gn", "Gênesis", 0, nil}, {"ex", "Êxodo", 1, []domain.Verse{verse}}, {"gn", "Gênesis", 2, []domain.Verse{verse}}, {"gn", "Gênesis", 1, []domain.Verse{{}}}, {"gn", "Gênesis", 1, []domain.Verse{verse, verse}},
	} {
		if _, err := domain.NewChapter(test.book, test.name, test.number, test.verses); err == nil {
			t.Fatalf("accepted invalid chapter: %+v", test)
		}
	}
	input := []domain.Verse{verse}
	chapter, err := domain.NewChapter("gn", "Gênesis", 1, input)
	if err != nil {
		t.Fatal(err)
	}
	input[0] = domain.Verse{}
	output := chapter.Verses()
	output[0] = domain.Verse{}
	if chapter.Verses()[0].Text() != "text" {
		t.Fatal("chapter backing array escaped")
	}
}
