package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	Get(pattern string, handler http.HandlerFunc)
	Post(pattern string, handler http.HandlerFunc)
}

func NewRouter(deps Dependencies) http.Handler {
	router := chi.NewRouter()

	registerHealthRoutes(router, deps.Health)
	registerCatalogRoutes(router, deps.Catalog)
	registerIdentityAccessRoutes(router, deps.IdentityAccess)

	return router
}
