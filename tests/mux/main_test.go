package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/darwinfroese/gowt/mux"
)

const expectedOutput string = "Hello World"

var srv *httptest.Server

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	shutdown()

	os.Exit(code)
}

func setup() {
	m := mux.NewMux()

	m.RegisterRoute("/hello", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, expectedOutput) })

	srv = httptest.NewServer(m)
}

func shutdown() {
	srv.Close()
}

func TestRoute(t *testing.T) {
	t.Log("Testing Routes")

	res, err := http.Get(srv.URL + "/hello")
	if err != nil {
		t.Errorf(" FAILED - %s", err.Error())
	}

	msg, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf(" FAILED - %s", err.Error())
	}

	if string(msg) != expectedOutput {
		t.Errorf(" FAILED - Expected \"%s\" but got \"%s\"", expectedOutput, msg)
	}
}
