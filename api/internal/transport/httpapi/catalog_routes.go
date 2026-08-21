package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/domain"
	"github.com/eduardoaugustolb/versum/api/internal/ports/httprouter"
)

func registerCatalogRoutes(router httprouter.Router, deps CatalogDependencies) {
	router.Get("/books", func(w http.ResponseWriter, r *http.Request) {
		books, err := deps.ListBooks.Execute(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := make([]bookResponse, 0, len(books))
		for _, book := range books {
			response = append(response, newBookResponse(book))
		}
		writeJSON(w, http.StatusOK, response)
	})

	router.Get("/books/{bookId}/chapters/{number}", func(w http.ResponseWriter, r *http.Request) {
		bookID := r.PathValue("bookId")

		number, err := strconv.Atoi(r.PathValue("number"))
		if err != nil || number <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chapter, err := deps.GetChapter.Execute(r.Context(), bookID, number)
		if err != nil {
			if errors.Is(err, domain.ErrChapterNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, newChapterResponse(chapter))
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}
