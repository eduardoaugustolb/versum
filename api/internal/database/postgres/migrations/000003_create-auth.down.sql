DROP INDEX IF EXISTS login_tokens_user_created_idx;
DROP INDEX IF EXISTS login_tokens_expires_idx;
DROP TABLE IF EXISTS login_tokens;

DROP INDEX IF EXISTS sessions_user_created_idx;
DROP INDEX IF EXISTS sessions_used_idx;
DROP INDEX IF EXISTS sessions_expires_idx;
DROP TABLE IF EXISTS sessions;

DROP TABLE IF EXISTS users;
