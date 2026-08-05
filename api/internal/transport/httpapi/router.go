package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(deps Dependencies) http.Handler {
	router := chi.NewRouter()

	registerHealthRoutes(router, deps.Health)
	registerCatalogRoutes(router, deps.Catalog)

	return router
}
