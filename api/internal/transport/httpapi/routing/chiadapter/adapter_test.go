package chiadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
)

func TestAdapterRoutesNestedRouter(t *testing.T) {
	adapter := New()
	adapter.Route("/v1", func(router routing.Router) {
		router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	})

	recorder := httptest.NewRecorder()
	adapter.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/health", nil))

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
}
