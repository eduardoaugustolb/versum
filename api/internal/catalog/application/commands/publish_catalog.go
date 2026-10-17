package commands

import (
	"context"
	"fmt"
	"strings"

	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
)

type PublishCatalog struct{ transactions ports.CatalogUnitOfWork }

func NewPublishCatalog(transactions ports.CatalogUnitOfWork) PublishCatalog {
	return PublishCatalog{transactions}
}
func (uc PublishCatalog) Execute(ctx context.Context, books []PublishBookInput, corpusSHA256 string) error {
	if strings.TrimSpace(corpusSHA256) == "" {
		return fmt.Errorf("empty corpus hash")
	}
	return uc.transactions.WithinCatalogTransaction(ctx, func(repository ports.CatalogRepository, versions ports.CatalogVersionRepository) error {
		for _, book := range books {
			entity, verses, err := buildPublication(book)
			if err != nil {
				return fmt.Errorf("building book %q: %w", book.ID, err)
			}
			if err := repository.ReplaceBook(ctx, entity, verses); err != nil {
				return fmt.Errorf("publishing book %q: %w", book.ID, err)
			}
		}
		return versions.Record(ctx, corpusSHA256)
	})
}
