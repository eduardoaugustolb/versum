package catalog

import (
	"context"

	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
)

type Repository struct {
	db dbexec.Executor
}

func NewRepository(db dbexec.Executor) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListBooks(ctx context.Context) ([]Book, error) {
	books := []Book{}
	rows, err := r.db.Query(ctx, ListBooksQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Order, &book.Name, &book.Testament, &book.ChapterCount); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *Repository) FindChapter(ctx context.Context, bookID string, number int) (Chapter, error) {
	chapter := Chapter{
		BookID: bookID,
		Number: number,
		Verses: []Verse{},
	}
	rows, err := r.db.Query(ctx, FindChapterVersesQuery, bookID, number)
	if err != nil {
		return Chapter{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var verse Verse
		var bookName string
		if err := rows.Scan(&verse.BookID, &verse.Chapter, &verse.Number, &verse.Text, &verse.Part, &bookName); err != nil {
			return Chapter{}, err
		}
		chapter.BookName = bookName
		chapter.Verses = append(chapter.Verses, verse)
	}

	if err := rows.Err(); err != nil {
		return Chapter{}, err
	}

	if len(chapter.Verses) == 0 {
		return Chapter{}, ErrChapterNotFound
	}

	return chapter, nil
}
