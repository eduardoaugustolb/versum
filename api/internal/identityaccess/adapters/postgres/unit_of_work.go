package postgres

import (
	"context"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	outboxPg "github.com/eduardoaugustolb/versum/api/internal/outboxevent/adapters/postgres"
	outboxPorts "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
)

type UnitOfWork struct {
	DbExecutor             dbexec.Executor
	EmailProtector         ports.EmailProtector
	OutboxPayloadProtector outboxPorts.PayloadProtector
}

var _ ports.IdentityAccessUnitOfWork = &UnitOfWork{}

func NewUnitOfWork(dbExecutor dbexec.Executor, emailProtector ports.EmailProtector, outboxPayloadProtector outboxPorts.PayloadProtector) *UnitOfWork {
	return &UnitOfWork{
		DbExecutor:             dbExecutor,
		EmailProtector:         emailProtector,
		OutboxPayloadProtector: outboxPayloadProtector,
	}
}

func (u *UnitOfWork) WithinTransaction(ctx context.Context, f func(repositories ports.IdentityAccessUnitOfWorkRepositories) error) error {
	tx, err := u.DbExecutor.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	repositories := ports.IdentityAccessUnitOfWorkRepositories{
		Users:       NewUserRepository(tx, u.EmailProtector),
		LoginTokens: NewLoginTokenRepository(tx),
		Sessions:    NewSessionRepository(tx),
		Outbox:      outboxPg.NewOutboxRepository(tx, u.OutboxPayloadProtector),
	}
	err = f(repositories)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
