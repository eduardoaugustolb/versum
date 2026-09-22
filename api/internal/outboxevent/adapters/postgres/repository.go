package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	dbexecports "github.com/eduardoaugustolb/versum/api/internal/database"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	outboxports "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type OutboxRepository struct {
	dbExecutor       dbexecports.Executor
	payloadProtector outboxports.PayloadProtector
}

var _ outboxports.OutboxRepository = (*OutboxRepository)(nil)

func NewOutboxRepository(dbExecutor dbexecports.Executor, payloadProtector outboxports.PayloadProtector) *OutboxRepository {
	return &OutboxRepository{dbExecutor: dbExecutor, payloadProtector: payloadProtector}
}

func (r *OutboxRepository) Publish(ctx context.Context, event *domain.Event) error {
	protectedPayload, err := r.payloadProtector.Protect(ctx, event.Payload())
	if err != nil {
		return fmt.Errorf("protecting payload: %w", err)
	}
	err = r.dbExecutor.Exec(
		ctx,
		PublishEventQuery,
		event.ID(),
		event.EventType().String(),
		protectedPayload.Ciphertext,
		protectedPayload.KeyVersion,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return application.ErrOutboxEventAlreadyExists
		}
		return fmt.Errorf("publishing outbox event: %w", err)
	}

	return nil
}
