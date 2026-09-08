// Package clock defines the application-facing source of current time.
package clock

import "time"

type Clock interface {
	Now() time.Time
}
