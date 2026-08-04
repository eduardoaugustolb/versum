package health_test

import (
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/health"
)

func TestCheckHealthExecute(t *testing.T) {
	got := health.CheckHealth{}.Execute()

	want := health.Status{State: "ok"}
	if got != want {
		t.Errorf("expected %+v, got %+v", want, got)
	}
}
