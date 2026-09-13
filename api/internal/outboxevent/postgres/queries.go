package postgres

const (
	PublishEventQuery = `
		INSERT INTO outbox_events (
			id,
			event_type,
			payload_ciphertext,
			payload_key_version,
			available_at,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
)
