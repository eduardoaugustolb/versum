package domain

import "errors"

var (
	ErrInvalidId          = errors.New("id is invalid")
	ErrInvalidEventType   = errors.New("eventType is invalid")
	ErrInvalidAvailableAt = errors.New("availableAt is invalid")
	ErrInvalidCreatedAt   = errors.New("createdAt is invalid")
	ErrInvalidPayload     = errors.New("payload is invalid")
)
