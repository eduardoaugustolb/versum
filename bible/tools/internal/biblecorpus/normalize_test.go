package biblecorpus

import "testing"

func TestNormalizeBookRemovesOnlyTheLeadingVerseMarker(t *testing.T) {
	raw := RawBook{
		Name: "Gênesis",
		Chapters: []RawChapter{{
			Number: 1,
			Verses: []RawVerse{{Number: 1, Text: " [1] No princípio. "}},
		}},
	}

	book, err := NormalizeBook(raw, BookMetadata{
		ID:        "gn",
		Order:     1,
		Testament: "old",
	})
	if err != nil {
		t.Fatalf("NormalizeBook() error = %v", err)
	}

	if got, want := book.Chapters[0].Verses[0].Text, "No princípio."; got != want {
		t.Fatalf("normalized text = %q, want %q", got, want)
	}
}

func TestNormalizeBookPreservesGapsInSourceVerseNumbering(t *testing.T) {
	raw := RawBook{
		Name: "Gênesis",
		Chapters: []RawChapter{{
			Number: 1,
			Verses: []RawVerse{
				{Number: 1, Text: "[1] Primeiro."},
				{Number: 3, Text: "[3] Terceiro."},
			},
		}},
	}

	book, err := NormalizeBook(raw, BookMetadata{ID: "gn", Order: 1, Testament: "old"})
	if err != nil {
		t.Fatalf("NormalizeBook() error = %v", err)
	}
	if got, want := book.Chapters[0].Verses[1].Number, 3; got != want {
		t.Fatalf("second verse number = %d, want %d", got, want)
	}
}

func TestNumberingGapCountCountsOnlyMissingRanges(t *testing.T) {
	book := Book{Chapters: []Chapter{{
		Number: 1,
		Verses: []Verse{
			{Number: 1},
			{Number: 3},
			{Number: 4},
			{Number: 7},
		},
	}}}

	if got, want := numberingGapCount(book), 3; got != want {
		t.Fatalf("numberingGapCount() = %d, want %d", got, want)
	}
}

func TestNormalizeBookPreservesRepeatedVerseNumbersAsParts(t *testing.T) {
	raw := RawBook{
		Name: "Josué",
		Chapters: []RawChapter{{
			Number: 1,
			Verses: []RawVerse{
				{Number: 9, Text: "[9] Primeira variante."},
				{Number: 9, Text: "[9] Segunda variante."},
			},
		}},
	}

	book, err := NormalizeBook(raw, BookMetadata{ID: "js", Order: 6, Testament: "old"})
	if err != nil {
		t.Fatalf("NormalizeBook() error = %v", err)
	}
	if got, want := book.Chapters[0].Verses[1].Part, 2; got != want {
		t.Fatalf("second verse part = %d, want %d", got, want)
	}
}
