package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestNewSession(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		id        string
		userID    string
		expiresAt time.Time
		wantErr   error
	}{
		{"valid", "session-1", "user-1", now.Add(time.Hour), nil},
		{"missing id", "", "user-1", now.Add(time.Hour), domain.ErrInvalidSessionID},
		{"missing user", "session-1", "", now.Add(time.Hour), domain.ErrInvalidSessionUserID},
		{"zero expiration", "session-1", "user-1", time.Time{}, domain.ErrInvalidSessionExpiresAt},
		{"expiration equal to now", "session-1", "user-1", now, domain.ErrInvalidSessionExpiresAt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := domain.NewSession(tt.id, tt.userID, "family-1", domain.SessionClient{}, tt.expiresAt, now)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && session == nil {
				t.Fatal("valid session was not created")
			}
		})
	}
}

func TestSessionRehydrateOptionalState(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	expiresAt := now.Add(time.Hour)
	session, err := domain.RehydrateSession("session-1", "user-1", "family-1", domain.SessionClient{}, nil, nil, expiresAt, "")
	if err != nil {
		t.Fatalf("rehydration of active session failed: %v", err)
	}
	if session.IsRevoked() {
		t.Fatal("active session was marked as revoked")
	}
	if _, ok := session.RevokedAt(); ok {
		t.Fatal("session without revoked_at returned a timestamp")
	}

	revokedAt := now.Add(10 * time.Minute)
	usedAt := now.Add(5 * time.Minute)
	revoked, err := domain.RehydrateSession("session-1", "user-1", "family-1", domain.SessionClient{}, &revokedAt, &usedAt, expiresAt, "")
	if err != nil {
		t.Fatalf("rehydration of revoked session failed: %v", err)
	}
	if !revoked.IsRevoked() {
		t.Fatal("revoked session was not marked as revoked")
	}
	if got, ok := revoked.LastUsedAt(); !ok || !got.Equal(usedAt) {
		t.Fatalf("incorrect used_at: got %v, ok=%v", got, ok)
	}

	zero := time.Time{}
	if _, err := domain.RehydrateSession("session-1", "user-1", "family-1", domain.SessionClient{}, &zero, nil, expiresAt, ""); !errors.Is(err, domain.ErrInvalidSessionRevokedAt) {
		t.Fatalf("expected zero revoked_at error, got %v", err)
	}
	if _, err := domain.RehydrateSession("session-1", "user-1", "family-1", domain.SessionClient{}, nil, &zero, expiresAt, ""); !errors.Is(err, domain.ErrInvalidSessionLastUsedAt) {
		t.Fatalf("expected zero used_at error, got %v", err)
	}
}

func TestSessionValidityAndUse(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	session, err := domain.NewSession("session-1", "user-1", "family-1", domain.SessionClient{}, now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	if !session.IsValidAt(now) {
		t.Fatal("new session should be valid")
	}
	if err := session.Use(now.Add(10 * time.Minute)); err != nil {
		t.Fatalf("using valid session failed: %v", err)
	}

	session.Revoke(now.Add(20 * time.Minute))
	if session.IsValidAt(now.Add(20 * time.Minute)) {
		t.Fatal("revoked session was considered valid")
	}
	if err := session.Use(now.Add(21 * time.Minute)); !errors.Is(err, domain.ErrInvalidSessionState) {
		t.Fatalf("expected revoked session error, got %v", err)
	}
}

func TestSessionRehydrateCopiesTimestamps(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	usedAt, revokedAt := now, now.Add(time.Minute)
	session, err := domain.RehydrateSession("session-1", "user-1", "family-1", domain.SessionClient{}, &revokedAt, &usedAt, now.Add(time.Hour), "")
	if err != nil {
		t.Fatal(err)
	}
	usedAt, revokedAt = time.Time{}, time.Time{}
	if at, ok := session.LastUsedAt(); !ok || !at.Equal(now) {
		t.Fatal("session retained the input used_at reference")
	}
	if at, ok := session.RevokedAt(); !ok || !at.Equal(now.Add(time.Minute)) {
		t.Fatal("session retained the input revoked_at reference")
	}
}
