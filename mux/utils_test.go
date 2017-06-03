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
	description:  "Testing: When providing a route without a namew the variable will not be extracted.",
	route:        "/test/{:string}/test",
	expected:     nil,
	errorMessage: "Missing the variable name in variable declaration",
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
