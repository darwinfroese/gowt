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

		mux.RegisterRoute(test.route, nil)
		l := len(mux.routes)

		if l != test.expectedCount {
			t.Errorf("FAIL - Expceted %d routes but have %d", test.expectedCount, l)
		}
	}
}

var errorHandlerTests = []struct {
	description, route           string
	handler                      http.HandlerFunc
	expectedResponse             response
	expectedRegistrationResponse bool
}{{
	description:                  "Testing: When not registering a not found handler the default not found handler's response is returned.",
	route:                        "/notfound",
	handler:                      nil,
	expectedResponse:             response{Body: http.StatusText(http.StatusNotFound), Code: http.StatusNotFound},
	expectedRegistrationResponse: false,
}, {
	description: "Testing: When registering a new not found handler the new not found handler's response is returned and we are informed of the overwrite.",
	route:       "/notfound",
	handler: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Couldn't find the handler")
	},
	expectedResponse:             response{Body: "Couldn't find the handler", Code: http.StatusNotFound},
	expectedRegistrationResponse: true,
}}

func TestErrorHandlerRegistration(t *testing.T) {
	t.Log("Testing error handler registration...")

	m := NewMux()

	for i, test := range errorHandlerTests {
		t.Logf("[ %02d ] %s", i+1, test.description)

		r := httptest.NewRequest("GET", test.route, nil)
		w := httptest.NewRecorder()

		// if the test doesn't register a handler we want it to pass
		overwritten := test.expectedRegistrationResponse

		if test.handler != nil {
			overwritten = m.RegisterErrorHandler(http.StatusNotFound, test.handler)
		}

		if overwritten != test.expectedRegistrationResponse {
			t.Logf("[FAIL] :: Expected the overwritten response to be %v but was %v.\n", test.expectedRegistrationResponse, overwritten)
		}

		m.ServeHTTP(w, r)

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
