package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

func TestLoginTokenRepositoryTranslatesUniqueViolationOnCreate(t *testing.T) {
	now := time.Now()
	token, err := domain.NewLoginToken("token-1", []byte("hash"), "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewLoginTokenRepository(userRepositoryExecutor{execErr: &pgconn.PgError{Code: uniqueViolationCode}})

	if err := repo.CreateLoginToken(t.Context(), token); !errors.Is(err, application.ErrLoginTokenAlreadyExists) {
		t.Fatalf("expected already-exists error, got %v", err)
	}
}

func TestSessionRepositoryTranslatesUniqueViolationOnCreate(t *testing.T) {
	now := time.Now()
	session, err := domain.NewSession("session-1", []byte("hash"), "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewSessionRepository(userRepositoryExecutor{execErr: &pgconn.PgError{Code: uniqueViolationCode}})

	if err := repo.CreateSession(t.Context(), session); !errors.Is(err, application.ErrSessionAlreadyExists) {
		t.Fatalf("expected already-exists error, got %v", err)
	}
}
