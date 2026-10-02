package postgres

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"testing"
	"time"

	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"github.com/eduardoaugustolb/versum/api/internal/identityaccess/domain"
	"github.com/jackc/pgx/v5"
)

type sessionRow struct {
	familyID, ipAddress, userAgent string
	replacement                    *string
	id, userID                     string
	revokedAt, usedAt              *time.Time
	expiresAt                      time.Time
}

func (r sessionRow) Scan(dest ...any) error {
	*dest[0].(*string) = r.id
	*dest[1].(*string) = r.userID
	*dest[2].(**time.Time) = r.revokedAt
	*dest[3].(**time.Time) = r.usedAt
	*dest[4].(*time.Time) = r.expiresAt
	familyID := r.familyID
	if familyID == "" {
		familyID = "family-1"
	}
	*dest[5].(*string) = familyID
	*dest[6].(*string) = r.ipAddress
	*dest[7].(*string) = r.userAgent
	*dest[8].(**string) = r.replacement
	return nil
}

type sessionRows struct {
	rows   []dbexec.Row
	index  int
	err    error
	closed bool
}

func (r *sessionRows) Next() bool {
	if r.index >= len(r.rows) {
		return false
	}
	r.index++
	return true
}
func (r *sessionRows) Scan(dest ...any) error { return r.rows[r.index-1].Scan(dest...) }
func (r *sessionRows) Err() error             { return r.err }
func (r *sessionRows) Close()                 { r.closed = true }

type sessionExecutor struct {
	lookupExecutor
	result   dbexec.Rows
	queryErr error
	listArgs []any
}

func (e *sessionExecutor) Query(_ context.Context, _ string, args ...any) (dbexec.Rows, error) {
	e.listArgs = args
	return e.result, e.queryErr
}

