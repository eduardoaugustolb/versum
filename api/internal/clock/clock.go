package clock

import "time"

// Clock provides the current instant to application use cases.
type Clock interface {
	Now() time.Time
}
