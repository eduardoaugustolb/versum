package domain

import (
	"time"
)

type Event struct {
	id                string
	eventType         string
	payload           []byte
	attempts          int
	availableAt       time.Time
	leaseToken        string
	leasedUntil       time.Time
	processedAt       time.Time
	failedAt          time.Time
	lastErrorRedacted string
	createdAt         time.Time
}

func NewEvent(
	id string,
	eventType string,
	payload []byte,
	availableAt time.Time,
	createdAt time.Time,
) (*Event, error) {
	e := &Event{
		id:          id,
		eventType:   eventType,
		payload:     append([]byte{}, payload...),
		availableAt: availableAt,
		createdAt:   createdAt,
	}
	if err := validateEvent(e); err != nil {
		return nil, err
	}
	return e, nil
}

func RehydrateEvent(
	id string,
	eventType string,
	payload []byte,
	attempts int,
	availableAt time.Time,
	leaseToken string,
	leasedUntil time.Time,
	processedAt time.Time,
	failedAt time.Time,
	lastErrorRedacted string,
	createdAt time.Time,
) (*Event, error) {
	e := &Event{
		id:                id,
		eventType:         eventType,
		payload:           append([]byte{}, payload...),
		attempts:          attempts,
		availableAt:       availableAt,
		leaseToken:        leaseToken,
		leasedUntil:       leasedUntil,
		processedAt:       processedAt,
		failedAt:          failedAt,
		lastErrorRedacted: lastErrorRedacted,
		createdAt:         createdAt,
	}
	if err := validateEvent(e); err != nil {
		return nil, err
	}
	return e, nil
}

func validateEvent(e *Event) error {
	if e.id == "" {
		return ErrInvalidId
	}
	if e.eventType == "" {
		return ErrInvalidEventType
	}
	if e.availableAt.IsZero() || e.availableAt.Before(e.createdAt) {
		return ErrInvalidAvailableAt
	}

	if e.createdAt.IsZero() {
		return ErrInvalidCreatedAt
	}

	if len(e.payload) == 0 {
		return ErrInvalidPayload
	}
	return nil
}

func (e *Event) ID() string {
	return e.id
}
func (e *Event) EventType() string {
	return e.eventType
}
func (e *Event) Payload() []byte {
	return append([]byte{}, e.payload...)
}
func (e *Event) Attempts() int {
	return e.attempts
}
func (e *Event) AvailableAt() time.Time {
	return e.availableAt
}
func (e *Event) LeaseToken() string {
	return e.leaseToken
}
func (e *Event) LeasedUntil() time.Time {
	return e.leasedUntil
}
func (e *Event) ProcessedAt() time.Time {
	return e.processedAt
}
func (e *Event) FailedAt() time.Time {
	return e.failedAt
}
func (e *Event) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Event) LastErrorRedacted() string {
	return e.lastErrorRedacted
}
