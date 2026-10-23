package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	outboxports "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type EventDAO interface {
	Insert(context.Context, EventRecord) error
}

type OutboxRepository struct {
	dao              EventDAO
	payloadProtector outboxports.PayloadProtector
}

var _ outboxports.OutboxRepository = (*OutboxRepository)(nil)

func NewOutboxRepository(dao EventDAO, payloadProtector outboxports.PayloadProtector) *OutboxRepository {
	return &OutboxRepository{dao: dao, payloadProtector: payloadProtector}
}

func (r *OutboxRepository) Publish(ctx context.Context, event *domain.Event) error {
	protectedPayload, err := r.payloadProtector.Protect(ctx, event.Payload())
	if err != nil {
		return fmt.Errorf("protecting payload: %w", err)
	}
	err = r.dao.Insert(ctx, EventRecord{ID: event.ID(), EventType: event.EventType().String(), UserID: event.UserID(), PayloadCiphertext: protectedPayload.Ciphertext, PayloadKeyVersion: protectedPayload.KeyVersion})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return application.ErrOutboxEventAlreadyExists
		}
		return fmt.Errorf("publishing outbox event: %w", err)
	}

	return nil
}
