package postgres

import (
	"context"
	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
	"time"
)

type EventRecord struct {
	ID, EventType, UserID string
	PayloadCiphertext     []byte
	PayloadKeyVersion     int
}

type DeliveryRecord struct {
	EventRecord
	Attempt int
}

// Claim must use an autocommit executor: its lock must end before delivery starts.
func (d *DAO) Claim(ctx context.Context, now time.Time, leaseToken string, leaseUntil time.Time) (DeliveryRecord, error) {
	var record DeliveryRecord
	err := d.db.QueryRow(ctx, ClaimEventQuery, now, leaseToken, leaseUntil).Scan(
		&record.ID, &record.EventType, &record.UserID, &record.PayloadCiphertext, &record.PayloadKeyVersion, &record.Attempt,
	)
	return record, err
}

func (d *DAO) Complete(ctx context.Context, eventID, leaseToken string, now time.Time) error {
	var id string
	return d.db.QueryRow(ctx, CompleteEventQuery, eventID, leaseToken, now).Scan(&id)
}

func (d *DAO) Retry(ctx context.Context, eventID, leaseToken string, now, availableAt time.Time) error {
	var id string
	return d.db.QueryRow(ctx, RetryEventQuery, eventID, leaseToken, now, availableAt).Scan(&id)
}

func (d *DAO) Fail(ctx context.Context, eventID, leaseToken string, now time.Time) error {
	var id string
	return d.db.QueryRow(ctx, FailEventQuery, eventID, leaseToken, now).Scan(&id)
}

type DAO struct{ db dbexec.Executor }

func NewDAO(db dbexec.Executor) *DAO { return &DAO{db: db} }
func (d *DAO) Insert(ctx context.Context, record EventRecord) error {
	return d.db.Exec(ctx, PublishEventQuery, record.ID, record.EventType, record.PayloadCiphertext, record.PayloadKeyVersion, record.UserID)
}
