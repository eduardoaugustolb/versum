package domain

import "errors"

var (
	ErrInvalidEmail = errors.New("invalid email")
)

var (
	ErrInvalidSessionFamilyID    = errors.New("invalid session family id")
	ErrInvalidSessionIPAddress   = errors.New("invalid session ip address")
	ErrInvalidSessionUserAgent   = errors.New("invalid session user agent")
	ErrInvalidSessionReplacement = errors.New("invalid session replacement")
	ErrSessionReused             = errors.New("replaced session reused")
	ErrInvalidSessionID          = errors.New("invalid session id")
	ErrInvalidSessionUserID      = errors.New("invalid session user id")
	ErrInvalidSessionExpiresAt   = errors.New("invalid session expires at")
	ErrInvalidSessionState       = errors.New("invalid session state")
	ErrInvalidSessionRevokedAt   = errors.New("invalid session revoked at")
	ErrInvalidSessionLastUsedAt  = errors.New("invalid session last used at")
)

var (
	ErrInvalidUserID = errors.New("invalid user id")
)

var (
	ErrInvalidLoginTokenID         = errors.New("invalid login token id")
	ErrInvalidLoginTokenUserID     = errors.New("invalid login token user id")
	ErrInvalidLoginTokenExpiresAt  = errors.New("invalid login token expires at")
	ErrInvalidLoginTokenConsumedAt = errors.New("invalid login token consumed at")
	ErrLoginTokenAlreadyConsumed   = errors.New("login token already consumed")
	ErrLoginTokenExpired           = errors.New("login token expired")
)
