package mux

import (
	"net/http"
)

// Mux - A multiplexer object
type Mux struct {
	Routes []Route
}

// Route - A Route Object
type Route struct {
	URL       string
	Handler   http.HandlerFunc
	SubRoutes []string
}

// NewMux returns a new Mux object
func NewMux() *Mux {
	return &Mux{}
}

// AddRoute adds a route to the mux
func (m *Mux) AddRoute(route string, handler http.HandlerFunc) bool {
	i, ok := m.containsRoute(route)
	if ok {
		m.Routes[i].Handler = handler
		return true
	}

	m.Routes = append(m.Routes, Route{URL: route, Handler: handler})

	return true
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
