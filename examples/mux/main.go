package main

import (
	"net/http"

	"github.com/darwinfroese/gowt/mux"
)

func main() {
	m := mux.NewMux()

	http.ListenAndServe(":8080", m)
}
