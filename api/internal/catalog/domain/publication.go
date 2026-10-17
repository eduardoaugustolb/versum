package domain

import "fmt"

// ValidateBookVerses checks the composition before a publication reaches storage.
func ValidateBookVerses(book Book, verses []Verse) error {
	if _, err := NewBook(NewBookParams{book.ID(), book.Order(), book.Name(), book.Testament(), book.ChapterCount()}); err != nil {
		return err
	}
	seen := make(map[[3]int]struct{}, len(verses))
	for _, v := range verses {
		if _, err := NewVerse(NewVerseParams{v.BookID(), v.Chapter(), v.Number(), v.Text(), v.Part()}); err != nil {
			return err
		}
		if v.BookID() != book.ID() || v.Chapter() > book.ChapterCount() {
			return fmt.Errorf("verse does not belong to book: %w", ErrInvalidVerse)
		}
		key := [3]int{v.Chapter(), v.Number(), v.Part()}
		if _, exists := seen[key]; exists {
			return fmt.Errorf("duplicate verse: %w", ErrInvalidVerse)
		}
		seen[key] = struct{}{}
	}
	return nil
}
