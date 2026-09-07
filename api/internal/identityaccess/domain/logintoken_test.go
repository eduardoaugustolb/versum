package domain_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestNewLoginToken(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		id        string
		hash      []byte
		userID    string
		expiresAt time.Time
		wantErr   error
	}{
		{"valid", "token-1", []byte("hash"), "user-1", now.Add(time.Hour), nil},
		{"missing id", "", []byte("hash"), "user-1", now.Add(time.Hour), domain.ErrInvalidLoginTokenID},
		{"empty hash", "token-1", []byte{}, "user-1", now.Add(time.Hour), domain.ErrInvalidLoginTokenHash},
		{"missing user", "token-1", []byte("hash"), "", now.Add(time.Hour), domain.ErrInvalidLoginTokenUserID},
		{"zero expiration", "token-1", []byte("hash"), "user-1", time.Time{}, domain.ErrInvalidLoginTokenExpiresAt},
		{"expiration equal to now", "token-1", []byte("hash"), "user-1", now, domain.ErrInvalidLoginTokenExpiresAt},
		{"expiration in the past", "token-1", []byte("hash"), "user-1", now.Add(-time.Second), domain.ErrInvalidLoginTokenExpiresAt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := domain.NewLoginToken(tt.id, tt.hash, tt.userID, tt.expiresAt, now)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && token == nil {
				t.Fatal("valid token was not created")
			}
		})
	}
}

func TestLoginTokenRehydrateOptionalConsumedAt(t *testing.T) {
	expiresAt := time.Date(2026, time.January, 1, 13, 0, 0, 0, time.UTC)
	token, err := domain.RehydrateLoginToken("token-1", []byte("hash"), "user-1", expiresAt, nil)
	if err != nil {
		t.Fatalf("rehydration failed: %v", err)
	}
	if token.IsConsumed() {
		t.Fatal("token with nil consumed_at was marked as consumed")
	}
	if _, ok := token.ConsumedAt(); ok {
		t.Fatal("token without consumed_at returned a timestamp")
	}

	consumedAt := expiresAt.Add(-time.Minute)
	consumed, err := domain.RehydrateLoginToken("token-1", []byte("hash"), "user-1", expiresAt, &consumedAt)
	if err != nil {
		t.Fatalf("rehydration of consumed token failed: %v", err)
	}
	got, ok := consumed.ConsumedAt()
	if !ok || !got.Equal(consumedAt) {
		t.Fatalf("incorrect consumed_at: got %v, ok=%v", got, ok)
	}

	zero := time.Time{}
	if _, err := domain.RehydrateLoginToken("token-1", []byte("hash"), "user-1", expiresAt, &zero); !errors.Is(err, domain.ErrInvalidLoginTokenConsumedAt) {
		t.Fatalf("expected zero consumed_at error, got %v", err)
	}
}

func TestLoginTokenConsume(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	token, err := domain.NewLoginToken("token-1", []byte("hash"), "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := token.Consume(now.Add(10 * time.Minute)); err != nil {
		t.Fatalf("valid consumption failed: %v", err)
	}
	if !token.IsConsumed() {
		t.Fatal("token was not marked as consumed")
	}
	if err := token.Consume(now.Add(20 * time.Minute)); !errors.Is(err, domain.ErrLoginTokenAlreadyConsumed) {
		t.Fatalf("expected duplicate consumption error, got %v", err)
	}

	expired, err := domain.NewLoginToken("token-2", []byte("hash"), "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := expired.Consume(now.Add(time.Hour)); !errors.Is(err, domain.ErrLoginTokenExpired) {
		t.Fatalf("expected expired token error, got %v", err)
	}
}

func TestLoginTokenDefensiveCopies(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	hash := []byte("hash")
	token, err := domain.NewLoginToken("token-1", hash, "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	hash[0] = 'X'
	if bytes.Equal(token.TokenHash(), hash) {
		t.Fatal("token retained the input hash reference")
	}
	returned := token.TokenHash()
	returned[0] = 'Y'
	if bytes.Equal(token.TokenHash(), returned) {
		t.Fatal("getter exposed the internal hash reference")
	}
}

func TestLoginTokenIsExpired(t *testing.T) {
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	token, err := domain.NewLoginToken("token-1", []byte("hash"), "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	if token.IsExpired(now.Add(59 * time.Minute)) {
		t.Fatal("token expired before its deadline")
	}
	if !token.IsExpired(now.Add(time.Hour)) {
		t.Fatal("token should expire at its deadline")
	}
}
