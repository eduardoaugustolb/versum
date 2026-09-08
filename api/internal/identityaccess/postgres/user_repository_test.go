package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identityports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/ports"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
)

type userRepositoryExecutor struct{ row dbexec.Row }

func (e userRepositoryExecutor) QueryRow(context.Context, string, ...any) dbexec.Row {
	return e.row
}

func (userRepositoryExecutor) Query(context.Context, string, ...any) (dbexec.Rows, error) {
	return nil, nil
}

func (userRepositoryExecutor) Exec(context.Context, string, ...any) error { return nil }

func (userRepositoryExecutor) CopyFrom(context.Context, string, []string, [][]any) (int64, error) {
	return 0, nil
}

type errorRow struct{ err error }

func (r errorRow) Scan(...any) error { return r.err }

type userRepositoryProtector struct{}

func (userRepositoryProtector) Protect(context.Context, domain.Email) (identityports.ProtectedEmail, error) {
	return identityports.ProtectedEmail{}, nil
}

func (userRepositoryProtector) Unprotect(context.Context, identityports.ProtectedEmail) (domain.Email, error) {
	return domain.ParseEmail("user@example.com")
}

func (userRepositoryProtector) LookupHMAC(context.Context, domain.Email) ([]byte, error) {
	return []byte("lookup"), nil
}

var errDatabaseUnavailable = errors.New("database unavailable")

func TestUserRepositoryTranslatesOnlyMissingRows(t *testing.T) {
	tests := []struct {
		name    string
		rowErr  error
		wantErr error
	}{
		{
			name:    "missing row",
			rowErr:  pgx.ErrNoRows,
			wantErr: application.ErrUserNotFound,
		},
		{
			name:    "database failure",
			rowErr:  errDatabaseUnavailable,
			wantErr: errDatabaseUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewUserRepository(userRepositoryExecutor{row: errorRow{err: tt.rowErr}}, userRepositoryProtector{})

			_, err := repo.FindUserByID(t.Context(), "user-1")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error matching %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestUserRepositoryRejectsInvalidLookupValues(t *testing.T) {
	repo := NewUserRepository(userRepositoryExecutor{}, userRepositoryProtector{})

	if _, err := repo.FindUserByID(t.Context(), ""); !errors.Is(err, domain.ErrInvalidUserID) {
		t.Fatalf("expected invalid user id, got %v", err)
	}
}
