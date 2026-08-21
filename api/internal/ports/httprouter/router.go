// Package httprouter is the routing port between domain route registration
// and whatever concrete HTTP router implements it. Its single method never
// references a third-party router package — request/response handling
// already has a driver-agnostic port in the stdlib (net/http.Handler), so
// this only needs to cover route registration itself.
package httprouter

import "net/http"

type Router interface {
	Get(pattern string, handler http.HandlerFunc)
}
