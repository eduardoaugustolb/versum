package postgres

import (
	"context"

	ports "github.com/eduardoaugustolb/versum/api/internal/catalog/application"
	adapterpostgres "github.com/eduardoaugustolb/versum/api/internal/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct{ pool *pgxpool.Pool }

func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return TransactionManager{pool: pool}
}

func (m TransactionManager) WithinTransaction(ctx context.Context, fn func(context.Context, ports.CatalogRepository) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(ctx, NewRepository(adapterpostgres.NewPgxExecutor(tx))); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
