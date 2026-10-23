package application

import (
	"context"
	"errors"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/clock"
	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type EventProcessorOptions struct {
	LeaseDuration time.Duration
	HandleTimeout time.Duration
	RetryDelay    time.Duration
	MaxAttempts   int
}

type EventProcessor struct {
	repository EventDeliveryRepository
	clock      clock.Clock
	handlers   map[domain.EventType]EventHandler
	options    EventProcessorOptions
}

func NewEventProcessor(repository EventDeliveryRepository, clock clock.Clock, handlers map[domain.EventType]EventHandler, options EventProcessorOptions) (*EventProcessor, error) {
	if repository == nil || clock == nil || len(handlers) == 0 ||
		options.HandleTimeout <= 0 || options.LeaseDuration <= options.HandleTimeout ||
		options.RetryDelay <= 0 || options.MaxAttempts <= 0 {
		return nil, ErrInvalidEventProcessor
	}
	registered := make(map[domain.EventType]EventHandler, len(handlers))
	for eventType, handler := range handlers {
		if _, err := domain.ParseEventType(eventType.String()); err != nil || handler == nil {
			return nil, ErrInvalidEventProcessor
		}
		registered[eventType] = handler
	}
	return &EventProcessor{repository: repository, clock: clock, handlers: registered, options: options}, nil
}

// ProcessNext handles at most one event. false, nil means the outbox is idle.
// A handled failure is rescheduled (or marked terminal) before being returned.
// On cancellation, leave the lease to expire for recovery by another worker.
// Errors may wrap handler errors: callers must sanitize them before logging.
func (p *EventProcessor) ProcessNext(ctx context.Context) (bool, error) {
	now := p.clock.Now().UTC()
	claim, err := p.repository.Claim(ctx, now, now.Add(p.options.LeaseDuration))
	if errors.Is(err, ErrNoEventAvailable) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if claim.Event == nil || claim.LeaseToken == "" || claim.Attempt < 1 {
		return true, ErrInvalidEventClaim
	}
	eventID := claim.Event.ID()
	// A crash on the final attempt must not create an unlimited retry cycle.
	if claim.Attempt > p.options.MaxAttempts {
		return true, p.repository.Fail(ctx, eventID, claim.LeaseToken, p.clock.Now().UTC())
	}
	handler, ok := p.handlers[claim.Event.EventType()]
	if !ok {
		return true, errors.Join(ErrEventHandlerNotFound, p.repository.Fail(ctx, eventID, claim.LeaseToken, p.clock.Now().UTC()))
	}
	handleCtx, cancel := context.WithTimeout(ctx, p.options.HandleTimeout)
	err = handler.Handle(handleCtx, claim.Event)
	if err == nil {
		err = handleCtx.Err()
	}
	cancel()
	if ctx.Err() != nil {
		return true, ctx.Err()
	}
	now = p.clock.Now().UTC()
	if err == nil {
		return true, p.repository.Complete(ctx, eventID, claim.LeaseToken, now)
	}
	if claim.Attempt >= p.options.MaxAttempts {
		return true, errors.Join(err, p.repository.Fail(ctx, eventID, claim.LeaseToken, now))
	}
	return true, errors.Join(err, p.repository.Retry(ctx, eventID, claim.LeaseToken, now, now.Add(p.options.RetryDelay)))
}
