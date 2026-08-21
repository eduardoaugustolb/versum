package domain

type Chapter struct {
	bookID   string
	bookName string
	number   int
	verses   []Verse
}

func NewChapter(bookID, bookName string, number int, verses []Verse) Chapter {
	return Chapter{bookID: bookID, bookName: bookName, number: number, verses: verses}
}

func (c Chapter) BookID() string   { return c.bookID }
func (c Chapter) BookName() string { return c.bookName }
func (c Chapter) Number() int      { return c.number }
func (c Chapter) Verses() []Verse  { return append([]Verse(nil), c.verses...) }
