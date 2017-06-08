package mux

import "net/http"

// Mux - A multiplexer object that is used for registering routes
//
// Routes []Route - The array of routes that have been registered to the multiplexer
// ErrorHandlers map[int]Route - A map of routes to HTTP status codes
type Mux struct {
	routes        []Route
	ErrorHandlers map[int]http.HandlerFunc
}

// NewMux returns a new Mux object with the default not found handler registered,
// this returns 404 a handler wasn't found for the route received
func NewMux() *Mux {
	errorHandlers := make(map[int]http.HandlerFunc, 1)
	errorHandlers[http.StatusNotFound] = DefaultNotFoundHandler

	return &Mux{
		ErrorHandlers: errorHandlers,
	}
}

// RegisterRoute adds a route to the multiplexer
func (m *Mux) RegisterRoute(route string, handler http.HandlerFunc) *Route {
	i, ok := m.containsRoute(route)
	if ok {
		m.routes[i].Handler = handler
		return &m.routes[i]
	}

	r := Route{URL: route, Handler: handler}
	m.routes = append(m.routes, r)

	return &r
}

// RegisterErrorHandler registers an http.HandlerFunc for a status code providing
// a central place for error handlers to live.
//
// The function returns true if an existing error handler was updated/overwritten
func (m *Mux) RegisterErrorHandler(statusCode int, handler http.HandlerFunc) bool {

	// if ok is true, the map contained a value
	_, ok := m.ErrorHandlers[statusCode]
	m.ErrorHandlers[statusCode] = handler

	return ok
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	for _, route := range m.routes {
		if route.URL == url {
			route.Handler(w, r)
			return
		}
	}

	h := m.ErrorHandlers[http.StatusNotFound]
	h(w, r)
}
