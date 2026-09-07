package domain

import "time"

type Session struct {
	id         string
	secretHash []byte
	userID     string
	revokedAt  *time.Time
	usedAt     *time.Time
	expiresAt  *time.Time
}

func NewSession(id string, secretHash []byte, userId string, expiresAt time.Time, now time.Time) (*Session, error) {
	if !expiresAt.After(now) {
		return nil, ErrInvalidSessionExpiresAt
	}

	session := &Session{
		id:         id,
		secretHash: cloneBytes(secretHash),
		userID:     userId,
		expiresAt:  cloneTime(&expiresAt),
	}

	if err := validateSessionOptions(session); err != nil {
		return nil, err
	}

	return session, nil
}

func RehydrateSession(id string, secretHash []byte, userID string, revokedAt *time.Time, usedAt *time.Time, expiresAt time.Time) (*Session, error) {
	session := &Session{
		id:         id,
		secretHash: cloneBytes(secretHash),
		userID:     userID,
		revokedAt:  cloneTime(revokedAt),
		usedAt:     cloneTime(usedAt),
		expiresAt:  cloneTime(&expiresAt),
	}

	if err := validateSessionOptions(session); err != nil {
		return nil, err
	}

	return session, nil
}

func validateSessionOptions(session *Session) error {
	if session.id == "" {
		return ErrInvalidSessionID
	}

	if len(session.secretHash) == 0 {
		return ErrInvalidSessionSecretHash
	}

	if session.usedAt != nil && session.usedAt.IsZero() {
		return ErrInvalidSessionUsedAt
	}

	if session.userID == "" {
		return ErrInvalidSessionUserID
	}

	if session.expiresAt == nil || session.expiresAt.IsZero() {
		return ErrInvalidSessionExpiresAt
	}

	if session.revokedAt != nil && session.revokedAt.IsZero() {
		return ErrInvalidSessionRevokedAt
	}

	return nil
}

func (s *Session) ID() string {
	return s.id
}
func (s *Session) SecretHash() []byte {
	return cloneBytes(s.secretHash)
}
func (s *Session) UserID() string {
	return s.userID
}
func (s *Session) RevokedAt() (time.Time, bool) {
	if s.revokedAt == nil {
		return time.Time{}, false
	}
	return *s.revokedAt, true
}
func (s *Session) UsedAt() (time.Time, bool) {
	if s.usedAt == nil {
		return time.Time{}, false
	}
	return *s.usedAt, true
}
func (s *Session) ExpiresAt() time.Time {
	return *s.expiresAt
}

func (s *Session) Revoke(now time.Time) {
	s.revokedAt = &now
}

func (s *Session) Use(now time.Time) error {
	if !s.IsValidAt(now) {
		return ErrInvalidSessionState
	}

	usedAt := now
	s.usedAt = &usedAt
	return nil
}

func (s *Session) IsRevoked() bool {
	return s.revokedAt != nil
}

func (s *Session) IsExpired(now time.Time) bool {
	return s.expiresAt == nil || !s.expiresAt.After(now)
}

func (s *Session) IsValidAt(now time.Time) bool {
	return !s.IsRevoked() && !s.IsExpired(now)
}
