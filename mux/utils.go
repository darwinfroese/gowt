package mux

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
)

// variableInfo contains the information about the variable
// that is extracted from the route
type variableInfo struct {
	name, route string
	kind        reflect.Kind
}

// containsRoute performs a simple check on if the route is
// already registered in the multiplexer. This is matched
// exactly so the same route can be registered if the variable
// names are different
func (m *Mux) containsRoute(route string) (int, bool) {
	for i, r := range m.routes {
		if r.url == route {
			return i, true
		}
	}

	return -1, false
}

// matchRoute attempts to match the request URL to the route
// that is the same.
//
// Exact matching is used of there are no variables in the route.
// If there are variables in the route then it matches around those
func matchRoute(route Route, requestURL string) bool {
	if !route.hasVariables {
		return route.url == requestURL
	}
	urlBlocks := cleanSlice(strings.Split(route.url, "/"))
	reqBlocks := cleanSlice(strings.Split(requestURL, "/"))

	if len(urlBlocks) != len(reqBlocks) {
		return false
	}

	for i, block := range urlBlocks {
		if block[0] == '{' && block[len(block)-1] == '}' {
			continue
		}

		if block != reqBlocks[i] {
			return false
		}
	}

	return true
}

// getVariablesFromRoute - Returns an array of variableInfo structs for
// the variables in the route
func getVariablesFromRoute(route string) ([]variableInfo, error) {
	// Check if we have a variable
	if !strings.Contains(route, "{") {
		return nil, nil
	}

	variables, err := getVariableStrings(route)

	// Bubble up error to the user
	if err != nil {
		return nil, err
	}

	infoSplice := []variableInfo{}
	for _, variable := range variables {
		info, err := getVariableInfo(variable)
		info.route = route

		if err != nil {
			return nil, err
		}

		infoSplice = append(infoSplice, info)
	}

	return infoSplice, nil
}

// getVariableFromRequest returns the value from the request
func getVariableFromRequest(info variableInfo, request string) interface{} {
	urlBlocks := cleanSlice(strings.Split(info.route, "/"))
	reqBlocks := cleanSlice(strings.Split(request, "/"))

	var val interface{}
	for i, block := range urlBlocks {
		if strings.Contains(block, info.name) {
			val = cast(info.kind, reqBlocks[i])
		}
	}

	return val
}

// getVariableStrings - Returns all the strings for the variables found
// inside of a route. For example:
//
// route: "test/{name: string}/testing"
// returns: "{name: string}"
func getVariableStrings(route string) ([]string, error) {
	err := checkVariableSyntax(route)

	// Bubble up error
	if err != nil {
		return nil, err
	}

	variables := []string{}
	rawSplice := getVariables(strings.SplitAfter(route, "}"))

	for _, s := range rawSplice {
		variables = append(variables, s[strings.Index(s, "{"):])
	}

	return variables, nil
}

// getVariableInfo extracts the information from the block of the
// route that contains the variable decleraton and will return
// an error if any information is missing.
func getVariableInfo(variable string) (variableInfo, error) {
	decon := variable[1 : len(variable)-1]

	if strings.Index(decon, ":") == 0 {
		return variableInfo{}, errors.New("Missing the variable name in variable declaration")
	}

	if decon == "" {
		return variableInfo{}, errors.New("Missing variable information in variable declaration")
	}

	pieces := strings.Split(decon, ":")
	// if kindString is empty, it'll default to interface
	kindString := ""
	if len(pieces) > 1 {
		kindString = strings.TrimSpace(pieces[1])
	}

	info := variableInfo{name: strings.TrimSpace(pieces[0]), kind: getKind(kindString)}

	return info, nil
}

// checkVariablSyntax - Checks if the number of braces matches up
// and throws an error if there's a missing brace
func checkVariableSyntax(route string) error {
	if strings.Count(route, "{") > strings.Count(route, "}") {
		return errors.New("Missing '}' in route variable declaration")
	}

	if strings.Count(route, "{") < strings.Count(route, "}") {
		return errors.New("Missing '{' in route variable declaration")
	}

	return nil
}

// getVariables - Removes all the strings from the slice that
// are empty or do not contain a variable
func getVariables(splitStrings []string) []string {
	newSlice := []string{}

	for _, s := range splitStrings {
		if s != "" && strings.Contains(s, "{") {
			newSlice = append(newSlice, s)
		}
	}

	return newSlice
}

// cleanSlice - Removes all the strings from the slice that are empty
func cleanSlice(slice []string) []string {
	newSlice := []string{}

	for _, s := range slice {
		if s != "" {
			newSlice = append(newSlice, s)
		}
	}

	return newSlice
}

// call sends the response writer and request to whichever handler
// type has been initialized
func call(route Route, w http.ResponseWriter, r *http.Request) {
	if route.handler.handler == nil {
		route.handler.handlerFunc(w, r)
		return
	}

	route.handler.handler.ServeHTTP(w, r)
}

// register does the actual registration of handlers to the multiplexer,
// this lets us have the same functionality between both of the
// registration methods while still providing two methods of registration.
func (m *Mux) register(route string, gh gowtHandler) (*Route, error) {
	i, ok := m.containsRoute(route)

	variables, err := getVariablesFromRoute(route)

	if err != nil {
		return nil, err
	}

	if ok {
		m.routes[i].handler = gh
		m.routes[i].variables = variables
		m.routes[i].hasVariables = len(variables) > 0
		return &m.routes[i], nil
	}

	r := Route{
		url:          route,
		handler:      gh,
		variables:    variables,
		hasVariables: len(variables) > 0,
	}
	m.routes = append(m.routes, r)

	return &r, nil
}
