package mux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type response struct {
	Body string
	Code int
}

var routeRegistrationTests = []struct {
	description, route string
	expectedCount      int
}{{
	description:   "Testing: Route registration on new mux should add one route",
	route:         "testroute",
	expectedCount: 1,
}, {
	description:   "Testing: Registering a second route should increment the route count by one",
	route:         "testroute2",
	expectedCount: 2,
}, {
	description:   "Testing: Registering a registered route should overwrite existing route",
	route:         "testroute2",
	expectedCount: 2,
}}

func TestRouteRegistration(t *testing.T) {
	t.Log("Testing route registration...")

	mux := NewMux()

	for i, test := range routeRegistrationTests {
		t.Logf("[ %02d ] %s", i+1, test.description)

		mux.AddRoute(test.route, nil)
		l := len(mux.Routes)

		if l != test.expectedCount {
			t.Errorf("FAIL - Expceted %d routes but have %d", test.expectedCount, l)
		}
	}
}

var errorHandlerTests = []struct {
	description, route string
	mux                Mux
	expectedResponse   response
}{{
	description:      "Testing: Default not found handler returns expected response.",
	route:            "/notfound",
	mux:              Mux{NotFoundHandler: DefaultNotFoundHandler},
	expectedResponse: response{Body: http.StatusText(http.StatusNotFound), Code: http.StatusNotFound},
}, {
	description:      "Testing: Default internal server error handler returns expected response.",
	route:            "/{}/servererror",
	mux:              Mux{InternalErrorHandler: DefaultInternalServerHandler},
	expectedResponse: response{Body: http.StatusText(http.StatusInternalServerError), Code: http.StatusInternalServerError},
}, {
	description: "Testing: Overwritten not found handler returns expected response.",
	route:       "/notfound",
	mux: Mux{NotFoundHandler: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Couldn't find the handler")
	}},
	expectedResponse: response{Body: "Couldn't find the handler", Code: http.StatusNotFound},
}, {
	description: "Testing: Overwritten internal server error handler returns expected response.",
	route:       "/{}/servererror",
	mux: Mux{InternalErrorHandler: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "There was an error")
	}},
	expectedResponse: response{Body: "There was an error", Code: http.StatusInternalServerError},
}}

func TestErrorHandlerRegistration(t *testing.T) {
	t.Log("Testing error handler registration...")

	for i, test := range errorHandlerTests {
		t.Logf("[ %02d ] %s", i+1, test.description)

		r := httptest.NewRequest("GET", test.route, nil)
		w := httptest.NewRecorder()

		test.mux.ServeHTTP(w, r)

		if w.Code != test.expectedResponse.Code {
			t.Logf("[FAIL] :: Expected status code %d but got status code %d.\n", test.expectedResponse.Code, w.Code)
			t.Fail()
		}

		body := strings.TrimSpace(w.Body.String())

		if body != test.expectedResponse.Body {
			t.Logf("[FAIL] :: Expected body \"%s\" but got body \"%s\"", test.expectedResponse.Body, body)
			t.Fail()
		}
	}
}
