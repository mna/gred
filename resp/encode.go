package resp

import (
	"errors"
	"io"
	"strconv"
)

var ErrInvalidValue = errors.New("resp: invalid value")

func Encode(w io.Writer, v interface{}) error {
	return encodeValue(w, v)
}

func encodeValue(w io.Writer, v interface{}) error {
	switch v := v.(type) {
	case SimpleString:
		return encodeSimpleString(w, v)
	case Error:
		return encodeError(w, v)
	case Integer:
		return encodeInteger(w, v)
	case BulkString:
		return encodeBulkString(w, v)
	case Array:
		return encodeArray(w, v)
	default:
		return ErrInvalidValue
	}
}

func encodeArray(w io.Writer, v Array) error {
	// Special case for a nil array
	if v == nil {
		err := encodePrefixed(w, '*', []byte("-1"))
		return err
	}

	// First encode the number of elements
	n := len(v)
	err := encodePrefixed(w, '*', []byte(strconv.Itoa(n)))
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

func encodeBulkString(w io.Writer, v BulkString) error {
	// Special case for a nil bulk string
	if v == nil {
		err := encodePrefixed(w, '$', []byte("-1"))
		return err
	}

	// First encode the length
	n := len(v)
	err := encodePrefixed(w, '$', []byte(strconv.Itoa(n)))
	if err != nil {
		return err
	}
	// Then the string
	_, err = w.Write(append(v, '\r', '\n'))
	return err
}

func encodeInteger(w io.Writer, v Integer) error {
	return encodePrefixed(w, ':', []byte(strconv.FormatInt(int64(v), 10)))
}

func encodeSimpleString(w io.Writer, v SimpleString) error {
	return encodePrefixed(w, '+', v)
}

func encodeError(w io.Writer, v Error) error {
	return encodePrefixed(w, '-', v)
}

func encodePrefixed(w io.Writer, prefix byte, v []byte) error {
	buf := make([]byte, len(v)+3)
	buf[0] = prefix
	copy(buf[1:], v)
	copy(buf[len(buf)-2:], "\r\n")
	_, err := w.Write(buf)
	return err
}
