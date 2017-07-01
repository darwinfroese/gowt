package mux

import "reflect"
import "strconv"

// getKind - Returns the reflect.Kind for a type defined by a string
func getKind(kind string) reflect.Kind {
	switch kind {
	case "string":
	default:
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
	}

	return reflect.String
}

// cast uses the reflect.Kind value to cast the interface to a type but
// still an interface. We need to do this so we get back an int or a string
// as the underyling type and because we can't use reflect.Kind to cast
func cast(kind reflect.Kind, val string) interface{} {
	ival, _ := strconv.Atoi(val)
	switch kind {
	case reflect.Int:
		return ival
	case reflect.Int8:
		return (int8)(ival)
	case reflect.Int16:
		return (int16)(ival)
	case reflect.Int32:
		return (int32)(ival)
	case reflect.Int64:
		return (int64)(ival)
	case reflect.Uint:
		return (uint)(ival)
	case reflect.Uint8:
		return (uint8)(ival)
	case reflect.Uint16:
		return (uint16)(ival)
	case reflect.Uint32:
		return (uint32)(ival)
	case reflect.Uint64:
		return (uint64)(ival)
	case reflect.String:
	default:
		return val
	}

	return val
}
