package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
)

func TestHealthEndpoint(t *testing.T) {
	router := httpapi.NewRouter(health.CheckHealth{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", got)
	}

	want := `{"status":"ok"}`
	if got := rec.Body.String(); got != want {
		t.Errorf("expected body %q, got %q", want, got)
	}
}
