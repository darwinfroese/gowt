package mux

import "strconv"

// cast uses the reflect.Kind value to cast the interface to a type but
// still an interface. We need to do this so we get back an int or a string
// as the underyling type and because we can't use reflect.Kind to cast
func cast(kind string, val string) (retval interface{}, err error) {
	switch kind {
	case "int":
		retval, err = convertInt(val, 0)
	case "int8":
		retval, err = convertInt(val, 8)
	case "int16":
		retval, err = convertInt(val, 16)
	case "int32":
		retval, err = convertInt(val, 32)
	case "int64":
		retval, err = convertInt(val, 64)
	case "uint":
		retval, err = convertUint(val, 0)
	case "uint8":
		retval, err = convertUint(val, 8)
	case "uint16":
		retval, err = convertUint(val, 16)
	case "uint32":
		retval, err = convertUint(val, 32)
	case "uint64":
		retval, err = convertUint(val, 64)
	case "string":
		retval = val
		err = nil
	default:
		retval = val
		err = nil
	}

	return
}

// convertInt sets retval to 'nil' on error since ParseInt returns 0 on error
// and we are returning an interface not an integer.
func convertInt(value string, base int) (retval interface{}, err error) {
	val, err := strconv.ParseInt(value, 10, base)
	if err != nil {
		retval = nil
		return
	}

	switch base {
	case 0:
		retval = int(val)
	case 8:
		retval = int8(val)
	case 16:
		retval = int16(val)
	case 32:
		retval = int32(val)
	case 64:
		retval = int64(val)
	}

	return
}

// convertUint sets retval to 'nil' on error since ParseUint returns 0 on error
// and we are returning an interface not an integer.
func convertUint(value string, base int) (retval interface{}, err error) {
	val, err := strconv.ParseUint(value, 10, base)
	if err != nil {
		retval = nil
		return
	}

	switch base {
	case 0:
		retval = uint(val)
	case 8:
		retval = uint8(val)
	case 16:
		retval = uint16(val)
	case 32:
		retval = uint32(val)
	case 64:
		retval = uint64(val)
	}

	return
}
