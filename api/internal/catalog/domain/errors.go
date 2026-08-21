package domain

import "errors"

var (
	ErrInvalidBookID       = errors.New("invalid book id")
	ErrInvalidBookName     = errors.New("invalid book name")
	ErrInvalidBookOrder    = errors.New("invalid book order")
	ErrInvalidTestament    = errors.New("invalid testament")
	ErrInvalidChapterCount = errors.New("invalid chapter count")
	ErrInvalidVerse        = errors.New("invalid verse")
	ErrChapterNotFound     = errors.New("chapter not found")
)
