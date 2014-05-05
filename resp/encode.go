package resp

import (
	"errors"
	"io"
	"strconv"
)

// ErrInvalidValue is returned if the value to encode is invalid.
var ErrInvalidValue = errors.New("resp: invalid value")

// Error represents an error string as defined by the RESP. It cannot
// contain \r or \n characters. It must be used as a type conversion
// so that Encode serializes the string as an Error.
type Error string

// SimpleString represents a simple string as defined by the RESP. It
// cannot contain \r or \n characters. It must be used as a type conversion
// so that Encode serializes the string as a SimpleString.
type SimpleString string

// BulkString represents a binary-safe string as defined by the RESP.
// It can be used as a type conversion so that Encode serializes the string
// as a BulkString, but this is the default encoding for a normal Go string.
type BulkString string

// Encode encode the value v and writes the serialized data to w.
func Encode(w io.Writer, v interface{}) error {
	return encodeValue(w, v)
}

// encodeValue encodes the value v and writes the serialized data to w.
func encodeValue(w io.Writer, v interface{}) error {
	switch v := v.(type) {
	case SimpleString:
		return encodeSimpleString(w, v)
	case Error:
		return encodeError(w, v)
	case int64:
		return encodeInteger(w, v)
	case string:
		return encodeBulkString(w, BulkString(v))
	case BulkString:
		return encodeBulkString(w, v)
	case Array:
		return encodeArray(w, v)
	case nil:
		return encodeNil(w)
	default:
		return ErrInvalidValue
	}
}

// encodeArray encodes an array value to w.
func encodeArray(w io.Writer, v Array) error {
	// Special case for a nil array
	if v == nil {
		err := encodePrefixed(w, '*', "-1")
		return err
	}

	// First encode the number of elements
	n := len(v)
	err := encodePrefixed(w, '*', strconv.Itoa(n))
	if err != nil {
		return err
	}

	// Then encode each value
	for _, el := range v {
		err = encodeValue(w, el)
		if err != nil {
			return err
		}
	}
	return nil
}

// encodeBulkString encodes a bulk string to w.
func encodeBulkString(w io.Writer, v BulkString) error {
	n := len(v)
	data := strconv.Itoa(n) + "\r\n" + string(v)
	return encodePrefixed(w, '$', data)
}

// encodeInteger encodes an integer value to w.
func encodeInteger(w io.Writer, v int64) error {
	return encodePrefixed(w, ':', strconv.FormatInt(v, 10))
}

// encodeSimpleString encodes a simple string value to w.
func encodeSimpleString(w io.Writer, v SimpleString) error {
	return encodePrefixed(w, '+', string(v))
}

// encodeError encodes an error value to w.
func encodeError(w io.Writer, v Error) error {
	return encodePrefixed(w, '-', string(v))
}

// encodeNil encodes a nil value as a nil bulk string.
func encodeNil(w io.Writer) error {
	return encodePrefixed(w, '$', "-1")
}

// encodePrefixed encodes the data v to w, with the specified prefix.
func encodePrefixed(w io.Writer, prefix byte, v string) error {
	buf := make([]byte, len(v)+3)
	buf[0] = prefix
	copy(buf[1:], v)
	copy(buf[len(buf)-2:], "\r\n")
	_, err := w.Write(buf)
	return err
}
