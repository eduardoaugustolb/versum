package domain_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestSessionClientAndFamilyValidation(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name, family, ip, agent string
		want                    error
	}{
		{"ipv4", "family", "192.0.2.1", "browser", nil},
		{"ipv6", "family", "2001:db8::1", "browser", nil},
		{"unknown client", "family", "", "", nil},
		{"empty family", "", "", "", domain.ErrInvalidSessionFamilyID},
		{"blank family", " ", "", "", domain.ErrInvalidSessionFamilyID},
		{"invalid ip", "family", "not-an-ip", "", domain.ErrInvalidSessionIPAddress},
		{"ip with port", "family", "192.0.2.1:443", "", domain.ErrInvalidSessionIPAddress},
		{"ipv6 zone", "family", "fe80::1%eth0", "", domain.ErrInvalidSessionIPAddress},
		{"agent too long", "family", "", strings.Repeat("a", 1025), domain.ErrInvalidSessionUserAgent},
		{"agent newline", "family", "", "browser\nheader", domain.ErrInvalidSessionUserAgent},
		{"agent nul", "family", "", "browser\x00", domain.ErrInvalidSessionUserAgent},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := domain.NewSession("session", "user", tt.family, tt.ip, tt.agent, now.Add(time.Hour), now)
			if !errors.Is(err, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, err)
			}
			if err == nil && (s.FamilyID() != tt.family || s.IPAddress() != tt.ip || s.UserAgent() != tt.agent) {
				t.Fatal("incorrect client or family")
			}
		})
	}
	s, err := domain.NewSession("session", "user", "family", "::ffff:192.0.2.1", "", now.Add(time.Hour), now)
	if err != nil || s.IPAddress() != "192.0.2.1" {
		t.Fatal("IPv4-mapped IPv6 must normalize")
	}
}

func TestSessionRotationTransitions(t *testing.T) {
	now := time.Now().UTC()
	old, _ := domain.NewSession("old", "user", "family", "", "", now.Add(time.Hour), now)
	next, _ := domain.NewSession("next", "user", "family", "192.0.2.2", "new browser", old.ExpiresAt(), now)
	if err := old.ReplaceWith(next, now); err != nil {
		t.Fatal(err)
	}
	if old.IsValidAt(now) || !next.IsValidAt(now) {
		t.Fatal("only the successor must remain valid")
	}
	if id, ok := old.ReplacedBySessionID(); !ok || id != next.ID() {
		t.Fatal("replacement link missing")
	}
	if at, ok := old.LastUsedAt(); !ok || !at.Equal(now) {
		t.Fatal("last use missing")
	}
	if at, ok := old.RevokedAt(); !ok || !at.Equal(now) {
		t.Fatal("revocation missing")
	}
	if err := old.Use(now); !errors.Is(err, domain.ErrSessionReused) {
		t.Fatal(err)
	}
	if err := old.ReplaceWith(next, now); !errors.Is(err, domain.ErrSessionReused) {
		t.Fatal(err)
	}
	if err := old.Revoke(now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if at, _ := old.RevokedAt(); !at.Equal(now) {
		t.Fatal("revocation must be idempotent")
	}
}

func TestSessionRejectsInvalidRotationsWithoutMutation(t *testing.T) {
	now := time.Now().UTC()
	for _, name := range []string{"nil", "self", "other user", "other family", "revoked successor", "expired successor", "extended expiry", "used successor", "expired predecessor", "revoked predecessor", "backward time"} {
		t.Run(name, func(t *testing.T) {
			old, _ := domain.NewSession("old", "user", "family", "", "", now.Add(time.Hour), now)
			next, _ := domain.NewSession("next", "user", "family", "", "", old.ExpiresAt(), now)
			at := now
			want := domain.ErrInvalidSessionReplacement
			switch name {
			case "nil":
				next = nil
			case "self":
				next = old
			case "other user":
				next, _ = domain.NewSession("next", "other", "family", "", "", old.ExpiresAt(), now)
			case "other family":
				next, _ = domain.NewSession("next", "user", "other", "", "", old.ExpiresAt(), now)
			case "revoked successor":
				next.Revoke(now)
			case "expired successor":
				next, _ = domain.NewSession("next", "user", "family", "", "", now.Add(time.Second), now)
				at = now.Add(time.Second)
			case "extended expiry":
				next, _ = domain.NewSession("next", "user", "family", "", "", old.ExpiresAt().Add(time.Second), now)
			case "used successor":
				next.Use(now)
			case "expired predecessor":
				at = old.ExpiresAt()
				want = domain.ErrInvalidSessionState
			case "revoked predecessor":
				old.Revoke(now)
				want = domain.ErrInvalidSessionState
			case "backward time":
				old.Use(now.Add(time.Minute))
				want = domain.ErrInvalidSessionLastUsedAt
			}
			before, used := old.LastUsedAt()
			if err := old.ReplaceWith(next, at); !errors.Is(err, want) {
				t.Fatalf("want %v, got %v", want, err)
			}
			if _, replaced := old.ReplacedBySessionID(); replaced {
				t.Fatal("failed rotation mutated replacement")
			}
			after, stillUsed := old.LastUsedAt()
			if used != stillUsed || !before.Equal(after) {
				t.Fatal("failed rotation mutated last use")
			}
		})
	}
}

func TestSessionTimeErrorsAndRehydration(t *testing.T) {
	now := time.Now().UTC()
	expiry := now.Add(time.Hour)
	zero := time.Time{}
	s, _ := domain.NewSession("session", "user", "family", "", "", expiry, now)
	if err := s.Use(now); err != nil {
		t.Fatal(err)
	}
	if err := s.Use(now.Add(-time.Second)); !errors.Is(err, domain.ErrInvalidSessionLastUsedAt) {
		t.Fatal(err)
	}
	if err := s.Use(expiry); !errors.Is(err, domain.ErrInvalidSessionState) {
		t.Fatal(err)
	}
	if err := s.Use(zero); !errors.Is(err, domain.ErrInvalidSessionState) {
		t.Fatal(err)
	}
	if err := s.Revoke(zero); !errors.Is(err, domain.ErrInvalidSessionRevokedAt) {
		t.Fatal(err)
	}
	if s.IsRevoked() {
		t.Fatal("invalid revocation mutated session")
	}
	if _, err := domain.NewSession("session", "user", "family", "", "", expiry, zero); !errors.Is(err, domain.ErrInvalidSessionExpiresAt) {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		replacement   string
		revoked, last *time.Time
		want          error
	}{
		{"next", &now, &now, nil}, {"next", nil, nil, domain.ErrInvalidSessionReplacement}, {"session", &now, nil, domain.ErrInvalidSessionReplacement}, {" ", &now, nil, domain.ErrInvalidSessionReplacement}, {"", nil, &expiry, domain.ErrInvalidSessionLastUsedAt},
	} {
		s, err := domain.RehydrateSession("session", "user", "family", "", "", tt.revoked, tt.last, expiry, tt.replacement)
		if !errors.Is(err, tt.want) {
			t.Fatalf("want %v, got %v", tt.want, err)
		}
		if err == nil {
			if id, _ := s.ReplacedBySessionID(); id != tt.replacement {
				t.Fatal("lost persisted replacement")
			}
		}
	}
}
