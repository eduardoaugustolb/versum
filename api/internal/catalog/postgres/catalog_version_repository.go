package postgres

import (
	"context"
	"fmt"

	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
)

type CatalogVersionRepository struct {
	db dbexec.Executor
}

func NewCatalogVersionRepository(db dbexec.Executor) *CatalogVersionRepository {
	return &CatalogVersionRepository{
		db: db,
	}
}

func (r *CatalogVersionRepository) Record(ctx context.Context, corpusSHA256 string) error {
	if err := r.db.Exec(ctx, UpsertCatalogVersionQuery, corpusSHA256); err != nil {
		return fmt.Errorf("upserting catalog version: %w", err)
	}

	return nil
}
