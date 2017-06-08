package mux

import (
	"net/http"
)

// DefaultNotFoundHandler - The default handler for NotFound errors
func DefaultNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
