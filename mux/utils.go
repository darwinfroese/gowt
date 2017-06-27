package mux

import (
	"errors"
	"fmt"
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
		if r.URL == route {
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
		return route.URL == requestURL
	}

	// match the url to the left of the variable declaration first
	leftIdx := strings.Index(route.URL, "{")
	if route.URL[:leftIdx-1] != requestURL[:leftIdx-1] {
		return false
	}

	// match the url to the right of the variable declaration now
	rightIdxRoute := strings.Index(route.URL[leftIdx:], "/")
	rightIdxRequest := strings.Index(requestURL[leftIdx:], "/")
	if route.URL[leftIdx+rightIdxRoute:] != requestURL[leftIdx+rightIdxRequest:] {
		return false
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

func getVariableFromRequest(info variableInfo, request string) interface{} {
	name := info.name

	startIdx := strings.Index(info.route, "{"+name)
	endIdx := startIdx + strings.Index(request[startIdx:], "/")

	fmt.Println("startIdx: ", startIdx, "endIdx: ", endIdx)

	return request[startIdx:endIdx]
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
