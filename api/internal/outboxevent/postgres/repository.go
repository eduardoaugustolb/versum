package postgres

import (
	"context"
	"fmt"

	outboxports "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
	dbexecports "github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
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
		return fmt.Errorf("publishing outbox event: %w", err)
	}

	return nil
}
