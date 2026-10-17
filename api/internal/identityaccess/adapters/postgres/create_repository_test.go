package postgres_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/adapters/postgres"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
)

const uniqueViolationCode = "23505"

func TestLoginTokenRepositoryTranslatesUniqueViolationOnCreate(t *testing.T) {
	now := time.Now()
	token, err := domain.NewLoginToken("token-1", "user-1", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	repo := postgres.NewLoginTokenRepository(userRepositoryExecutor{execErr: &pgconn.PgError{Code: uniqueViolationCode}})

	if err := repo.Create(t.Context(), token, "secret"); !errors.Is(err, application.ErrLoginTokenAlreadyExists) {
		t.Fatalf("expected already-exists error, got %v", err)
	}
}

func TestSessionRepositoryTranslatesUniqueViolationOnCreate(t *testing.T) {
	now := time.Now()
	session, err := domain.NewSession("session-1", "user-1", "family-1", "", "", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	repo := postgres.NewSessionRepository(userRepositoryExecutor{execErr: &pgconn.PgError{Code: uniqueViolationCode}})

	if err := repo.CreateSession(t.Context(), session, "secret"); !errors.Is(err, application.ErrSessionAlreadyExists) {
		t.Fatalf("expected already-exists error, got %v", err)
	}
}
