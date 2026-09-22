package httpapi_test

import (
	"net/http"

	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi"
)

func router(br booksRepository, cr chaptersRepository) http.Handler {
	return httpapi.NewHandler(httpapi.Dependencies{
		Health: health.CheckHealth{},
		Catalog: httpapi.CatalogDependencies{
			ListBooks:  queries.NewListBooks(br),
			GetChapter: queries.NewGetChapter(cr),
		},
	})
}
