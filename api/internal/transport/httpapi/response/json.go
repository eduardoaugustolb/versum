// Package response contains HTTP response serialization shared by adapters.
package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// WriteJSON serializes v as a JSON HTTP response.
func WriteJSON(w http.ResponseWriter, request *http.Request, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		slog.ErrorContext(request.Context(), "failed to marshal HTTP response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(body); err != nil {
		slog.ErrorContext(request.Context(), "failed to write HTTP response", "error", err)
	}
}
