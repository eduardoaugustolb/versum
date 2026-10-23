package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/outboxevent/domain"
)

type deliveryStub struct {
	claim       ClaimedEvent
	claimErr    error
	finishErr   error
	action      string
	availableAt time.Time
	id, lease   string
}

func (s *deliveryStub) Claim(context.Context, time.Time, time.Time) (ClaimedEvent, error) {
	return s.claim, s.claimErr
}
func (s *deliveryStub) Complete(_ context.Context, id, lease string, _ time.Time) error {
	s.action, s.id, s.lease = "complete", id, lease
	return s.finishErr
}
func (s *deliveryStub) Retry(_ context.Context, id, lease string, _, availableAt time.Time) error {
	s.action, s.id, s.lease, s.availableAt = "retry", id, lease, availableAt
	return s.finishErr
}
func (s *deliveryStub) Fail(_ context.Context, id, lease string, _ time.Time) error {
	s.action, s.id, s.lease = "fail", id, lease
	return s.finishErr
}

type processorClock struct{ now time.Time }

func (c processorClock) Now() time.Time { return c.now }

type handlerFunc func(context.Context, *domain.Event) error

func (f handlerFunc) Handle(ctx context.Context, event *domain.Event) error { return f(ctx, event) }

func TestEventProcessorDeliveryOutcomes(t *testing.T) {
	handlerErr := errors.New("delivery failed")
	for _, tt := range []struct {
		name                  string
		attempt               int
		missing               bool
		handlerErr, finishErr error
		action                string
		calls                 int
	}{
		{name: "success", attempt: 1, action: "complete", calls: 1},
		{name: "retry", attempt: 1, handlerErr: handlerErr, action: "retry", calls: 1},
		{name: "final attempt", attempt: 3, handlerErr: handlerErr, action: "fail", calls: 1},
		{name: "crash exhausted attempts", attempt: 4, action: "fail"},
		{name: "unknown type", attempt: 1, missing: true, action: "fail"},
		{name: "stale completion", attempt: 1, finishErr: ErrEventLeaseLost, action: "complete", calls: 1},
		{name: "retry fails", attempt: 1, handlerErr: handlerErr, finishErr: ErrEventLeaseLost, action: "retry", calls: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			event, err := domain.NewEvent("event-1", "mail.requested", "", []byte("payload"))
			if err != nil {
				t.Fatal(err)
			}
			repo := &deliveryStub{claim: ClaimedEvent{Event: event, LeaseToken: "lease-1", Attempt: tt.attempt}, finishErr: tt.finishErr}
			calls := 0
			handlers := map[domain.EventType]EventHandler{"mail.requested": handlerFunc(func(ctx context.Context, _ *domain.Event) error {
				calls++
				if _, ok := ctx.Deadline(); !ok {
					t.Error("missing handler deadline")
				}
				return tt.handlerErr
			})}
			if tt.missing {
				handlers["other"] = handlers["mail.requested"]
				delete(handlers, "mail.requested")
			}
			now := time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)
			processor, err := NewEventProcessor(repo, processorClock{now}, handlers, EventProcessorOptions{LeaseDuration: time.Minute, HandleTimeout: 30 * time.Second, RetryDelay: time.Second, MaxAttempts: 3})
			if err != nil {
				t.Fatal(err)
			}
			processed, err := processor.ProcessNext(context.Background())
			if !processed || repo.action != tt.action || calls != tt.calls || repo.id != event.ID() || repo.lease != "lease-1" {
				t.Fatalf("unexpected outcome: processed=%v repo=%+v calls=%d", processed, repo, calls)
			}
			for _, want := range []error{tt.handlerErr, tt.finishErr} {
				if want != nil && !errors.Is(err, want) {
					t.Fatalf("expected %v, got %v", want, err)
				}
			}
			if tt.missing && !errors.Is(err, ErrEventHandlerNotFound) {
				t.Fatalf("expected missing handler error, got %v", err)
			}
			if !tt.missing && tt.handlerErr == nil && tt.finishErr == nil && err != nil {
				t.Fatal(err)
			}
			if tt.action == "retry" && !repo.availableAt.Equal(now.Add(time.Second)) {
				t.Fatal("incorrect retry instant")
			}
		})
	}
}

func TestEventProcessorIdleAndCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	repo := &deliveryStub{claimErr: ErrNoEventAvailable}
	handlers := map[domain.EventType]EventHandler{"mail.requested": handlerFunc(func(context.Context, *domain.Event) error { cancel(); return nil })}
	processor, err := NewEventProcessor(repo, processorClock{time.Now()}, handlers, EventProcessorOptions{LeaseDuration: time.Minute, HandleTimeout: time.Second, RetryDelay: time.Second, MaxAttempts: 3})
	if err != nil {
		t.Fatal(err)
	}
	if processed, err := processor.ProcessNext(ctx); processed || err != nil {
		t.Fatalf("idle: %v %v", processed, err)
	}
	event, _ := domain.NewEvent("event-1", "mail.requested", "", []byte("payload"))
	repo.claimErr = nil
	repo.claim = ClaimedEvent{Event: event, LeaseToken: "lease-1", Attempt: 1}
	if processed, err := processor.ProcessNext(ctx); !processed || !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel: %v %v", processed, err)
	}
	if repo.action != "" {
		t.Fatal("cancellation must leave lease for recovery")
	}
}
