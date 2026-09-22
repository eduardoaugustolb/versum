// Package chiadapter adapts Chi to the HTTP routing port.
package chiadapter

import (
	"net/http"

	"github.com/eduardoaugustolb/versum/api/internal/transport/httpapi/routing"
	"github.com/go-chi/chi/v5"
)

// Adapter is the Chi implementation of the inbound HTTP routing port.
type Adapter struct {
	router chi.Router
}

// New builds an adapter backed by a new Chi router.
func New() *Adapter {
	return &Adapter{router: chi.NewRouter()}
}

func (a *Adapter) Get(pattern string, handler http.HandlerFunc) {
	a.router.Get(pattern, handler)
}

func (a *Adapter) Post(pattern string, handler http.HandlerFunc) {
	a.router.Post(pattern, handler)
}

func (a *Adapter) Use(middlewares ...func(http.Handler) http.Handler) {
	a.router.Use(middlewares...)
}

func (a *Adapter) Route(pattern string, fn func(routing.Router)) {
	a.router.Route(pattern, func(child chi.Router) {
		fn(&Adapter{router: child})
	})
}

func (a *Adapter) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	a.router.ServeHTTP(w, request)
}
