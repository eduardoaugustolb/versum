package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	identityports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/jackc/pgx/v5"
)

type userRepositoryExecutor struct {
	row     dbexec.Row
	execErr error
}

func (e userRepositoryExecutor) QueryRow(context.Context, string, ...any) dbexec.Row {
	return e.row
}

func (userRepositoryExecutor) Query(context.Context, string, ...any) (dbexec.Rows, error) {
	return nil, nil
}

func (e userRepositoryExecutor) Exec(context.Context, string, ...any) error { return e.execErr }

func (userRepositoryExecutor) CopyFrom(context.Context, string, []string, [][]any) (int64, error) {
	return 0, nil
}

func (userRepositoryExecutor) Begin(context.Context) (dbexec.Transaction, error) {
	return nil, nil
}

type errorRow struct{ err error }

func (r errorRow) Scan(...any) error { return r.err }

type userRow struct {
	id                string
	ciphertext        []byte
	lookupHMAC        []byte
	encryptionVersion int
	lookupVersion     int
}

func (r userRow) Scan(destinations ...any) error {
	*destinations[0].(*string) = r.id
	*destinations[1].(*[]byte) = r.ciphertext
	*destinations[2].(*[]byte) = r.lookupHMAC
	*destinations[3].(*int) = r.encryptionVersion
	*destinations[4].(*int) = r.lookupVersion
	return nil
}

type lookupExecutor struct {
	rows      []dbexec.Row
	queryArgs [][]any
	execQuery []string
	execArgs  [][]any
	execErr   error
}

func (e *lookupExecutor) QueryRow(_ context.Context, _ string, args ...any) dbexec.Row {
	e.queryArgs = append(e.queryArgs, args)
	rowIndex := len(e.queryArgs) - 1
	if rowIndex >= len(e.rows) {
		return errorRow{err: pgx.ErrNoRows}
	}
	return e.rows[rowIndex]
}

func (e *lookupExecutor) Query(context.Context, string, ...any) (dbexec.Rows, error) { return nil, nil }

func (e *lookupExecutor) Exec(_ context.Context, query string, args ...any) error {
	e.execQuery = append(e.execQuery, query)
	e.execArgs = append(e.execArgs, args)
	return e.execErr
}

func (*lookupExecutor) CopyFrom(context.Context, string, []string, [][]any) (int64, error) {
	return 0, nil
}

func (*lookupExecutor) Begin(context.Context) (dbexec.Transaction, error) { return nil, nil }

type lookupCandidateProtector struct {
	candidates []identityports.LookupCandidate
	lookupErr  error
}

func (p lookupCandidateProtector) Protect(context.Context, domain.Email) (identityports.ProtectedEmail, error) {
	return identityports.ProtectedEmail{}, nil
}

func (p lookupCandidateProtector) Unprotect(context.Context, identityports.ProtectedEmail) (domain.Email, error) {
	return domain.ParseEmail("user@example.com")
}

func (p lookupCandidateProtector) LookupHMAC(context.Context, domain.Email) ([]byte, error) {
	return nil, errors.New("FindUserByEmail must use LookupCandidates")
}

func (p lookupCandidateProtector) LookupCandidates(context.Context, domain.Email) ([]identityports.LookupCandidate, error) {
	if p.lookupErr != nil {
		return nil, p.lookupErr
	}
	return p.candidates, nil
}

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

