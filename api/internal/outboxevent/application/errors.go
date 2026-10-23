// Package application defines use cases and their observable outcomes.
package application

import "errors"

var ErrOutboxEventAlreadyExists = errors.New("outbox event already exists")

var (
	ErrNoEventAvailable      = errors.New("no outbox event available")
	ErrEventLeaseLost        = errors.New("outbox event lease lost")
	ErrEventHandlerNotFound  = errors.New("outbox event handler not found")
	ErrInvalidEventProcessor = errors.New("invalid outbox event processor configuration")
	ErrInvalidEventClaim     = errors.New("invalid outbox event claim")
)
