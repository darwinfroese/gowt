package mux

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
)

// logger is an interface that should be capable of satisfying
// most common log interfaces? This lets the mux log what happens
// if a logger is provided by the consumer.
type logger interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Debug(string, ...interface{})
	Error(string, ...interface{})
}

// variableInfo contains the information about the variable
// that is extracted from the route
type variableInfo struct {
	name, route string
	kind        reflect.Kind
}

// routeNode is a struct that contains information about a "block"
// of a request or route and is used for constructing a tree of
// registered routes and walking down that tree to match routes.
type routeNode struct {
	path       string
	isVariable bool
	variableInfo

	subroutes []*routeNode
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

	varCount := len(route.variables)
	url := route.url
	req := requestURL

	for i := 0; i < varCount; i++ {
		leftIdx := strings.Index(url, "{")

		if leftIdx == -1 || leftIdx > len(req) {
			break
		}

		// match the url to the left of the variable declaration first
		if url[:leftIdx] != req[:leftIdx] {
			return false
		}

		// remove the portion of the URL that we've matched
		urlIdx := strings.Index(url[leftIdx:], "/")
		reqIdx := strings.Index(req[leftIdx:], "/")
		var rightIdxRoute int
		var rightIdxRequest int

		if urlIdx == -1 {
			rightIdxRoute = len(url[leftIdx:]) - 1
			rightIdxRequest = len(req[leftIdx:]) - 1
		} else {
			rightIdxRoute = urlIdx
			rightIdxRequest = reqIdx
		}

		url = url[leftIdx+rightIdxRoute+1:]
		req = req[leftIdx+rightIdxRequest+1:]
	}

	// if the URL and RequestURL aren't empty, lets check the end of it
	if len(url) >= 0 || len(req) >= 0 {
		return url == req
	}

	return false
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

func getVariableFromRequest(info variableInfo, request string) interface{} {
	name := info.name

	leftIdx := strings.Index(info.route, "{"+name)
	lessVar := info.route[:leftIdx]
	count := strings.Count(lessVar, "/")
	req := request

	for i := 0; i < count; i++ {
		idx := strings.Index(req, "/")
		req = req[idx+1:]
	}

	// if we still have something on the right hand side, return everything else
	i := strings.Index(req, "/")
	if i > 0 {
		return req[:i]
	}

	return req
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
	rawSplice := cleanSplice(strings.SplitAfter(route, "}"))

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

// getKind - Returns the reflect.Kind for a type defined by a string
//
// TODO: Can this be done without a switch/case?
func getKind(kind string) reflect.Kind {
	switch kind {
	case "string":
		return reflect.String
	case "int":
		return reflect.Int
	case "int8":
		return reflect.Int8
	case "int16":
		return reflect.Int16
	case "int32":
		return reflect.Int32
	case "int64":
		return reflect.Int64
	case "uint":
		return reflect.Uint
	case "uint8":
		return reflect.Uint8
	case "uint16":
		return reflect.Uint16
	case "uint32":
		return reflect.Uint32
	case "uint64":
		return reflect.Uint64
	case "interface":
	default:
		return reflect.Interface
	}

	// Default should catch this but it's a compiler error otherwise
	return reflect.Interface
}

// cleanSplice - Removes all the strings from the slice that
// are empty or do not contain a variable
func cleanSplice(splitStrings []string) []string {
	newSlice := []string{}

	for _, s := range splitStrings {
		if s != "" && strings.Contains(s, "{") {
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

// register does the actual registration of the multiplexer, this lets us have
// the same functionality between both of the registration methods while still
// providing two methods of registration.
func register(m *Mux, route string, gh gowtHandler) (*Route, error) {
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
