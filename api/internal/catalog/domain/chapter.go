package domain

import "strings"

type Chapter struct {
	bookID   string
	bookName string
	number   int
	verses   []Verse
}

func NewChapter(bookID, bookName string, number int, verses []Verse) (Chapter, error) {
	if strings.TrimSpace(bookID) == "" {
		return Chapter{}, ErrInvalidBookID
	}
	if strings.TrimSpace(bookName) == "" {
		return Chapter{}, ErrInvalidBookName
	}
	if number <= 0 {
		return Chapter{}, ErrInvalidChapter
	}
	seen := make(map[[2]int]struct{}, len(verses))
	for _, verse := range verses {
		if _, err := NewVerse(NewVerseParams{verse.BookID(), verse.Chapter(), verse.Number(), verse.Text(), verse.Part()}); err != nil {
			return Chapter{}, err
		}
		if verse.BookID() != bookID || verse.Chapter() != number {
			return Chapter{}, ErrInvalidChapter
		}
		key := [2]int{verse.Number(), verse.Part()}
		if _, ok := seen[key]; ok {
			return Chapter{}, ErrInvalidVerse
		}
		seen[key] = struct{}{}
	}
	return Chapter{bookID: bookID, bookName: bookName, number: number, verses: append([]Verse(nil), verses...)}, nil
}

func (c Chapter) BookID() string   { return c.bookID }
func (c Chapter) BookName() string { return c.bookName }
func (c Chapter) Number() int      { return c.number }
func (c Chapter) Verses() []Verse  { return append([]Verse(nil), c.verses...) }
