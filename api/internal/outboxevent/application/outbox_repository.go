package application

import (
	"context"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type OutboxRepository interface {
	Publish(ctx context.Context, event *domain.Event) error
}

// ClaimedEvent contains a decrypted event and its exclusive delivery lease.
// Attempt counts claims, including claims recovered after a process crash.
type ClaimedEvent struct {
	Event      *domain.Event
	LeaseToken string
	Attempt    int
}

// EventDeliveryRepository is separate from transactional publication so that
// producers do not depend on worker operations.
type EventDeliveryRepository interface {
	// Claim atomically acquires one due event or an event with an expired lease.
	// It increments Attempt and issues a new, unique LeaseToken. No work returns
	// ErrNoEventAvailable. It must commit before returning, without holding locks
	// during Handle. now and leaseUntil are UTC instants.
	Claim(ctx context.Context, now, leaseUntil time.Time) (ClaimedEvent, error)
	// Completion, retry and failure must conditionally update only an unexpired,
	// matching lease; otherwise return ErrEventLeaseLost. Completion and terminal
	// failure must discard the secret payload. Never persist raw handler errors.
	Complete(ctx context.Context, eventID, leaseToken string, now time.Time) error
	Retry(ctx context.Context, eventID, leaseToken string, now, availableAt time.Time) error
	Fail(ctx context.Context, eventID, leaseToken string, now time.Time) error
}