func TestSessionRepositoryCreatesAndFindsWithSecret(t *testing.T) {
	now := time.Now().UTC()
	session, err := domain.NewSession("session-1", "user-1", "family-1", "", "", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	db := &lookupExecutor{rows: []dbexec.Row{
		sessionRow{id: session.ID(), userID: session.UserID(), expiresAt: session.ExpiresAt()},
		sessionRow{id: session.ID(), userID: session.UserID(), expiresAt: session.ExpiresAt()},
	}}
	repo := NewSessionRepository(db)
	if err := repo.CreateSession(t.Context(), session, "secret"); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256([]byte("secret"))
	args := db.execArgs[0]
	if !bytes.Equal(args[1].([]byte), digest[:]) {
		t.Fatal("must persist the digest, not the secret")
	}
	if args[3] != nil || args[4] != nil {
		t.Fatal("unset timestamps must persist as NULL")
	}
	for _, find := range []func() (*domain.Session, error){
		func() (*domain.Session, error) { return repo.FindSessionByID(t.Context(), session.ID()) },
		func() (*domain.Session, error) { return repo.FindSessionBySecret(t.Context(), "secret") },
	} {
		got, err := find()
		if err != nil {
			t.Fatal(err)
		}
		if got.ID() != session.ID() || got.UserID() != session.UserID() || !got.ExpiresAt().Equal(session.ExpiresAt()) {
			t.Fatal("incorrect rehydrated session")
		}
	}
	if !bytes.Equal(db.queryArgs[1][0].([]byte), digest[:]) {
		t.Fatal("lookup must use the same digest as creation")
	}
}

func TestSessionRepositoryPersistsAndRehydratesOptionalState(t *testing.T) {
	now := time.Now().UTC()
	session, err := domain.NewSession("session-1", "user-1", "family-1", "", "", now.Add(time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := session.Use(now); err != nil {
		t.Fatal(err)
	}
	session.Revoke(now.Add(time.Minute))
	revokedAt, _ := session.RevokedAt()
	usedAt, _ := session.LastUsedAt()
	db := &lookupExecutor{rows: []dbexec.Row{sessionRow{id: session.ID(), userID: session.UserID(), revokedAt: &revokedAt, usedAt: &usedAt, expiresAt: session.ExpiresAt()}}}
	repo := NewSessionRepository(db)
	if err := repo.CreateSession(t.Context(), session, "secret"); err != nil {
		t.Fatal(err)
	}
	if db.execArgs[0][3] != revokedAt || db.execArgs[0][4] != usedAt {
		t.Fatal("state timestamps were not persisted")
	}
	got, err := repo.FindSessionByID(t.Context(), session.ID())
	if err != nil {
		t.Fatal(err)
	}
	if at, ok := got.RevokedAt(); !ok || !at.Equal(revokedAt) {
		t.Fatal("revocation was not rehydrated")
	}
	if at, ok := got.LastUsedAt(); !ok || !at.Equal(usedAt) {
		t.Fatal("usage was not rehydrated")
	}
}

func TestSessionRepositoryRejectsInvalidArguments(t *testing.T) {
	zero := time.Time{}
	now := time.Now()
	repo := NewSessionRepository(&lookupExecutor{})
	tests := []struct {
		name string
		run  func() error
		want error
	}{
		{"create empty secret", func() error { return repo.CreateSession(t.Context(), nil, "") }, application.ErrInvalidSessionSecret},
		{"find empty secret", func() error { _, err := repo.FindSessionBySecret(t.Context(), ""); return err }, application.ErrInvalidSessionSecret},
		{"find empty id", func() error { _, err := repo.FindSessionByID(t.Context(), ""); return err }, domain.ErrInvalidSessionID},
		{"list empty user", func() error { _, err := repo.ListSessionsByUserID(t.Context(), ""); return err }, domain.ErrInvalidSessionUserID},
		{"revoke empty id", func() error { return repo.RevokeSession(t.Context(), "", &now) }, domain.ErrInvalidSessionID},
		{"revoke nil time", func() error { return repo.RevokeSession(t.Context(), "session", nil) }, domain.ErrInvalidSessionRevokedAt},
		{"revoke zero time", func() error { return repo.RevokeSession(t.Context(), "session", &zero) }, domain.ErrInvalidSessionRevokedAt},
		{"revoke all empty user", func() error { return repo.RevokeAllSessions(t.Context(), "", &now) }, domain.ErrInvalidSessionUserID},
		{"revoke all nil time", func() error { return repo.RevokeAllSessions(t.Context(), "user", nil) }, domain.ErrInvalidSessionRevokedAt},
		{"revoke all zero time", func() error { return repo.RevokeAllSessions(t.Context(), "user", &zero) }, domain.ErrInvalidSessionRevokedAt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); !errors.Is(err, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, err)
			}
		})
	}
}

func TestSessionRepositoryFindErrors(t *testing.T) {
	failure := errors.New("database unavailable")
	for _, tt := range []struct {
		name string
		row  dbexec.Row
		want error
	}{
		{"not found", errorRow{err: pgx.ErrNoRows}, application.ErrSessionNotFound},
		{"database failure", errorRow{err: failure}, failure},
		{"invalid persisted state", sessionRow{userID: "user", expiresAt: time.Now()}, domain.ErrInvalidSessionID},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, bySecret := range []bool{false, true} {
				repo := NewSessionRepository(userRepositoryExecutor{row: tt.row})
				var err error
				if bySecret {
					_, err = repo.FindSessionBySecret(t.Context(), "secret")
				} else {
					_, err = repo.FindSessionByID(t.Context(), "session")
				}
				if !errors.Is(err, tt.want) {
					t.Fatalf("want %v, got %v", tt.want, err)
				}
			}
		})
	}
}

func TestSessionRepositoryListsAndClosesRows(t *testing.T) {
	now := time.Now()
	rows := &sessionRows{rows: []dbexec.Row{sessionRow{id: "first", userID: "user", expiresAt: now}, sessionRow{id: "second", userID: "user", expiresAt: now}}}
	db := &sessionExecutor{result: rows}
	sessions, err := NewSessionRepository(db).ListSessionsByUserID(t.Context(), "user")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 || sessions[0].ID() != "first" || sessions[1].ID() != "second" {
		t.Fatal("incorrect list order or contents")
	}
	if !rows.closed || db.listArgs[0] != "user" {
		t.Fatal("must filter by user and close rows")
	}
}

