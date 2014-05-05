// Package resp implements an efficient decoder for the Redis Serialization Protocol (RESP).
//
// See http://redis.io/topics/protocol for the reference.
package resp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrInvalidPrefix is returned if the data contains an unrecognized prefix.
	ErrInvalidPrefix = errors.New("resp: invalid prefix")

	// ErrMissingCRLF is returned if a \r\n is missing in the data slice.
	ErrMissingCRLF = errors.New("resp: missing CRLF")

	// ErrInvalidInteger is returned if an invalid character is found while parsing an integer.
	ErrInvalidInteger = errors.New("resp: invalid integer character")

	// ErrInvalidBulkString is returned if the bulk string data cannot be decoded.
	ErrInvalidBulkString = errors.New("resp: invalid bulk string")

	// ErrInvalidArray is returned if the array data cannot be decoded.
	ErrInvalidArray = errors.New("resp: invalid array")

	// ErrNotAnArray is returned if the DecodeRequest function is called and
	// the decoded value is not an array.
	ErrNotAnArray = errors.New("resp: expected an array type")

	// ErrInvalidRequest is returned if the DecodeRequest function is called and
	// the decoded value is not an array containing only bulk strings, and at least 1 element.
	ErrInvalidRequest = errors.New("resp: invalid request, must be an array of bulk strings with at least one element")
)

// Integer represents a signed, 64-bit integer as defined by the RESP.
type Integer int64

func (i Integer) String() string {
	return fmt.Sprintf("%d", i)
}

// Error represents an error string as defined by the RESP. It cannot
// contain \r or \n characters.
type Error string

// SimpleString represents a simple string as defined by the RESP. It
// cannot contain \r or \n characters.
type SimpleString string

// BulkString represents a binary-safe string as defined by the RESP.
type BulkString string

// Array represents an array of values, as defined by the RESP.
type Array []interface{}

func (a Array) String() string {
	var buf bytes.Buffer
	for i, v := range a {
		buf.WriteString(fmt.Sprintf("[%2d] %[2]s (%[2]T)\n", i, v))
	}
	return buf.String()
}

// DecodeRequest decodes the provided byte slice and returns the array
// representing the request. If the encoded value is not an array, it
// returns ErrNotAnArray, and if it is not a valid request, it returns ErrInvalidRequest.
func DecodeRequest(r io.Reader) ([]string, error) {
	// Decode the value
	val, err := Decode(r)
	if err != nil {
		return nil, err
	}

	// Must be an array
	ar, ok := val.(Array)
	if !ok {
		return nil, ErrNotAnArray
	}

	// Must have at least one element
	if len(ar) < 1 {
		return nil, ErrInvalidRequest
	}

	// Must have only strings
	strs := make([]string, len(ar))
	for i, v := range ar {
		if v, ok := v.(string); !ok {
			return nil, ErrInvalidRequest
		} else {
			strs[i] = v
		}
	}
	return strs, nil
}

// Decode decodes the provided byte slice and returns the parsed value.
func Decode(r io.Reader) (interface{}, error) {
	br := bufio.NewReader(r)
	val, _, err := decodeValue(br)
	return val, err
}

// decodeValue parses the byte slice and decodes the value based on its
// prefix, as defined by the RESP protocol.
func decodeValue(r *bufio.Reader) (val interface{}, n int, err error) {
	ch, err := r.ReadByte()
	if err != nil {
		return val, 0, err
	}
	switch ch {
	case '+':
		// Simple string
		val, n, err = decodeSimpleString(r)
	case '-':
		// Error
		val, n, err = decodeError(r)
	case ':':
		// Integer
		val, n, err = decodeInteger(r)
	case '$':
		// Bulk string
		val, n, err = decodeBulkString(r)
	case '*':
		// Array
		val, n, err = decodeArray(r)
	default:
		err = ErrInvalidPrefix
	}

	// n+1 for the prefix consumed by this func
	return val, n + 1, err
}

// decodeArray decodes the byte slice as an array. It assumes the
// '*' prefix is already consumed.
func decodeArray(r *bufio.Reader) (Array, int, error) {
	// First comes the number of elements in the array
	cnt, n, err := decodeInteger(r)
	if err != nil {
		return nil, n, err
	}
	switch {
	case cnt == -1:
		// Nil array
		return nil, n, nil

	case cnt == 0:
		// Empty, but allocated, array
		return Array{}, n, nil

	case cnt < 0:
		// Invalid length
		return nil, n, ErrInvalidArray

	default:
		// Allocate the array
		ar := make(Array, cnt)

		// Decode each value
		for i := 0; i < int(cnt); i++ {
			val, nn, err := decodeValue(r)
			n += nn
			if err != nil {
				return nil, n, err
			}
			ar[i] = val
		}
		return ar, n, nil
	}
}

// decodeBulkString decodes the byte slice as a binary-safe string. The
// '$' prefix is assumed to be already consumed.
func decodeBulkString(r *bufio.Reader) (interface{}, int, error) {
	// First comes the length of the bulk string, an integer
	cnt, n, err := decodeInteger(r)
	if err != nil {
		return nil, n, err
	}
	switch {
	case cnt == -1:
		// Special case to represent a nil bulk string
		return nil, n, nil

	case cnt < -1:
		return nil, n, ErrInvalidBulkString

	default:
		// Then the string is cnt long, and bytes read is cnt+n+2 (for ending CRLF)
		buf := make([]byte, cnt+2)
		nb, err := r.Read(buf)
		if nb < int(cnt)+2 {
			return nil, n + nb, ErrInvalidBulkString
		}
		return string(buf[:nb-2]), nb + n, err
	}
}

// decodeInteger decodes the byte slice as a singed 64bit integer. The
// ':' prefix is assumed to be already consumed.
func decodeInteger(r *bufio.Reader) (val int64, n int, err error) {
	var cr bool
	var sign int64 = 1

loop:
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return 0, n, err
		}
		n++

		switch ch {
		case '\r':
			cr = true
			break loop

		case '\n':
			break loop

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			val = val*10 + int64(ch-'0')

		case '-':
			if n == 1 {
				sign = -1
				continue
			}
			fallthrough
		default:
			return 0, n, ErrInvalidInteger
		}
	}

	if !cr {
		return 0, n, ErrMissingCRLF
	}
	// Presume next byte was \n
	r.ReadByte()
	return sign * val, n + 1, nil
}

// decodeSimpleString decodes the byte slice as a SimpleString. The
// '+' prefix is assumed to be already consumed.
func decodeSimpleString(r *bufio.Reader) (interface{}, int, error) {
	v, err := r.ReadBytes('\r')
	if err != nil {
		return nil, len(v), err
	}
	// Presume next byte was \n
	r.ReadByte()
	return string(v[:len(v)-1]), len(v) + 1, nil
}

// decodeError decodes the byte slice as an Error. The '-' prefix
// is assumed to be already consumed.
func decodeError(r *bufio.Reader) (interface{}, int, error) {
	return decodeSimpleString(r)
}
