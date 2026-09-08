package clock

import "time"

// SystemClock returns the current instant normalized to UTC.
type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