func TestSessionRepositoryListErrorsAndEmptyResult(t *testing.T) {
	failure := errors.New("database unavailable")
	for _, tt := range []struct {
		name           string
		rows           *sessionRows
		queryErr, want error
	}{
		{"empty", &sessionRows{}, nil, nil},
		{"query failure", nil, failure, failure},
		{"scan failure", &sessionRows{rows: []dbexec.Row{errorRow{err: failure}}}, nil, failure},
		{"iteration failure", &sessionRows{err: failure}, nil, failure},
	} {
		t.Run(tt.name, func(t *testing.T) {
			db := &sessionExecutor{result: tt.rows, queryErr: tt.queryErr}
			sessions, err := NewSessionRepository(db).ListSessionsByUserID(t.Context(), "user")
			if !errors.Is(err, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, err)
			}
			if tt.want == nil && (sessions == nil || len(sessions) != 0) {
				t.Fatal("must return an empty list")
			}
			if tt.rows != nil && !tt.rows.closed {
				t.Fatal("rows must close even on failure")
			}
		})
	}
}

func TestSessionRepositoryRevokesAndPreservesWriteErrors(t *testing.T) {
	now := time.Now()
	failure := errors.New("database unavailable")
	for _, writeErr := range []error{nil, failure} {
		db := &lookupExecutor{execErr: writeErr}
		repo := NewSessionRepository(db)
		session, err := domain.NewSession("session", "user", "family-1", "", "", now.Add(time.Hour), now)
		if err != nil {
			t.Fatal(err)
		}
		for _, run := range []func() error{
			func() error { return repo.CreateSession(t.Context(), session, "secret") },
			func() error { return repo.RevokeSession(t.Context(), "session", &now) },
			func() error { return repo.RevokeAllSessions(t.Context(), "user", &now) },
		} {
			if err := run(); !errors.Is(err, writeErr) {
				t.Fatalf("want %v, got %v", writeErr, err)
			}
		}
		if db.execQuery[1] != RevokeSessionQuery || db.execArgs[1][0] != "session" || db.execArgs[1][1] != now {
			t.Fatal("incorrect session revocation")
		}
		if db.execQuery[2] != RevokeAllSessionsQuery || db.execArgs[2][0] != "user" || db.execArgs[2][1] != now {
			t.Fatal("incorrect bulk revocation")
		}
	}
}

type sessionIDRow struct{ id string }

func (r sessionIDRow) Scan(dest ...any) error { *dest[0].(*string) = r.id; return nil }

func TestSessionRepositoryRotationPersistenceAndErrors(t *testing.T) {
	now := time.Now().UTC()
	old, _ := domain.NewSession("old", "user", "family", "", "", now.Add(time.Hour), now)
	next, _ := domain.NewSession("next", "user", "family", "", "", old.ExpiresAt(), now)
	if err := old.ReplaceWith(next, now); err != nil {
		t.Fatal(err)
	}
	failure := errors.New("database failure")
	for _, tt := range []struct {
		name               string
		row                dbexec.Row
		lockWant, saveWant error
	}{
		{"success", sessionIDRow{id: "old"}, nil, nil},
		{"missing", errorRow{err: pgx.ErrNoRows}, application.ErrUserNotFound, domain.ErrInvalidSessionState},
		{"failure", errorRow{err: failure}, failure, failure},
	} {
		t.Run(tt.name, func(t *testing.T) {
			db := &lookupExecutor{rows: []dbexec.Row{tt.row, tt.row}}
			repo := NewSessionRepository(db)
			if err := repo.LockSessionsByUserID(t.Context(), "user"); !errors.Is(err, tt.lockWant) {
				t.Fatal(err)
			}
			if err := repo.SaveRotation(t.Context(), old); !errors.Is(err, tt.saveWant) {
				t.Fatal(err)
			}
			args := db.queryArgs[1]
			if args[0] != "old" || args[1] != now || args[2] != now || args[3] != "next" {
				t.Fatal("incorrect rotation persistence")
			}
		})
	}
	repo := NewSessionRepository(&lookupExecutor{})
	if err := repo.LockSessionsByUserID(t.Context(), ""); !errors.Is(err, domain.ErrInvalidSessionUserID) {
		t.Fatal(err)
	}
	if err := repo.SaveRotation(t.Context(), next); !errors.Is(err, domain.ErrInvalidSessionReplacement) {
		t.Fatal(err)
	}
}

