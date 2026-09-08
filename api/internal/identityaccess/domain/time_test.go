package domain_test

import (
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestIdentityAccessTimesAreStoredInUTC(t *testing.T) {
	brt := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, time.January, 1, 10, 0, 0, 0, brt)
	expiresAt := now.Add(time.Hour)

	session, err := domain.NewSession("session-1", []byte("hash"), "user-1", expiresAt, now)
	if err != nil {
		t.Fatal(err)
	}
	if session.ExpiresAt().Location() != time.UTC {
		t.Fatalf("expected UTC session expiration, got %s", session.ExpiresAt().Location())
	}
	session.Revoke(now)
	if revokedAt, ok := session.RevokedAt(); !ok || revokedAt.Location() != time.UTC {
		t.Fatalf("expected UTC session revocation, got %s", revokedAt.Location())
	}

	token, err := domain.NewLoginToken("token-1", []byte("hash"), "user-1", expiresAt, now)
	if err != nil {
		t.Fatal(err)
	}
	if token.ExpiresAt().Location() != time.UTC {
		t.Fatalf("expected UTC token expiration, got %s", token.ExpiresAt().Location())
	}
	if err := token.Consume(now); err != nil {
		t.Fatal(err)
	}
	if consumedAt, ok := token.ConsumedAt(); !ok || consumedAt.Location() != time.UTC {
		t.Fatalf("expected UTC token consumption, got %s", consumedAt.Location())
	}
}
