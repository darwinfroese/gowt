package mux

import (
	"net/http"
)

// Route - A Route Object, only the object itself is exposed
// since the user needs it to call successive methods but we
// don't want to provide them with internal information.
type Route struct {
	url            string
	handler        gowtHandler
	allowedMethods []string
	hasVariables   bool
	variables      []variableInfo
}

// gowtHandler wraps around http.Handler and http.HandlerFunc
// allowing the Route object to expose one Handler that can
// internally be one or the other.
type gowtHandler struct {
	handler     http.Handler
	handlerFunc http.HandlerFunc
}