func TestSessionRepositoryRevokesOnlyRequestedFamily(t *testing.T) {
	now := time.Now().UTC()
	failure := errors.New("database failure")
	for _, writeErr := range []error{nil, failure} {
		db := &lookupExecutor{execErr: writeErr}
		if err := NewSessionRepository(db).RevokeSessionFamily(t.Context(), "user", "family", now); !errors.Is(err, writeErr) {
			t.Fatal(err)
		}
		if db.execQuery[0] != RevokeSessionFamilyQuery || db.execArgs[0][0] != "user" || db.execArgs[0][1] != "family" || db.execArgs[0][2] != now {
			t.Fatal("must scope revocation by user and family")
		}
	}
	repo := NewSessionRepository(&lookupExecutor{})
	for _, tt := range []struct {
		user, family string
		at           time.Time
		want         error
	}{
		{"", "family", now, domain.ErrInvalidSessionUserID}, {"user", "", now, domain.ErrInvalidSessionFamilyID}, {"user", "family", time.Time{}, domain.ErrInvalidSessionRevokedAt},
	} {
		if err := repo.RevokeSessionFamily(t.Context(), tt.user, tt.family, tt.at); !errors.Is(err, tt.want) {
			t.Fatal(err)
		}
	}
}

func TestSessionRepositoryRehydratesRotationAndClient(t *testing.T) {
	now := time.Now().UTC()
	replacement := "next"
	db := &lookupExecutor{rows: []dbexec.Row{sessionRow{id: "old", userID: "user", familyID: "family", ipAddress: "2001:db8::1", userAgent: "browser", revokedAt: &now, usedAt: &now, expiresAt: now.Add(time.Hour), replacement: &replacement}}}
	got, err := NewSessionRepository(db).FindSessionByID(t.Context(), "old")
	if err != nil {
		t.Fatal(err)
	}
	if id, ok := got.ReplacedBySessionID(); !ok || id != "next" {
		t.Fatal("replacement not rehydrated")
	}
	if got.FamilyID() != "family" || got.IPAddress() != "2001:db8::1" || got.UserAgent() != "browser" || got.IsValidAt(now) {
		t.Fatal("incorrect persisted state")
	}
}

func TestSessionRepositorySavesUsageConditionally(t *testing.T) {
	now := time.Now().UTC()
	session, _ := domain.NewSession("session", "user", "family", "", "", now.Add(time.Hour), now)
	repo := NewSessionRepository(&lookupExecutor{})
	if err := repo.SaveUsage(t.Context(), session); !errors.Is(err, domain.ErrInvalidSessionLastUsedAt) {
		t.Fatal(err)
	}
	if err := session.Use(now); err != nil {
		t.Fatal(err)
	}
	failure := errors.New("database failure")
	for _, tt := range []struct {
		row  dbexec.Row
		want error
	}{{sessionIDRow{id: session.ID()}, nil}, {errorRow{err: pgx.ErrNoRows}, domain.ErrInvalidSessionState}, {errorRow{err: failure}, failure}} {
		db := &lookupExecutor{rows: []dbexec.Row{tt.row}}
		if err := NewSessionRepository(db).SaveUsage(t.Context(), session); !errors.Is(err, tt.want) {
			t.Fatal(err)
		}
		if db.queryArgs[0][0] != session.ID() || db.queryArgs[0][1] != now {
			t.Fatal("incorrect usage persistence")
		}
	}
	session.Revoke(now)
	if err := repo.SaveUsage(t.Context(), session); !errors.Is(err, domain.ErrInvalidSessionState) {
		t.Fatal(err)
	}
}

func TestSessionRepositoryRejectsCreationOfReplacedSession(t *testing.T) {
	now := time.Now().UTC()
	session, _ := domain.RehydrateSession("old", "user", "family", "", "", &now, &now, now.Add(time.Hour), "next")
	if err := NewSessionRepository(&lookupExecutor{}).CreateSession(t.Context(), session, "secret"); !errors.Is(err, domain.ErrInvalidSessionReplacement) {
		t.Fatal(err)
	}
}
