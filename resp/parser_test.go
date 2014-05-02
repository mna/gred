package resp

import (
	"bytes"
	"testing"
)

var simpleStrings = []struct {
	in  []byte
	out interface{}
	err error
}{
	0: {[]byte{'+', '\r', '\n'}, []byte{}, nil},
	1: {[]byte{'+', 'a', '\r', '\n'}, []byte{'a'}, nil},
	2: {[]byte{'+', 'O', 'K', '\r', '\n'}, []byte{'O', 'K'}, nil},
	3: {[]byte("+ceci n'est pas un string\r\n"), []byte("ceci n'est pas un string"), nil},
	4: {[]byte("+ceci n'est pas un string"), nil, ErrMissingCRLF},
	5: {[]byte("+"), nil, ErrMissingCRLF},
}

var errs = []struct {
	in  []byte
	out interface{}
	err error
}{
	0: {[]byte{'-', '\r', '\n'}, []byte{}, nil},
	1: {[]byte{'-', 'a', '\r', '\n'}, []byte{'a'}, nil},
	2: {[]byte{'-', 'K', 'O', '\r', '\n'}, []byte{'K', 'O'}, nil},
	3: {[]byte("-ceci n'est pas un string\r\n"), []byte("ceci n'est pas un string"), nil},
	4: {[]byte("-ceci n'est pas un string"), nil, ErrMissingCRLF},
	5: {[]byte("-"), nil, ErrMissingCRLF},
}

var integers = []struct {
	in  []byte
	out interface{}
	err error
}{
	0: {[]byte(":\r\n"), Integer(0), nil},
	1: {[]byte(":1\r\n"), Integer(1), nil},
	2: {[]byte(":123\r\n"), Integer(123), nil},
	3: {[]byte(":123\n"), nil, ErrMissingCRLF},
	4: {[]byte(":123a\r\n"), nil, ErrInvalidInteger},
	5: {[]byte(":-123\r\n"), Integer(-123), nil},
	6: {[]byte(":123"), nil, ErrMissingCRLF},
	7: {[]byte(":-1-3\r\n"), nil, ErrInvalidInteger},
}

func TestSimpleString(t *testing.T) {
	for i, c := range simpleStrings {
		got, err := decodePayload(c.in, false)
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.out == nil {
			continue
		}
		if ss, ok := got.(SimpleString); !ok {
			t.Errorf("%d: expected a simple string, got %T", i, got)
		} else if bytes.Compare(ss, c.out.([]byte)) != 0 {
			t.Errorf("%d: expected output %X (%q), got %X (%q)", i, c.out, string(c.out.([]byte)), ss, string(ss))
		}
	}
}

func TestError(t *testing.T) {
	for i, c := range errs {
		got, err := decodePayload(c.in, false)
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.out == nil {
			continue
		}
		if er, ok := got.(Error); !ok {
			t.Errorf("%d: expected an error, got %T", i, got)
		} else if bytes.Compare(er, c.out.([]byte)) != 0 {
			t.Errorf("%d: expected output %X (%q), got %X (%q)", i, c.out, string(c.out.([]byte)), er, string(er))
		}
	}
}

func TestInteger(t *testing.T) {
	for i, c := range integers {
		got, err := decodePayload(c.in, false)
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.out == nil {
			continue
		}
		if ii, ok := got.(Integer); !ok {
			t.Errorf("%d: expected an integer, got %T", i, got)
		} else if c.out.(Integer) != ii {
			t.Errorf("%d: expected output %d, got %d", i, c.out, ii)
		}
	}
}

func BenchmarkSimpleString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decodePayload(simpleStrings[3].in, false)
	}
}

func BenchmarkError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decodePayload(errs[3].in, false)
	}
}
