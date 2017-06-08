package mux

import "net/http"

// Mux - A multiplexer object that is used for registering routes
//
// Routes []Route - The array of routes that have been registered to the multiplexer
// ErrorHandlers map[int]Route - A map of routes to HTTP status codes
type Mux struct {
	routes          []Route
	NotFoundHandler http.HandlerFunc
}

// NewMux returns a new Mux object
func NewMux() *Mux {
	return &Mux{
		NotFoundHandler: DefaultNotFoundHandler,
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

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	for _, route := range m.routes {
		if route.URL == url {
			route.Handler(w, r)
			return
		}
	}

	m.NotFoundHandler(w, r)
}
