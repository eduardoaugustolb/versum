package domain

import "time"

func cloneBytes(b []byte) []byte {
	return append([]byte(nil), b...)
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := value.UTC()
	return &copied
}
