// Package routing defines the HTTP routing port used by inbound adapters.
package routing

import "net/http"

// Router is the small routing contract required by HTTP endpoint registrars.
// Framework-specific implementations belong to infrastructure adapters.
type Router interface {
	Get(pattern string, handler http.HandlerFunc)
	Post(pattern string, handler http.HandlerFunc)
	Use(middlewares ...func(http.Handler) http.Handler)
	Route(pattern string, fn func(Router))
}
