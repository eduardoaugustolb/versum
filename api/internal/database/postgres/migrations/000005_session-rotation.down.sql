DROP INDEX sessions_family_id_idx;
ALTER TABLE sessions DROP CONSTRAINT sessions_user_agent_length;
ALTER TABLE sessions DROP CONSTRAINT sessions_replacement_state;
ALTER TABLE sessions DROP COLUMN replaced_by_session_id;
ALTER TABLE sessions DROP COLUMN user_agent;
ALTER TABLE sessions DROP COLUMN ip_address;
ALTER TABLE sessions DROP COLUMN family_id;
ALTER INDEX sessions_last_used_idx RENAME TO sessions_used_idx;
ALTER TABLE sessions RENAME COLUMN last_used_at TO used_at;
