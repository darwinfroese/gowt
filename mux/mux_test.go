package mux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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
	description:           "Testing: registering a route with one variable should have a count of one variable stored.",
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
}, {
	description:  "Testing: Registering a route with two variables should have a count of two variables stored.",
	route:        "/testing/{name: string}/variable/{age: int}",
	requestRoute: "/testing/darwin/variable/1234",
	handler: func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Two variable function")
	},
	expectedResponse:      response{Body: "Two variable function", Code: 200},
	expectedVariableCount: 2,
}, {
	description:  "Testing: Registering a route with multiple variables should return the correct number of variables stored.",
	route:        "/testing/{name}/{age}/{otherID}",
	requestRoute: "/testing/d/1/d",
	handler: func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Multi variable function")
	},
	expectedResponse:      response{Body: "Multi variable function", Code: 200},
	expectedVariableCount: 3,
}}

func TestVariableRouteRegistration(t *testing.T) {
	t.Log("Testing registering routes with variables in route.")

	m := NewMux()

	for i, test := range variableRouteRegistrationTests {
		t.Logf("[ %02d ] %s", i+1, test.description)

		route, err := m.RegisterRoute(test.route, test.handler)

		if err != nil {
			t.Logf("[FAIL] :: Failed to register the route. Error: \"%s\".", err.Error())
			t.FailNow()
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

var variableListRetrievalTests = []struct {
	description, routeURL, requestURL string
	expectedValues                    []interface{}
	expectedErrorMessage              string
}{{
	description:          "Testing: When registering no variables in the route no variables are returned.",
	routeURL:             "/test/no/variables",
	requestURL:           "/test/no/variables",
	expectedValues:       nil,
	expectedErrorMessage: "No variables matched for the route and request",
}, {
	description:          "Testing: When registering a route with a variable in the route the expected variable is returned.",
	routeURL:             "/test/{name: string}/test",
	requestURL:           "/test/darwin/test",
	expectedValues:       []interface{}{"darwin"},
	expectedErrorMessage: "",
}, {
	description:          "Testing: When registering a route with a variable at the end of the route the expceted variable is returned.",
	routeURL:             "/test/{name:string}",
	requestURL:           "/test/darwin",
	expectedValues:       []interface{}{"darwin"},
	expectedErrorMessage: "",
}, {
	description:          "Testing: When registering a route with two variables in the route the expected variables are returned.",
	routeURL:             "/test/{name: string}/{age: int}",
	requestURL:           "/test/darwin/1234",
	expectedValues:       []interface{}{"darwin", 1234},
	expectedErrorMessage: "",
}, {
	description:          "Testing: When registering any number of variables in a route the expected variables are returned.",
	routeURL:             "/test/{name}/{age:int}/{profile}/{count}",
	requestURL:           "/test/d/1/e/2",
	expectedValues:       []interface{}{"d", 1, "e", "2"},
	expectedErrorMessage: "",
}}

func TestVariableListRetrieval(t *testing.T) {
	t.Log("Testing getting all variables for a route.")

	m := NewMux()
	var values []interface{}
	var retrieveError error
	dummy := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Function was called")
		values, retrieveError = m.GetVariables(r)
	}

	for i, test := range variableListRetrievalTests {
		t.Logf("[ %02d ] %s", i+1, test.description)

		_, err := m.RegisterRoute(test.routeURL, dummy)

		if err != nil {
			t.Logf("[FAIL] :: Failed to register the route. Error: \"%s\".", err.Error())
			t.FailNow()
		}

		r := httptest.NewRequest("GET", test.requestURL, nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, r)

		if strings.TrimSpace(w.Body.String()) != "Function was called" {
			t.Logf("[FAIL] :: The dummy function was never called.")
			t.FailNow()
		}

		if retrieveError != nil {
			if retrieveError.Error() != test.expectedErrorMessage {
				t.Logf("[FAIL] :: Got an unexpected error retrieving the variables. Expected \"%s\" but got \"%s\".", test.expectedErrorMessage, retrieveError.Error())
				t.Fail()
			}
		}

		if len(test.expectedValues) != len(values) {
			t.Logf("[FAIL] :: Expected %d variables but got %d variables.", len(test.expectedValues), len(values))
			t.Fail()
		}

		if !reflect.DeepEqual(test.expectedValues, values) {
			t.Logf("[FAIL] :: Expected to get %+v but got %+v.", test.expectedValues, values)
			t.Fail()
		}
	}
}

var variableByNameRetrievalTests = []struct {
	description, routeURL, requestURL string
	expectedValue                     interface{}
	expectedErrorMessage              string
}{{
	description:          "Testing: When no variable is registered in the route no variable information is returned.",
	routeURL:             "/test/no/variables",
	requestURL:           "/test/no/variables",
	expectedValue:        nil,
	expectedErrorMessage: "No variables found for url \"/test/no/variables\"",
}, {
	description:          "Testing: When a variable is registered in a route but the wrong variable name is used no information is returned.",
	routeURL:             "/test/{profile: string}/bad",
	requestURL:           "/test/darwin/bad",
	expectedValue:        nil,
	expectedErrorMessage: "No variable was found that matched for \"name\"",
}, {
	description:          "Testing: When a variable is registered in a route and the right variable name is used the correct information is returned.",
	routeURL:             "/test/{name: string}/good",
	requestURL:           "/test/darwin/good",
	expectedValue:        "darwin",
	expectedErrorMessage: "",
}, {
	description:          "Testing: When a variable is registered at the end of the route and the right variable name is used the correct information is returned.",
	routeURL:             "/test/{name:string}",
	requestURL:           "/test/darwin",
	expectedValue:        "darwin",
	expectedErrorMessage: "",
}, {
	description:          "Testing: When two variables are registered in a route and the right variable name is used the correct value is returned.",
	routeURL:             "/test/{age: int}/{name: string}",
	requestURL:           "/test/1234/darwin",
	expectedValue:        "darwin",
	expectedErrorMessage: "",
}, {
	description:          "Testing: when multiple variables are registered in a route and the right variable name is used the correct value is returned.",
	routeURL:             "/test/{age}/{count}/{name}",
	requestURL:           "/test/1/2/darwin",
	expectedValue:        "darwin",
	expectedErrorMessage: "",
}}

func TestVariableRetrievalByName(t *testing.T) {
	t.Log("Testing getting a variable by name.")

	m := NewMux()
	var value interface{}
	var retrievalError error
	dummy := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Function was called")
		value, retrievalError = m.GetVariableByName("name", r)
	}

	for i, test := range variableByNameRetrievalTests {
		t.Logf("[ %02d ] %s", i+1, test.description)

		_, err := m.RegisterRoute(test.routeURL, dummy)

		if err != nil {
			t.Logf("[FAIL] :: Failed to register the route. Error: \"%s\".", err.Error())
			t.FailNow()
		}

		r := httptest.NewRequest("GET", test.requestURL, nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, r)

		if strings.TrimSpace(w.Body.String()) != "Function was called" {
			t.Logf("[FAIL] :: The dummy function was never called.")
			t.FailNow()
		}

		if retrievalError != nil {
			if retrievalError.Error() != test.expectedErrorMessage {
				t.Logf("[FAIL] :: Got an unexpected error retrieving the variable. Expected \"%s\" but got \"%s\".", test.expectedErrorMessage, retrievalError.Error())
				t.Fail()
			}
		}

		if test.expectedValue != value {
			t.Logf("[FAIL] :: Expected %+v but got %+v instead.", test.expectedValue, value)
			t.Fail()
		}
	}
}
