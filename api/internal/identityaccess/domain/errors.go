package domain

import "errors"

var (
	ErrInvalidEmail = errors.New("invalid email")
)

var (
	ErrInvalidSessionID         = errors.New("invalid session id")
	ErrInvalidSessionSecretHash = errors.New("invalid session secret hash")
	ErrInvalidSessionUserID     = errors.New("invalid session user id")
	ErrInvalidSessionExpiresAt  = errors.New("invalid session expires at")
	ErrInvalidSessionState      = errors.New("invalid session state")
	ErrInvalidSessionRevokedAt  = errors.New("invalid session revoked at")
	ErrInvalidSessionUsedAt     = errors.New("invalid session used at")
)

var (
	ErrInvalidUserID = errors.New("invalid user id")
)

var (
	ErrInvalidLoginTokenID         = errors.New("invalid login token id")
	ErrInvalidLoginTokenHash       = errors.New("invalid login token hash")
	ErrInvalidLoginTokenUserID     = errors.New("invalid login token user id")
	ErrInvalidLoginTokenExpiresAt  = errors.New("invalid login token expires at")
	ErrInvalidLoginTokenConsumedAt = errors.New("invalid login token consumed at")
	ErrLoginTokenAlreadyConsumed   = errors.New("login token already consumed")
	ErrLoginTokenExpired           = errors.New("login token expired")
)
