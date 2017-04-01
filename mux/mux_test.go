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

		ok := mux.AddRoute(test.route)
		l := len(mux.Routes)

		if ok != test.expectedOutcome {
			t.Errorf("FAIL - Expected %v but got %v", test.expectedOutcome, ok)
		}
		if l != test.expectedCount {
			t.Errorf("FAIL - Expceted %d routes but have %d", test.expectedCount, l)
		}
	}
}
