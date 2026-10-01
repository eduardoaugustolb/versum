package domain

import (
	"net/netip"
	"strings"
	"time"
)

// SessionClient describes the client observed when this session was issued.
// These attributes are risk signals, not proof of possession of a credential.
type SessionClient struct {
	IPAddress string
	UserAgent string
}

type Session struct {
	id, userID, familyID  string
	client                SessionClient
	revokedAt, lastUsedAt *time.Time
	expiresAt             time.Time
	replacedBySessionID   string
}

func NewSession(id, userID, familyID string, client SessionClient, expiresAt, now time.Time) (*Session, error) {
	if now.IsZero() || !expiresAt.After(now) {
		return nil, ErrInvalidSessionExpiresAt
	}
	return RehydrateSession(id, userID, familyID, client, nil, nil, expiresAt, "")
}

func RehydrateSession(id, userID, familyID string, client SessionClient, revokedAt, lastUsedAt *time.Time, expiresAt time.Time, replacedBySessionID string) (*Session, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidSessionID
	}
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidSessionUserID
	}
	if strings.TrimSpace(familyID) == "" {
		return nil, ErrInvalidSessionFamilyID
	}
	if expiresAt.IsZero() {
		return nil, ErrInvalidSessionExpiresAt
	}
	if revokedAt != nil && revokedAt.IsZero() {
		return nil, ErrInvalidSessionRevokedAt
	}
	if lastUsedAt != nil && (lastUsedAt.IsZero() || !lastUsedAt.Before(expiresAt)) {
		return nil, ErrInvalidSessionLastUsedAt
	}
	if replacedBySessionID != "" && (strings.TrimSpace(replacedBySessionID) == "" || replacedBySessionID == id || revokedAt == nil) {
		return nil, ErrInvalidSessionReplacement
	}
	if client.IPAddress != "" {
		addr, err := netip.ParseAddr(client.IPAddress)
		if err != nil || addr.Zone() != "" {
			return nil, ErrInvalidSessionIPAddress
		}
		client.IPAddress = addr.Unmap().String()
	}
	if len(client.UserAgent) > 1024 || strings.ContainsAny(client.UserAgent, "\r\n\x00") {
		return nil, ErrInvalidSessionUserAgent
	}
	return &Session{id: id, userID: userID, familyID: familyID, client: client, revokedAt: cloneTime(revokedAt), lastUsedAt: cloneTime(lastUsedAt), expiresAt: expiresAt.UTC(), replacedBySessionID: replacedBySessionID}, nil
}

func (s *Session) ID() string        { return s.id }
func (s *Session) UserID() string    { return s.userID }
func (s *Session) FamilyID() string  { return s.familyID }
func (s *Session) IPAddress() string { return s.client.IPAddress }
func (s *Session) UserAgent() string { return s.client.UserAgent }
func (s *Session) ReplacedBySessionID() (string, bool) {
	return s.replacedBySessionID, s.replacedBySessionID != ""
}
func (s *Session) RevokedAt() (time.Time, bool) {
	if s.revokedAt == nil {
		return time.Time{}, false
	}
	return *s.revokedAt, true
}
func (s *Session) LastUsedAt() (time.Time, bool) {
	if s.lastUsedAt == nil {
		return time.Time{}, false
	}
	return *s.lastUsedAt, true
}
func (s *Session) ExpiresAt() time.Time         { return s.expiresAt }
func (s *Session) IsRevoked() bool              { return s.revokedAt != nil }
func (s *Session) IsExpired(now time.Time) bool { return !s.expiresAt.After(now) }
func (s *Session) IsValidAt(now time.Time) bool {
	return !now.IsZero() && !s.IsRevoked() && s.replacedBySessionID == "" && !s.IsExpired(now)
}

// Revoke is idempotent and preserves the first revocation timestamp.
func (s *Session) Revoke(now time.Time) error {
	if now.IsZero() {
		return ErrInvalidSessionRevokedAt
	}
	if s.revokedAt == nil {
		now = now.UTC()
		s.revokedAt = &now
	}
	return nil
}

func (s *Session) Use(now time.Time) error {
	if s.replacedBySessionID != "" {
		return ErrSessionReused
	}
	if !s.IsValidAt(now) {
		return ErrInvalidSessionState
	}
	if s.lastUsedAt != nil && now.Before(*s.lastUsedAt) {
		return ErrInvalidSessionLastUsedAt
	}
	now = now.UTC()
	s.lastUsedAt = &now
	return nil
}

// ReplaceWith invalidates this credential without extending the family's lifetime.
func (s *Session) ReplaceWith(next *Session, now time.Time) error {
	if s.replacedBySessionID != "" {
		return ErrSessionReused
	}
	if !s.IsValidAt(now) {
		return ErrInvalidSessionState
	}
	if next == nil || next.id == s.id || next.familyID != s.familyID || next.userID != s.userID || !next.IsValidAt(now) || next.expiresAt.After(s.expiresAt) || next.lastUsedAt != nil {
		return ErrInvalidSessionReplacement
	}
	if err := s.Use(now); err != nil {
		return err
	}
	if err := s.Revoke(now); err != nil {
		return err
	}
	s.replacedBySessionID = next.id
	return nil
}