func (userRepositoryProtector) LookupCandidates(context.Context, domain.Email) ([]identityports.LookupCandidate, error) {
	return []identityports.LookupCandidate{{LookupKeyVersion: 1, LookupHMAC: []byte("lookup")}}, nil
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

func TestUserRepositoryReturnsAlreadyExistsWhenInsertDoesNotReturnUser(t *testing.T) {
	email, err := domain.ParseEmail("user@example.com")
	if err != nil {
		t.Fatal(err)
	}
	user, err := domain.NewUser("user-1", email.String())
	if err != nil {
		t.Fatal(err)
	}
	repo := NewUserRepository(userRepositoryExecutor{row: errorRow{err: pgx.ErrNoRows}}, userRepositoryProtector{})

	if err := repo.CreateUser(t.Context(), user); !errors.Is(err, application.ErrUserAlreadyExists) {
		t.Fatalf("expected already-exists error, got %v", err)
	}
}

func TestUserRepositoryFindUserByEmailFallsBackToPreviousLookupKey(t *testing.T) {
	executor := &lookupExecutor{rows: []dbexec.Row{
		userRow{id: "user-1", ciphertext: []byte("ciphertext"), lookupHMAC: []byte("old-hmac"), encryptionVersion: 1, lookupVersion: 1},
	}}
	repo := NewUserRepository(executor, lookupCandidateProtector{candidates: []identityports.LookupCandidate{
		{LookupKeyVersion: 2, LookupHMAC: []byte("current-hmac")},
		{LookupKeyVersion: 1, LookupHMAC: []byte("old-hmac")},
	}})
	email, err := domain.ParseEmail("user@example.com")
	if err != nil {
		t.Fatal(err)
	}

	user, err := repo.FindUserByEmail(t.Context(), email)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID() != "user-1" {
		t.Fatalf("expected user-1, got %q", user.ID())
	}
	if len(executor.queryArgs) != 1 {
		t.Fatalf("expected one lookup query, got %d", len(executor.queryArgs))
	}
	lookupHMACs := executor.queryArgs[0][0].([][]byte)
	lookupKeyVersions := executor.queryArgs[0][1].([]int16)
	if string(lookupHMACs[0]) != "current-hmac" || lookupKeyVersions[0] != 2 || string(lookupHMACs[1]) != "old-hmac" || lookupKeyVersions[1] != 1 {
		t.Fatalf("unexpected lookup candidates: hmacs=%q versions=%v", lookupHMACs, lookupKeyVersions)
	}
}

func TestUserRepositoryFindUserByEmailStopsOnDatabaseError(t *testing.T) {
	executor := &lookupExecutor{rows: []dbexec.Row{errorRow{err: errDatabaseUnavailable}}}
	repo := NewUserRepository(executor, lookupCandidateProtector{candidates: []identityports.LookupCandidate{
		{LookupKeyVersion: 2, LookupHMAC: []byte("current-hmac")},
		{LookupKeyVersion: 1, LookupHMAC: []byte("old-hmac")},
	}})
	email, err := domain.ParseEmail("user@example.com")
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.FindUserByEmail(t.Context(), email)
	if !errors.Is(err, errDatabaseUnavailable) {
		t.Fatalf("expected database error, got %v", err)
	}
	if len(executor.queryArgs) != 1 {
		t.Fatalf("expected one lookup query after database error, got %d", len(executor.queryArgs))
	}
}

func TestUserRepositoryFindUserByEmailMigratesPreviousLookupKey(t *testing.T) {
	executor := &lookupExecutor{rows: []dbexec.Row{
		userRow{id: "user-1", ciphertext: []byte("ciphertext"), lookupHMAC: []byte("old-hmac"), encryptionVersion: 1, lookupVersion: 1},
	}}
	repo := NewUserRepository(executor, lookupCandidateProtector{candidates: []identityports.LookupCandidate{
		{LookupKeyVersion: 2, LookupHMAC: []byte("current-hmac")},
		{LookupKeyVersion: 1, LookupHMAC: []byte("old-hmac")},
	}})
	email, err := domain.ParseEmail("user@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := repo.FindUserByEmail(t.Context(), email); err != nil {
		t.Fatal(err)
	}
	if len(executor.execQuery) != 1 || !strings.Contains(executor.execQuery[0], "UPDATE users") {
		t.Fatalf("expected a lookup-key migration update, got queries %#v", executor.execQuery)
	}
	if got := executor.execArgs[0]; len(got) != 4 || got[0] != "user-1" || string(got[1].([]byte)) != "current-hmac" || got[2] != 2 || got[3] != 1 {
		t.Fatalf("unexpected migration args: %#v", got)
	}
}
