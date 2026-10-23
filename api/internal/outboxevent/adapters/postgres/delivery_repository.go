package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/id"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/application"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
	"github.com/jackc/pgx/v5"
)

type DeliveryDAO interface {
	Claim(context.Context, time.Time, string, time.Time) (DeliveryRecord, error)
	Complete(context.Context, string, string, time.Time) error
	Retry(context.Context, string, string, time.Time, time.Time) error
	Fail(context.Context, string, string, time.Time) error
}

type EventDeliveryRepository struct {
	dao       DeliveryDAO
	protector application.PayloadProtector
	ids       id.Generator
}

var _ application.EventDeliveryRepository = (*EventDeliveryRepository)(nil)

func NewEventDeliveryRepository(dao DeliveryDAO, protector application.PayloadProtector, ids id.Generator) *EventDeliveryRepository {
	return &EventDeliveryRepository{dao: dao, protector: protector, ids: ids}
}

func (r *EventDeliveryRepository) Claim(ctx context.Context, now, leaseUntil time.Time) (application.ClaimedEvent, error) {
	lease := r.ids.Generate()
	record, err := r.dao.Claim(ctx, now, lease, leaseUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		return application.ClaimedEvent{}, application.ErrNoEventAvailable
	}
	if err != nil {
		return application.ClaimedEvent{}, err
	}
	payload, err := r.protector.Unprotect(ctx, application.ProtectedPayload{Ciphertext: record.PayloadCiphertext, KeyVersion: record.PayloadKeyVersion})
	if err != nil {
		// Preserve ciphertext and let the lease expire: missing historical keys can
		// be repaired operationally. Never log ciphertext or decrypted payloads.
		return application.ClaimedEvent{}, errors.New("cannot decrypt claimed outbox event")
	}
	event, err := domain.RehydrateEvent(record.ID, record.EventType, record.UserID, payload)
	if err != nil {
		return application.ClaimedEvent{}, err
	}
	return application.ClaimedEvent{Event: event, LeaseToken: lease, Attempt: record.Attempt}, nil
}

func deliveryResult(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return application.ErrEventLeaseLost
	}
	return err
}
func (r *EventDeliveryRepository) Complete(ctx context.Context, eventID, lease string, now time.Time) error {
	return deliveryResult(r.dao.Complete(ctx, eventID, lease, now))
}
func (r *EventDeliveryRepository) Retry(ctx context.Context, eventID, lease string, now, availableAt time.Time) error {
	return deliveryResult(r.dao.Retry(ctx, eventID, lease, now, availableAt))
}
func (r *EventDeliveryRepository) Fail(ctx context.Context, eventID, lease string, now time.Time) error {
	return deliveryResult(r.dao.Fail(ctx, eventID, lease, now))
}
