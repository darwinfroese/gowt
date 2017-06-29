package mux

import (
	"errors"
	"net/http"
)

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

// RegisterHandler adds a Handler to the multiplexer for the route specified. If the
// route that is being used has already been added, the existing route will be
// replaced.
// If the route was of type HandlerFunc, the HandlerFunc will be replaced with a
// Handler.
func (m *Mux) RegisterHandler(route string, handler http.Handler) (*Route, error) {
	gh := gowtHandler{handler: handler}

	return register(m, route, gh)
}

// RegisterRoute adds a HandlerFunc to the multiplexer for the route specified. If
// the route that is being used has already been added, the existing route will be
// replaced.
// If the route was of type Handler, the Handler will be replaced with a HandlerFunc.
func (m *Mux) RegisterRoute(route string, handler http.HandlerFunc) (*Route, error) {
	gh := gowtHandler{handlerFunc: handler}

	return register(m, route, gh)
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

// GetVariables returns a slice of interface{} that contains all the variables for
// request.
func (m *Mux) GetVariables(request *http.Request) ([]interface{}, error) {
	var infoList []variableInfo
	for _, route := range m.routes {
		if matchRoute(route, request.URL.Path) {
			infoList = append(infoList, route.variables...)
		}
	}

	if len(infoList) == 0 {
		return nil, errors.New("No variables matched for the route and request")
	}

	var vars []interface{}
	for _, v := range infoList {
		i := getVariableFromRequest(v, request.URL.Path)
		vars = append(vars, i)
	}

	return vars, nil
}

// GetVariableByName returns an interface{} that contains the value for the request
func (m *Mux) GetVariableByName(name, request string) (interface{}, error) {
	return nil, nil
}

// ServeHTTP matches the route incoming to the routes registered and calls the
// matched handler. If the route contains a variable, the match is based around
// the variable value
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	for _, route := range m.routes {
		if matchRoute(route, url) {
			call(route, w, r)
			return
		}
	}

	h := m.ErrorHandlers[http.StatusNotFound]
	h(w, r)
}
