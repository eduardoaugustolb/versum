package postgres

import (
	"context"

	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
)

type TransactionManager struct{ db dbexec.Executor }

func NewTransactionManager(db dbexec.Executor) TransactionManager {
	return TransactionManager{db: db}
}

func (m TransactionManager) WithinTransaction(ctx context.Context, fn func(context.Context, ports.CatalogRepository) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := fn(ctx, NewRepository(NewDAO(tx))); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m TransactionManager) WithinCatalogTransaction(ctx context.Context, fn func(ports.CatalogRepository, ports.CatalogVersionRepository) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := fn(NewRepository(NewDAO(tx)), NewCatalogVersionRepository(NewDAO(tx))); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
