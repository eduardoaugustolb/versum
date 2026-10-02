package commands_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/commands"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestRevokeSessionFamily_RevokesSessionFamily(t *testing.T) {
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	repo := &FakeSessionRepository{}
	for _, tt := range []struct{ id, user, family, secret string }{
		{"first", "user", "family", "secret"},
		{"second", "user", "family", "another-secret"},
		{"unrelated", "user", "other-family", "other-secret"},
		{"other-owner", "other-user", "family", "other-owner-secret"},
	} {
		session, err := domain.NewSession(tt.id, tt.user, tt.family, "192.0.2.1", "browser", now.Add(time.Hour), now)
		if err != nil {
			t.Fatal(err)
		}
		if err := repo.CreateSession(t.Context(), session, tt.secret); err != nil {
			t.Fatal(err)
		}
	}
	uc := commands.NewRevokeSessionFamily(repo, FakeClock{NowTime: now})
	if err := uc.Execute(t.Context(), "secret"); err != nil {
		t.Fatal(err)
	}
	sessions, err := repo.ListSessionsByUserID(t.Context(), "user")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 3 {
		t.Fatalf("expected 3 sessions, got %d", len(sessions))
	}
	for _, s := range repo.ByID {
		shouldRevoke := s.UserID() == "user" && s.FamilyID() == "family"
		if s.IsRevoked() != shouldRevoke {
			t.Fatalf("incorrect revocation for session %s", s.ID())
		}
		if shouldRevoke {
			if at, ok := s.RevokedAt(); !ok || !at.Equal(now) {
				t.Fatal("incorrect revocation timestamp")
			}
		}
	}
	if len(repo.RevokeFamilyCalls) != 1 || repo.RevokeFamilyCalls[0].UserID != "user" || repo.RevokeFamilyCalls[0].FamilyID != "family" {
		t.Fatalf("unexpected family revocation: %+v", repo.RevokeFamilyCalls)
	}
	if err := uc.Execute(t.Context(), "secret"); err != nil {
		t.Fatalf("logout must be idempotent: %v", err)
	}
}

func TestRevokeSessionFamily_ReturnsErrorWhenSessionExpired(t *testing.T) {
	now := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	for _, expiry := range []time.Time{now.Add(-time.Hour), now} {
		session, err := domain.NewSession("id", "user", "family", "192.0.2.1", "browser", expiry, now.Add(-2*time.Hour))
		if err != nil {
			t.Fatal(err)
		}
		repo := &FakeSessionRepository{}
		if err := repo.CreateSession(t.Context(), session, "secret"); err != nil {
			t.Fatal(err)
		}
		uc := commands.NewRevokeSessionFamily(repo, FakeClock{NowTime: now})
		if err := uc.Execute(t.Context(), "secret"); !errors.Is(err, application.ErrSessionExpired) {
			t.Fatalf("expected expired session, got %v", err)
		}
		found, err := repo.FindSessionBySecret(t.Context(), "secret")
		if err != nil {
			t.Fatal(err)
		}
		if found == nil || !found.IsExpired(now) || found.IsRevoked() {
			t.Fatal("expired session must remain stored without mutation")
		}
		if len(repo.RevokeFamilyCalls) != 0 {
			t.Fatal("expired credentials must not revoke a family")
		}
	}
}

func TestRevokeSessionFamily_ReturnsErrorWhenSecretDoesNotExist(t *testing.T) {
	repo := &FakeSessionRepository{}
	uc := commands.NewRevokeSessionFamily(repo, FakeClock{})
	if err := uc.Execute(t.Context(), "missing"); !errors.Is(err, application.ErrSessionNotFound) {
		t.Fatalf("expected session not found, got %v", err)
	}
	if len(repo.RevokeFamilyCalls) != 0 {
		t.Fatal("missing credentials must not revoke a family")
	}
}

func TestRevokeSessionFamilyPropagatesRepositoryErrors(t *testing.T) {
	now := time.Now().UTC()
	failure := errors.New("repository unavailable")
	for _, name := range []string{"lookup", "nil session", "revocation"} {
		t.Run(name, func(t *testing.T) {
			session, err := domain.NewSession("id", "user", "family", "", "", now.Add(time.Hour), now)
			if err != nil {
				t.Fatal(err)
			}
			repo := &FakeSessionRepository{ByID: map[string]*domain.Session{session.ID(): session}, BySecret: map[string]*domain.Session{"secret": session}}
			want := failure
			switch name {
			case "lookup":
				repo.FindBySecretErr = failure
			case "nil session":
				repo.OnFindBySecret = func(context.Context, string) (*domain.Session, error) { return nil, nil }
				want = application.ErrSessionNotFound
			case "revocation":
				repo.RevokeErr = failure
			}
			uc := commands.NewRevokeSessionFamily(repo, FakeClock{NowTime: now})
			if err := uc.Execute(t.Context(), "secret"); !errors.Is(err, want) {
				t.Fatalf("expected %v, got %v", want, err)
			}
			if session.IsRevoked() {
				t.Fatal("failed revocation must not mutate the session")
			}
		})
	}
}

func TestRevokeSessionFamilyAcceptsReplacedCredential(t *testing.T) {
	now := time.Now().UTC()
	previous, err := domain.RehydrateSession("old", "user", "family", "", "", &now, &now, now.Add(time.Hour), "next")
	if err != nil {
		t.Fatal(err)
	}
	next, err := domain.NewSession("next", "user", "family", "", "", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	repo := &FakeSessionRepository{ByID: map[string]*domain.Session{"old": previous, "next": next}, BySecret: map[string]*domain.Session{"old-secret": previous}}
	uc := commands.NewRevokeSessionFamily(repo, FakeClock{NowTime: now})
	if err := uc.Execute(t.Context(), "old-secret"); err != nil {
		t.Fatal(err)
	}
	if !next.IsRevoked() {
		t.Fatal("logout using a replaced credential must revoke its successor")
	}
}
