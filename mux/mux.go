package mux

import (
	"errors"
	"fmt"
	"net/http"
)

// Mux - A multiplexer object that is used for registering routes
//
// routes []Route - The array of routes that have been registered to the multiplexer
// errorHandlers map[int]Route - A map of routes to HTTP status codes
// rootNode - The root node that maps to the "/" node as the root of
// the tree of routes
// logger - A logger interface that can be set by a consumer so that
// the mux can log actions to the users logging system
type Mux struct {
	routes        []Route
	errorHandlers map[int]http.HandlerFunc

	logger
}

// NewMux returns a new Mux object with the default not found handler registered,
// this returns 404 a handler wasn't found for the route received.
func NewMux() *Mux {
	errorHandlers := make(map[int]http.HandlerFunc, 1)
	errorHandlers[http.StatusNotFound] = DefaultNotFoundHandler

	return &Mux{
		errorHandlers: errorHandlers,
	}
}

// RegisterLogger registers a logger for the multiplexer that can log actions
// to the provided logging system. logger is a simple interface that provides
// a (hopefully) commonplace log functionality.
//
// RegisterLogger will make an attempt to write an info level log entry to
// verify that the logger is working.
func (m *Mux) RegisterLogger(l logger) {
	m.logger = l

	m.log(infoLevel, "Logger has been registered for GOWT Mux.")
}

// RegisterHandler adds a Handler to the multiplexer for the route specified. If the
// route that is being used has already been added, the existing route will be
// replaced.
// If the route was of type HandlerFunc, the HandlerFunc will be replaced with a
// Handler.
func (m *Mux) RegisterHandler(route string, handler http.Handler) (*Route, error) {
	gh := gowtHandler{handler: handler}

	return m.register(route, gh)
}

// RegisterRoute adds a HandlerFunc to the multiplexer for the route specified. If
// the route that is being used has already been added, the existing route will be
// replaced.
// If the route was of type Handler, the Handler will be replaced with a HandlerFunc.
func (m *Mux) RegisterRoute(route string, handler http.HandlerFunc) (*Route, error) {
	gh := gowtHandler{handlerFunc: handler}

	return m.register(route, gh)
}

// RegisterErrorHandler registers an http.HandlerFunc for a status code providing
// a central place for error handlers to live.
//
// The function returns true if an existing error handler was updated/overwritten
func (m *Mux) RegisterErrorHandler(statusCode int, handler http.HandlerFunc) bool {

	// if ok is true, the map contained a value
	_, ok := m.errorHandlers[statusCode]
	m.errorHandlers[statusCode] = handler

	return ok
}

// GetVariables returns a slice of interface{} that contains all the variables for
// request.
func (m *Mux) GetVariables(request *http.Request) (variables []interface{}, err error) {
	var infoList []variableInfo
	for _, route := range m.routes {
		if matchRoute(route, request.URL.Path) {
			infoList = append(infoList, route.variables...)
		}
	}

	if len(infoList) == 0 {
		err = errors.New("No variables matched for the route and request")
		return
	}

	for _, v := range infoList {
		val, e := getVariableFromRequest(v, request.URL.Path)

		if e != nil {
			variables = nil
			err = e
			return
		}

		variables = append(variables, val)
	}

	return
}

// GetVariableByName returns an interface{} that contains the value for the request
func (m *Mux) GetVariableByName(name string, request *http.Request) (variable interface{}, err error) {
	var infoList []variableInfo
	for _, route := range m.routes {
		if matchRoute(route, request.URL.Path) {
			infoList = append(infoList, route.variables...)
		}
	}

	if len(infoList) == 0 {
		err = fmt.Errorf("No variables found for url \"%s\"", request.URL.Path)
		return
	}

	for _, v := range infoList {
		if v.name == name {
			variable, err = getVariableFromRequest(v, request.URL.Path)
		}
	}

	return
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

	h := m.errorHandlers[http.StatusNotFound]
	h(w, r)
}
