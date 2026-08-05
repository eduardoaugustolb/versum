package catalog

import "context"

type ListBooks struct {
	repo BookRepository
}

func NewListBooks(repo BookRepository) *ListBooks {
	return &ListBooks{repo: repo}
}

func (uc *ListBooks) Execute(ctx context.Context) ([]Book, error) {
	return uc.repo.ListBooks(ctx)
}
