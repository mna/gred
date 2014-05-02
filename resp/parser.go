package resp

import "errors"

var (
	ErrInvalidInteger = errors.New("resp: invalid integer character")
	ErrMissingCRLF    = errors.New("resp: missing CRLF")
)

type Integer int64

type Error []byte

type SimpleString []byte

func DecodeRequest(b []byte) ([][]byte, error) {
	return nil, nil
}

func Decode(b []byte) ([]interface{}, error) {
	return nil, nil
}

func decodePayload(b []byte, isreq bool) (interface{}, error) {
	var val interface{}
	var err error

	for i := 0; i < len(b); {
		ch := b[i]
		i++
		switch ch {
		case '+':
			// Simple string
			val, _, err = decodeSimpleString(b[i:])
		case '-':
			// Error
			val, _, err = decodeError(b[i:])
		case ':':
			// Integer
			val, _, err = decodeInteger(b[i:])
		}
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	return nil, nil
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
			// Will return with ErrMissingCRLF
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
	return sign * val, n, nil
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
	return SimpleString(b[:end]), n, nil
}

// decodeError decodes the byte slice as an Error. The '-' prefix
// is assumed to be already consumed.
func decodeError(b []byte) (Error, int, error) {
	val, n, err := decodeSimpleString(b)
	return Error(val), n, err
}
