package domain

type Verse struct {
	bookID  string
	chapter int
	number  int
	text    string
	part    int
}

type NewVerseParams struct {
	BookID  string
	Chapter int
	Number  int
	Text    string
	Part    int
}

func NewVerse(params NewVerseParams) (Verse, error) {
	if params.BookID == "" || params.Chapter <= 0 || params.Number <= 0 || params.Text == "" || params.Part <= 0 {
		return Verse{}, ErrInvalidVerse
	}
	return Verse{bookID: params.BookID, chapter: params.Chapter, number: params.Number, text: params.Text, part: params.Part}, nil
}

func (v Verse) BookID() string { return v.bookID }
func (v Verse) Chapter() int   { return v.chapter }
func (v Verse) Number() int    { return v.number }
func (v Verse) Text() string   { return v.text }
func (v Verse) Part() int      { return v.part }
