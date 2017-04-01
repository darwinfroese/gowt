package mux

import (
	"net/http"
)

// Mux is a multiplexer
type Mux struct {
	Routes []string
}

// NewMux returns a new Mux object
func NewMux() *Mux {
	return &Mux{}
}

// AddRoute adds a route to the mux
func (m *Mux) AddRoute(route string) bool {
	i, ok := m.containsRoute(route)
	if ok {
		m.Routes[i] = route
		return true
	}

	m.Routes = append(m.Routes, route)

	return true
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
