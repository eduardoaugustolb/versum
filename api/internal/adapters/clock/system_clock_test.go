package clock_test

import (
	"testing"
	"time"

	"github.com/eduardoaugustolb/versum/api/internal/adapters/clock"
)

func TestSystemClockNowIsUTC(t *testing.T) {
	if got := (clock.SystemClock{}).Now(); got.Location() != time.UTC {
		t.Fatalf("expected UTC instant, got %s", got.Location())
	}
}
