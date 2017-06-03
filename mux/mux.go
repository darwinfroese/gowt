package mux

import (
	"net/http"
	"strings"
)

// Mux - A multiplexer object that is used for registering routes
//
// Routes []Route - The array of routes that have been registered to the multiplexer
// ErrorHandlers map[int]Route - A map of routes to HTTP status codes
type Mux struct {
	Routes               []Route
	NotFoundHandler      http.HandlerFunc
	InternalErrorHandler http.HandlerFunc
}

// NewMux returns a new Mux object
func NewMux() *Mux {
	return &Mux{
		NotFoundHandler:      DefaultNotFoundHandler,
		InternalErrorHandler: DefaultInternalServerHandler,
	}
}

// AddRoute adds a route to the multiplexer
func (m *Mux) AddRoute(route string, handler http.HandlerFunc) *Route {
	i, ok := m.containsRoute(route)
	if ok {
		m.Routes[i].Handler = handler
		return &m.Routes[i]
	}

	r := Route{URL: route, Handler: handler}
	m.Routes = append(m.Routes, r)

	return &r
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	// Temporary - implement proper variable handling
	if strings.Contains(url, "{") {
		m.InternalErrorHandler(w, r)
		return
	}

	for _, route := range m.Routes {
		if route.URL == url {
			route.Handler(w, r)
			return
		}
	}

	m.NotFoundHandler(w, r)
}
