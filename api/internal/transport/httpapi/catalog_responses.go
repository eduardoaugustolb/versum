package httpapi

import "github.com/eduardoaugustolb/versum/api/internal/catalog/domain"

type bookResponse struct {
	ID           string `json:"id"`
	Order        int    `json:"order"`
	Name         string `json:"name"`
	Testament    string `json:"testament"`
	ChapterCount int    `json:"chapter_count"`
}

type chapterResponse struct {
	BookID   string          `json:"book_id"`
	BookName string          `json:"book_name"`
	Number   int             `json:"number"`
	Verses   []verseResponse `json:"verses"`
}

type verseResponse struct {
	BookID  string `json:"book_id"`
	Chapter int    `json:"chapter"`
	Number  int    `json:"number"`
	Text    string `json:"text"`
	Part    int    `json:"part"`
}

func newBookResponse(book domain.Book) bookResponse {
	return bookResponse{ID: book.ID(), Order: book.Order(), Name: book.Name(), Testament: string(book.Testament()), ChapterCount: book.ChapterCount()}
}

func newChapterResponse(chapter domain.Chapter) chapterResponse {
	verses := chapter.Verses()
	response := chapterResponse{BookID: chapter.BookID(), BookName: chapter.BookName(), Number: chapter.Number(), Verses: make([]verseResponse, 0, len(verses))}
	for _, verse := range verses {
		response.Verses = append(response.Verses, verseResponse{BookID: verse.BookID(), Chapter: verse.Chapter(), Number: verse.Number(), Text: verse.Text(), Part: verse.Part()})
	}
	return response
}
