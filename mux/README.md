# GOWT Mux
GOWT Mux is a simple multiplexer that satisfies the standard libraries
multiplexer and handlers allowing it to be easily dropped into most
projects that exist.

## Usage

- Routes ending in a trailing "/" are the same as routes without the trailing "/"
	- `/url/test/` is the same as `/url/test`
- `GetVariables(request)` will return an error if the variables couldn't be retrieved or if the 
variables trying to be retrieved couldn't be converted to the type specified in the route