package mux

import (
	"reflect"
	"testing"
)

var routeExtractionTests = []struct {
	description, route string
	expected           []variableInfo
	errorMessage       string
}{{
	description: "Testing: When providing a route with now variables, no variable should be extracted.",
	route:       "/test/testing",
	expected:    nil,
}, {
	description: "Testing: When providing a good route, the variable will be extracted.",
	route:       "/test/{name: string}/test",
	expected:    []variableInfo{variableInfo{name: "name", kind: reflect.String}},
}, {
	description: "Testing: When providing a good route, any variable type will be extracted.",
	route:       "/test/{age: int}/test",
	expected:    []variableInfo{variableInfo{name: "age", kind: reflect.Int}},
}, {
	description: "Testing: When providing a route without a type, a default type of interface{} will be used.",
	route:       "/test/{name}/test",
	expected:    []variableInfo{variableInfo{name: "name", kind: reflect.Interface}},
}, {
	description: "Testing: When providing multiple variables, all variables are returned",
	route:       "/test/{name: string}/{age: int}",
	expected:    []variableInfo{variableInfo{name: "name", kind: reflect.String}, variableInfo{name: "age", kind: reflect.Int}},
}, {
	description:  "Testing: When providing a malformed route, the variable will not be extracted.",
	route:        "/test/{name: string/test",
	expected:     nil,
	errorMessage: "Missing '}' in route variable declaration",
}, {
	description:  "Testing: When providing a route without a name the variable will not be extracted.",
	route:        "/test/{:string}/test",
	expected:     nil,
	errorMessage: "Missing the variable name in variable declaration",
}, {
	description:  "Testing: When providing an empty variable decleration the variable should not be extracted.",
	route:        "/test/{}/test",
	expected:     nil,
	errorMessage: "Missing variable information in variable declaration",
}}

func TestRouteExtraction(t *testing.T) {
	t.Log("Testing route extraction...")

	for i, test := range routeExtractionTests {
		t.Logf("[ %02d ] %s", i, test.description)

		results, err := getVariablesFromRoute(test.route)

		if err != nil && err.Error() != test.errorMessage {
			t.Logf("[FAIL] An error occured getting the variables: %s", err.Error())
			t.FailNow()
		}

		if len(results) != len(test.expected) {
			t.Logf("[FAIL] Expected %d results but got %d results.", len(test.expected), len(results))
			t.FailNow()
		}

		for i, result := range results {
			exp := test.expected[i]

			if result.name != exp.name {
				t.Logf("[FAIL] Expected a name of \"%s\" but got a name of \"%s\".", exp.name, result.name)
				t.Fail()
			}

			if result.kind != exp.kind {
				t.Logf("[FAIL] Expected a kind of %v but got a kind of %v.", exp.kind, result.kind)
				t.Fail()
			}
		}

	}
}

var routeMatchingTests = []struct {
	description, requestURL string
	route                   Route
	expectedMatch           bool
}{{
	description:   "Testing: Matching routes without variables should match.",
	requestURL:    "/test/route",
	route:         Route{URL: "/test/route", hasVariables: false},
	expectedMatch: true,
}, {
	description:   "Testing: Non-matching routes without variables shouldn't match.",
	requestURL:    "/test/route/one",
	route:         Route{URL: "/test/route/two", hasVariables: false},
	expectedMatch: false,
}, {
	description:   "Testing: Matching routes with variables should match.",
	requestURL:    "/profile/darwin/name",
	route:         Route{URL: "/profile/{name: string}/name", hasVariables: true},
	expectedMatch: true,
}, {
	description:   "Testing: Non-matching routes with variables shouldn't match.",
	requestURL:    "/profile/darwin/account",
	route:         Route{URL: "/profile/{name: string}/name", hasVariables: true},
	expectedMatch: false,
}}

func TestRouteMatching(t *testing.T) {
	t.Log("Testing route matching function.")

	for i, test := range routeMatchingTests {
		t.Logf("[ %02d ] %s", i, test.description)

		match := matchRoute(test.route, test.requestURL)

		if match != test.expectedMatch {
			t.Logf("[FAIL] :: Expected %v but got %v instead.", test.expectedMatch, match)
			t.Fail()
		}
	}
}
