package mux

import (
	"testing"
)

var routeRegistrationTests = []struct {
	description, route string
	expectedOutcome    bool
	expectedCount      int
}{{
	description:     "Testing: Route registration on new mux should succeed",
	route:           "testroute",
	expectedOutcome: true,
	expectedCount:   1,
}, {
	description:     "Testing: Registering a second route should succeed",
	route:           "testroute2",
	expectedOutcome: true,
	expectedCount:   2,
}, {
	description:     "Testing: Registering a registered route should overwrite existing route",
	route:           "testroute2",
	expectedOutcome: true,
	expectedCount:   2,
}}

func TestRouteRegistration(t *testing.T) {
	t.Log("Testing route registration...")

	mux := NewMux()

	for i, test := range routeRegistrationTests {
		t.Logf("%02d %s", i, test.description)

		ok := mux.AddRoute(test.route, nil)
		l := len(mux.Routes)

		if ok != test.expectedOutcome {
			t.Errorf("FAIL - Expected %v but got %v", test.expectedOutcome, ok)
		}
		if l != test.expectedCount {
			t.Errorf("FAIL - Expceted %d routes but have %d", test.expectedCount, l)
		}
	}
}

var errorRegistrationTests = []struct {
	description     string
	errorCode       int
	expectedOutcome bool
	expectedCount   int
}{{
	description:     "Testing: Registering an error handler on a new mux should succeed.",
	errorCode:       404,
	expectedOutcome: true,
	expectedCount:   1,
}, {
	description:     "Testing: Registering a second error handler on a new mux should succeed.",
	errorCode:       500,
	expectedOutcome: true,
	expectedCount:   2,
}, {
	description:     "Testing: Registering an error handler for a registered error code should replace the handler.",
	errorCode:       404,
	expectedOutcome: true,
	expectedCount:   2,
}}

func TestErrorHandlerRegistration(t *testing.T) {
	t.Log("Testing error handler registration...")
	mux := NewMux()

	for i, test := range errorRegistrationTests {
		t.Logf("%02d %s", i, test.description)

		ok := mux.AddErrorHandler(test.errorCode, nil)
		l := len(mux.ErrorHandlers)

		if ok != test.expectedOutcome {
			t.Errorf("FAIL - Expected %v but got %v", test.expectedOutcome, ok)
		}

		if l != test.expectedCount {
			t.Errorf("FAIL - Expected %d routes but have %d routes", test.expectedCount, l)
		}
	}
}
