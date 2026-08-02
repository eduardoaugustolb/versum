package biblecorpus

import (
	"fmt"
	"regexp"
	"strings"
)

var leadingVerseMarker = regexp.MustCompile(`^\[\d+\]\s*`)

type RawBook struct {
	Name     string       `json:"livro"`
	Chapters []RawChapter `json:"capitulos"`
}

type RawChapter struct {
	Number int        `json:"capitulo"`
	Verses []RawVerse `json:"versiculos"`
}

type RawVerse struct {
	Number int    `json:"numero"`
	Text   string `json:"texto"`
}

type BookMetadata struct {
	ID        string
	Order     int
	Testament string
}

type Book struct {
	ID        string    `json:"id"`
	Order     int       `json:"order"`
	Name      string    `json:"name"`
	Testament string    `json:"testament"`
	Chapters  []Chapter `json:"chapters"`
}

type Chapter struct {
	Number int     `json:"number"`
	Verses []Verse `json:"verses"`
}

type Verse struct {
	Number int    `json:"number"`
	Part   int    `json:"part"`
	Text   string `json:"text"`
}

func NormalizeBook(raw RawBook, metadata BookMetadata) (Book, error) {
	if strings.TrimSpace(raw.Name) == "" {
		return Book{}, fmt.Errorf("book %d: name is required", metadata.Order)
	}

	book := Book{
		ID:        metadata.ID,
		Order:     metadata.Order,
		Name:      strings.TrimSpace(raw.Name),
		Testament: metadata.Testament,
		Chapters:  make([]Chapter, 0, len(raw.Chapters)),
	}

	for chapterIndex, rawChapter := range raw.Chapters {
		expectedChapter := chapterIndex + 1
		if rawChapter.Number != expectedChapter {
			return Book{}, fmt.Errorf("book %s: chapter %d, want %d", book.ID, rawChapter.Number, expectedChapter)
		}

		chapter := Chapter{Number: rawChapter.Number, Verses: make([]Verse, 0, len(rawChapter.Verses))}
		previousVerseNumber := 0
		part := 0
		for verseIndex, rawVerse := range rawChapter.Verses {
			if rawVerse.Number < 1 {
				return Book{}, fmt.Errorf("book %s chapter %d: verse number must be positive", book.ID, chapter.Number)
			}
			if verseIndex > 0 && rawVerse.Number < previousVerseNumber {
				return Book{}, fmt.Errorf("book %s chapter %d: verse %d must not be less than %d", book.ID, chapter.Number, rawVerse.Number, previousVerseNumber)
			}
			if rawVerse.Number == previousVerseNumber {
				part++
			} else {
				part = 1
			}

			text := strings.TrimSpace(leadingVerseMarker.ReplaceAllString(strings.TrimSpace(rawVerse.Text), ""))
			if text == "" {
				return Book{}, fmt.Errorf("book %s chapter %d verse %d: text is required", book.ID, chapter.Number, rawVerse.Number)
			}

			chapter.Verses = append(chapter.Verses, Verse{Number: rawVerse.Number, Part: part, Text: text})
			previousVerseNumber = rawVerse.Number
		}

		book.Chapters = append(book.Chapters, chapter)
	}

	return book, nil
}
