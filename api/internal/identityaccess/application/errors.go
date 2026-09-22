// Package application defines use cases and their observable outcomes.
package application

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrSessionNotFound         = errors.New("session not found")
	ErrSessionAlreadyExists    = errors.New("session already exists")
	ErrLoginTokenNotFound      = errors.New("login token not found")
	ErrLoginTokenAlreadyExists = errors.New("login token already exists")
	ErrInvalidMagicLinkTTL     = errors.New("invalid magic link ttl")
	ErrInvalidMagicLinkPolicy  = errors.New("invalid magic link policy")
)
