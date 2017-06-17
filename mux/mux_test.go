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

var variableRouteRegistrationTests = []struct {
	description, route, requestRoute string
	handler                          http.HandlerFunc
	expectedResponse                 response
	expectedVariableCount            int
}{{
	description:           "Testing: registering a route with no variables should have no variables stored.",
	route:                 "/testingnovariables",
	requestRoute:          "/testingnovariables",
	handler:               func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello World") },
	expectedResponse:      response{Body: "Hello World", Code: 200},
	expectedVariableCount: 0,
}, {
	description:           "Testing: registering a route with one variable should have one variable stored.",
	route:                 "/testing/{name: string}/variable",
	requestRoute:          "/testing/darwin/variable",
	handler:               func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello World") },
	expectedResponse:      response{Body: "Hello World", Code: 200},
	expectedVariableCount: 1,
}, {
	description:           "Testing: updating a route with a variable should keep the count at one.",
	route:                 "/testing/{name: string}/variable",
	requestRoute:          "/testing/darwin/variable",
	handler:               func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "New function") },
	expectedResponse:      response{Body: "New function", Code: 200},
	expectedVariableCount: 1,
}}

func TestVariableRouteRegistration(t *testing.T) {
	t.Log("Testing registering routes with variables in route.")

	m := NewMux()

	for i, test := range variableRouteRegistrationTests {
		t.Logf("[ %02d ] %s", i+1, test.description)

		route, err := m.RegisterRoute(test.route, test.handler)

		if err != nil {
			t.Logf("[FAIL] :: Failed to register the route. Error: \"%s\".", err.Error())
			t.Fail()
		}

		r := httptest.NewRequest("GET", test.requestRoute, nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, r)

		if test.expectedResponse.Code != w.Code {
			t.Logf("[FAIL] :: Expected status code %d but got status code %d.", test.expectedResponse.Code, w.Code)
			t.Fail()
		}

		body := strings.TrimSpace(w.Body.String())
		if test.expectedResponse.Body != body {
			t.Logf("[FAIL] :: Expected body \"%s\" but got body \"%s\".", test.expectedResponse.Body, body)
			t.Fail()
		}

		if test.expectedVariableCount != len(route.variables) {
			t.Logf("[FAIL] :: Expected %d variables but got %d variables instead.", test.expectedVariableCount, len(route.variables))
			t.Fail()
		}
	}
}
