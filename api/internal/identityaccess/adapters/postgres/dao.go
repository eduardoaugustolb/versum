package postgres

import (
	"context"
	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	ports "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application"
	"time"
)

type UserRecord struct {
	ID    string
	Email ports.ProtectedEmail
}
type LoginTokenRecord struct {
	ID, UserID string
	ExpiresAt  time.Time
	ConsumedAt *time.Time
}
type SessionRecord struct {
	ID, UserID, FamilyID, IPAddress, UserAgent string
	RevokedAt, LastUsedAt                      *time.Time
	ExpiresAt                                  time.Time
	Replacement                                *string
}
type DAO struct{ db dbexec.Executor }

func NewDAO(db dbexec.Executor) *DAO { return &DAO{db: db} }
func scanUserRecord(row dbexec.Row) (UserRecord, error) {
	var v UserRecord
	err := row.Scan(&v.ID, &v.Email.Ciphertext, &v.Email.LookupHMAC, &v.Email.EncryptionKeyVersion, &v.Email.LookupKeyVersion)
	return v, err
}
func scanLoginTokenRecord(row dbexec.Row) (LoginTokenRecord, error) {
	var v LoginTokenRecord
	err := row.Scan(&v.ID, &v.UserID, &v.ExpiresAt, &v.ConsumedAt)
	return v, err
}
func scanSessionRecord(row dbexec.Row) (SessionRecord, error) {
	var v SessionRecord
	err := row.Scan(&v.ID, &v.UserID, &v.RevokedAt, &v.LastUsedAt, &v.ExpiresAt, &v.FamilyID, &v.IPAddress, &v.UserAgent, &v.Replacement)
	return v, err
}
func (d *DAO) InsertUser(ctx context.Context, id string, email ports.ProtectedEmail) error {
	var returnedID string
	return d.db.QueryRow(ctx, CreateUserQuery, id, email.Ciphertext, email.LookupHMAC, email.EncryptionKeyVersion, email.LookupKeyVersion).Scan(&returnedID)
}
func (d *DAO) UserByID(ctx context.Context, id string) (UserRecord, error) {
	return scanUserRecord(d.db.QueryRow(ctx, FindUserByIDQuery, id))
}
func (d *DAO) UserByLookupCandidates(ctx context.Context, hashes [][]byte, versions []int16) (UserRecord, error) {
	return scanUserRecord(d.db.QueryRow(ctx, FindUserByLookupCandidatesQuery, hashes, versions))
}
func (d *DAO) UpdateUserLookup(ctx context.Context, id string, hash []byte, version, previousVersion int) error {
	return d.db.Exec(ctx, UpdateUserLookupKeyQuery, id, hash, version, previousVersion)
}
func (d *DAO) DeleteUser(ctx context.Context, id string) error {
	return d.db.Exec(ctx, DeleteUserQuery, id)
}
func (d *DAO) InsertLoginToken(ctx context.Context, id string, hash []byte, userID string, expiresAt time.Time, consumedAt any) error {
	return d.db.Exec(ctx, CreateLoginTokenQuery, id, hash, userID, expiresAt, consumedAt)
}
func (d *DAO) LoginTokenByID(ctx context.Context, id string) (LoginTokenRecord, error) {
	return scanLoginTokenRecord(d.db.QueryRow(ctx, FindLoginTokenByIDQuery, id))
}
func (d *DAO) LoginTokenByHash(ctx context.Context, hash []byte) (LoginTokenRecord, error) {
	return scanLoginTokenRecord(d.db.QueryRow(ctx, FindLoginTokenByTokenHashQuery, hash))
}
func (d *DAO) ConsumeLoginToken(ctx context.Context, id string, at time.Time) error {
	var returnedID string
	return d.db.QueryRow(ctx, ConsumeLoginTokenByIDQuery, id, at).Scan(&returnedID)
}
func (d *DAO) InsertSession(ctx context.Context, id string, hash []byte, userID string, revokedAt, lastUsedAt any, expiresAt time.Time, familyID, ip, ua string) error {
	return d.db.Exec(ctx, CreateSessionQuery, id, hash, userID, revokedAt, lastUsedAt, expiresAt, familyID, ip, ua)
}
func (d *DAO) SessionByID(ctx context.Context, id string) (SessionRecord, error) {
	return scanSessionRecord(d.db.QueryRow(ctx, FindSessionByIDQuery, id))
}
func (d *DAO) SessionByHash(ctx context.Context, hash []byte) (SessionRecord, error) {
	return scanSessionRecord(d.db.QueryRow(ctx, FindSessionBySecretHashQuery, hash))
}
func (d *DAO) RevokeSession(ctx context.Context, id string, at time.Time) error {
	return d.db.Exec(ctx, RevokeSessionQuery, id, at)
}
func (d *DAO) RevokeAllSessions(ctx context.Context, userID string, at time.Time) error {
	return d.db.Exec(ctx, RevokeAllSessionsQuery, userID, at)
}
func (d *DAO) LockUser(ctx context.Context, userID string) error {
	var returnedID string
	return d.db.QueryRow(ctx, LockSessionsByUserIDQuery, userID).Scan(&returnedID)
}
func (d *DAO) RotateSession(ctx context.Context, id string, revokedAt, lastUsedAt time.Time, replacement string) error {
	var returnedID string
	return d.db.QueryRow(ctx, SaveSessionRotationQuery, id, revokedAt, lastUsedAt, replacement).Scan(&returnedID)
}
func (d *DAO) RevokeSessionFamily(ctx context.Context, userID, familyID string, at time.Time) error {
	return d.db.Exec(ctx, RevokeSessionFamilyQuery, userID, familyID, at)
}
func (d *DAO) UpdateSessionUsage(ctx context.Context, id string, at time.Time) error {
	var returnedID string
	return d.db.QueryRow(ctx, SaveSessionUsageQuery, id, at).Scan(&returnedID)
}
func (d *DAO) SessionsByUserID(ctx context.Context, userID string) ([]SessionRecord, error) {
	rows, err := d.db.Query(ctx, ListSessionsByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	records := []SessionRecord{}
	for rows.Next() {
		record, err := scanSessionRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
