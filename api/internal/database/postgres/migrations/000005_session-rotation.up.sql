ALTER TABLE sessions RENAME COLUMN used_at TO last_used_at;
ALTER INDEX sessions_used_idx RENAME TO sessions_last_used_idx;
ALTER TABLE sessions ADD COLUMN family_id TEXT;
-- Existing sessions start independent families; their old rotation history is unknown.
UPDATE sessions SET family_id = id;
ALTER TABLE sessions ALTER COLUMN family_id SET NOT NULL;
ALTER TABLE sessions ADD CONSTRAINT sessions_family_id_nonempty CHECK (length(trim(family_id)) > 0);
ALTER TABLE sessions ADD COLUMN replaced_by_session_id TEXT REFERENCES sessions(id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE sessions ADD COLUMN ip_address TEXT NOT NULL DEFAULT '';
ALTER TABLE sessions ADD COLUMN user_agent TEXT NOT NULL DEFAULT '';
ALTER TABLE sessions ADD CONSTRAINT sessions_replacement_state CHECK (replaced_by_session_id IS NULL OR (replaced_by_session_id <> id AND revoked_at IS NOT NULL));
ALTER TABLE sessions ADD CONSTRAINT sessions_user_agent_length CHECK (octet_length(user_agent) <= 1024);
CREATE INDEX sessions_family_id_idx ON sessions (family_id);
