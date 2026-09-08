// Package application defines use cases and their observable outcomes.
package application

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrLoginTokenNotFound = errors.New("login token not found")
)
