// Package application defines use cases and their observable outcomes.
package application

import "errors"

var ErrOutboxEventAlreadyExists = errors.New("outbox event already exists")
