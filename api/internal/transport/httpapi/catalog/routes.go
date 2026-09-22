// Package catalog adapts catalog use cases to HTTP endpoints.
package catalog

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/response"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
)

// Dependencies are the catalog use cases required by this HTTP adapter.
type Dependencies struct {
	ListBooks  *queries.ListBooks
	GetChapter *queries.GetChapter
}

// RegisterRoutes registers catalog endpoints.
func RegisterRoutes(router routing.Router, deps Dependencies) {
	handler := handler{deps: deps}
	router.Route("/books", func(router routing.Router) {
		router.Get("/", handler.listBooks)
		router.Get("/{bookId}/chapters/{number}", handler.getChapter)
	})
}

type handler struct{ deps Dependencies }

func (h handler) listBooks(w http.ResponseWriter, request *http.Request) {
	books, err := h.deps.ListBooks.Execute(request.Context())
	if err != nil {
		slog.ErrorContext(request.Context(), "failed to list books", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	result := make([]bookResponse, 0, len(books))
	for _, book := range books {
		result = append(result, newBookResponse(book))
	}
	response.WriteJSON(w, request, http.StatusOK, result)
}

func (h handler) getChapter(w http.ResponseWriter, request *http.Request) {
	bookID := request.PathValue("bookId")
	number, err := strconv.Atoi(request.PathValue("number"))
	if err != nil || number <= 0 {
		http.Error(w, "invalid chapter number", http.StatusBadRequest)
		return
	}

	chapter, err := h.deps.GetChapter.Execute(request.Context(), bookID, number)
	if err != nil {
		if errors.Is(err, domain.ErrChapterNotFound) {
			http.Error(w, "chapter not found", http.StatusNotFound)
			return
		}
		slog.ErrorContext(request.Context(), "failed to get chapter", "book_id", bookID, "chapter_number", number, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.WriteJSON(w, request, http.StatusOK, newChapterResponse(chapter))
}
