package resp

import "errors"

var ErrMissingCRLF = errors.New("resp: missing CRLF")

type Error []byte

type SimpleString []byte

func DecodeRequest(b []byte) ([][]byte, error) {
	return nil, nil
}

func Decode(b []byte) ([]interface{}, error) {
	return nil, nil
}

func decodePayload(b []byte, isreq bool) (interface{}, error) {
	for i := 0; i < len(b); {
		ch := b[i]
		i++
		switch ch {
		case '+':
			// Simple string
			val, _, err := decodeSimpleString(b[i:])
			if err != nil {
				return nil, err
			}
			return val, nil
		case '-':
			// Error
			val, _, err := decodeError(b[i:])
			if err != nil {
				return nil, err
			}
			return val, nil
		}
	}
	return nil, nil
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
