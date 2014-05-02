package resp

import "errors"

var (
	ErrNoData            = errors.New("resp: no data")
	ErrInvalidPrefix     = errors.New("resp: invalid prefix")
	ErrMissingCRLF       = errors.New("resp: missing CRLF")
	ErrInvalidInteger    = errors.New("resp: invalid integer character")
	ErrInvalidBulkString = errors.New("resp: invalid bulk string")
	ErrInvalidArray      = errors.New("resp: invalid array")
)

type Integer int64

type Error []byte

type SimpleString []byte

type BulkString []byte

type Array []interface{}

func DecodeRequest(b []byte) ([][]byte, error) {
	return nil, nil
}

func Decode(b []byte) ([]interface{}, error) {
	return nil, nil
}

func decodeValue(b []byte) (val interface{}, n int, err error) {
	if len(b) == 0 {
		return nil, 0, ErrNoData
	}

	switch b[0] {
	case '+':
		// Simple string
		val, n, err = decodeSimpleString(b[1:])
	case '-':
		// Error
		val, n, err = decodeError(b[1:])
	case ':':
		// Integer
		val, n, err = decodeInteger(b[1:])
	case '$':
		// Bulk string
		val, n, err = decodeBulkString(b[1:])
	case '*':
		// Array
		val, n, err = decodeArray(b[1:])
	default:
		err = ErrInvalidPrefix
	}

	// n+1 for the prefix consumed by this func
	return val, n + 1, err
}

func decodeArray(b []byte) (Array, int, error) {
	// First comes the number of elements in the array
	cnt, n, err := decodeInteger(b)
	if err != nil {
		return nil, n, err
	}
	switch {
	case cnt == -1:
		// Nil array
		return Array(nil), n, nil

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
			val, nn, err := decodeValue(b[n:])
			n += nn
			if err != nil {
				return nil, n, err
			}
			ar[i] = val
		}
		return ar, n, nil
	}
}

func decodeBulkString(b []byte) (BulkString, int, error) {
	// First comes the length of the bulk string, an integer
	cnt, n, err := decodeInteger(b)
	if err != nil {
		return nil, n, err
	}
	switch {
	case cnt == -1:
		// Special case to represent a nil bulk string
		return BulkString(nil), n, nil

	case cnt < -1:
		return nil, n, ErrInvalidBulkString

	case len(b) < int(cnt)+n+2:
		return nil, n, ErrInvalidBulkString

	default:
		// Then the string is cnt long, and bytes read is cnt+n+2 (for ending CRLF)
		return BulkString(b[n : int(cnt)+n]), int(cnt) + n + 2, nil
	}
}

// decodeInteger decodes the byte slice as a singed 64bit integer. The
// ':' prefix is assumed to be already consumed.
func decodeInteger(b []byte) (val Integer, n int, err error) {
	var cr bool
	var sign Integer = 1

loop:
	for i := 0; i < len(b); i++ {
		ch := b[i]
		n++

		switch ch {
		case '\r':
			cr = true
			break loop

		case '\n':
			break loop

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			val = val*10 + Integer(ch-'0')

		case '-':
			if i == 0 {
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
	return sign * val, n + 1, nil
}

// decodeSimpleString decodes the byte slice as a SimpleString. The
// '+' prefix is assumed to be already consumed.
func decodeSimpleString(b []byte) (SimpleString, int, error) {
	end := -1
	n := 0
	for i := 0; i < len(b); i++ {
		ch := b[i]
		n++
		if ch == '\r' {
			// Simple strings cannot contain \r nor \n, so at the first \r we know
			// the string is over.
			end = i
			break
		}
	}
	if end == -1 {
		return nil, n, ErrMissingCRLF
	}
	// Presume next byte was \n
	return SimpleString(b[:end]), n + 1, nil
}

// decodeError decodes the byte slice as an Error. The '-' prefix
// is assumed to be already consumed.
func decodeError(b []byte) (Error, int, error) {
	val, n, err := decodeSimpleString(b)
	return Error(val), n, err
}
