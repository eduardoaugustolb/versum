package postgres

import (
	"context"
	dbexec "github.com/eduardoaugustolb/versum/api/internal/database"
)

type EventRecord struct {
	ID, EventType, UserID string
	PayloadCiphertext     []byte
	PayloadKeyVersion     int
}
type DAO struct{ db dbexec.Executor }

func NewDAO(db dbexec.Executor) *DAO { return &DAO{db: db} }
func (d *DAO) Insert(ctx context.Context, record EventRecord) error {
	return d.db.Exec(ctx, PublishEventQuery, record.ID, record.EventType, record.PayloadCiphertext, record.PayloadKeyVersion, record.UserID)
}
