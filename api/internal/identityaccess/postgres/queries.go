package postgres

const (
	FindUserByIDQuery                 = `SELECT id, email_ciphertext, email_lookup_hmac, email_encryption_key_version, email_lookup_key_version FROM users WHERE id = $1`
	FindUserByEmailLookupHMACQuery    = `SELECT id, email_ciphertext, email_lookup_hmac, email_encryption_key_version, email_lookup_key_version FROM users WHERE email_lookup_hmac = $1`
	CreateUserQuery                   = `INSERT INTO users (id, email_ciphertext, email_lookup_hmac, email_encryption_key_version, email_lookup_key_version) VALUES ($1, $2, $3, $4, $5)`
	CreateSessionQuery                = `INSERT INTO sessions (id, secret_hash, user_id, revoked_at, used_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6)`
	FindSessionByIDQuery              = `SELECT id, secret_hash, user_id, revoked_at, used_at, expires_at FROM sessions WHERE id = $1`
	FindSessionBySecretHashQuery      = `SELECT id, secret_hash, user_id, revoked_at, used_at, expires_at FROM sessions WHERE secret_hash = $1`
	ListSessionsByUserIDQuery         = `SELECT id, secret_hash, user_id, revoked_at, used_at, expires_at FROM sessions WHERE user_id = $1 ORDER BY created_at ASC, id ASC`
	RevokeSessionQuery                = `UPDATE sessions SET revoked_at = $2 WHERE id = $1`
	RevokeAllSessionsQuery            = `UPDATE sessions SET revoked_at = $2 WHERE user_id = $1 AND revoked_at IS NULL`
	CreateLoginTokenQuery             = `INSERT INTO login_tokens (id, token_hash, user_id, expires_at, consumed_at) VALUES ($1, $2, $3, $4, $5)`
	FindLoginTokenByIDQuery           = `SELECT id, token_hash, user_id, expires_at, consumed_at FROM login_tokens WHERE id = $1`
	FindLoginTokenByTokenHashQuery    = `SELECT id, token_hash, user_id, expires_at, consumed_at FROM login_tokens WHERE token_hash = $1`
	ConsumeLoginTokenByTokenHashQuery = `UPDATE login_tokens SET consumed_at = $2 WHERE token_hash = $1`
)
