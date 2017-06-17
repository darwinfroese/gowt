package mux

import (
	"net/http"
)

// Route - A Route Object
type Route struct {
	URL     string
	Handler http.HandlerFunc

	allowedMethods []string
	hasVariables   bool
	variables      []variableInfo
}
