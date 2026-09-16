package httpapi

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/httprouter"
)

func registerCatalogRoutes(router httprouter.Router, deps CatalogDependencies) {
	router.Get("/books", func(w http.ResponseWriter, r *http.Request) {
		books, err := deps.ListBooks.Execute(r.Context())
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"failed to list books",
				"error", err,
			)

			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
			return
		}

		response := make([]bookResponse, 0, len(books))
		for _, book := range books {
			response = append(response, newBookResponse(book))
		}

		writeJSON(w, r, http.StatusOK, response)
	})

	router.Get("/books/{bookId}/chapters/{number}", func(w http.ResponseWriter, r *http.Request) {
		bookID := r.PathValue("bookId")

		number, err := strconv.Atoi(r.PathValue("number"))
		if err != nil || number <= 0 {
			http.Error(
				w,
				"invalid chapter number",
				http.StatusBadRequest,
			)
			return
		}

		chapter, err := deps.GetChapter.Execute(r.Context(), bookID, number)
		if err != nil {
			if errors.Is(err, domain.ErrChapterNotFound) {
				http.Error(
					w,
					"chapter not found",
					http.StatusNotFound,
				)
				return
			}

			slog.ErrorContext(
				r.Context(),
				"failed to get chapter",
				"book_id", bookID,
				"chapter_number", number,
				"error", err,
			)

			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
			return
		}

		writeJSON(w, r, http.StatusOK, newChapterResponse(chapter))
	})
}
