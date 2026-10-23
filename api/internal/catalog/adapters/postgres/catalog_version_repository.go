package postgres

import (
	"context"
	"fmt"

	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
)

type CatalogVersionDAO interface {
	RecordVersion(context.Context, string) error
}

type CatalogVersionRepository struct{ dao CatalogVersionDAO }

var _ ports.CatalogVersionRepository = (*CatalogVersionRepository)(nil)

func NewCatalogVersionRepository(dao CatalogVersionDAO) *CatalogVersionRepository {
	return &CatalogVersionRepository{
		dao: dao,
	}
}

func (r *CatalogVersionRepository) Record(ctx context.Context, corpusSHA256 string) error {
	if err := r.dao.RecordVersion(ctx, corpusSHA256); err != nil {
		return fmt.Errorf("upserting catalog version: %w", err)
	}

	return nil
}
