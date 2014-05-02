package resp

import (
	"bytes"
	"reflect"
	"testing"
)

var simpleStrings = []struct {
	in  []byte
	out interface{}
	err error
}{
	0: {[]byte{'+', '\r', '\n'}, SimpleString{}, nil},
	1: {[]byte{'+', 'a', '\r', '\n'}, SimpleString{'a'}, nil},
	2: {[]byte{'+', 'O', 'K', '\r', '\n'}, SimpleString{'O', 'K'}, nil},
	3: {[]byte("+ceci n'est pas un string\r\n"), SimpleString("ceci n'est pas un string"), nil},
	4: {[]byte("+ceci n'est pas un string"), SimpleString(nil), ErrMissingCRLF},
	5: {[]byte("+"), SimpleString(nil), ErrMissingCRLF},
}

var errs = []struct {
	in  []byte
	out interface{}
	err error
}{
	0: {[]byte{'-', '\r', '\n'}, Error{}, nil},
	1: {[]byte{'-', 'a', '\r', '\n'}, Error{'a'}, nil},
	2: {[]byte{'-', 'K', 'O', '\r', '\n'}, Error{'K', 'O'}, nil},
	3: {[]byte("-ceci n'est pas un string\r\n"), Error("ceci n'est pas un string"), nil},
	4: {[]byte("-ceci n'est pas un string"), Error(nil), ErrMissingCRLF},
	5: {[]byte("-"), Error(nil), ErrMissingCRLF},
}

var integers = []struct {
	in  []byte
	out interface{}
	err error
}{
	0: {[]byte(":\r\n"), Integer(0), nil},
	1: {[]byte(":1\r\n"), Integer(1), nil},
	2: {[]byte(":123\r\n"), Integer(123), nil},
	3: {[]byte(":123\n"), Integer(0), ErrMissingCRLF},
	4: {[]byte(":123a\r\n"), Integer(0), ErrInvalidInteger},
	5: {[]byte(":-123\r\n"), Integer(-123), nil},
	6: {[]byte(":123"), Integer(0), ErrMissingCRLF},
	7: {[]byte(":-1-3\r\n"), Integer(0), ErrInvalidInteger},
	8: {[]byte(":"), Integer(0), ErrMissingCRLF},
}

var bulkStrings = []struct {
	in  []byte
	out interface{}
	err error
}{
	0:  {[]byte("$0\r\n\r\n"), BulkString(""), nil},
	1:  {[]byte("$"), BulkString(nil), ErrMissingCRLF},
	2:  {[]byte("$\r\n\r\n"), BulkString(""), nil},
	3:  {[]byte("$24\r\nceci n'est pas un string\r\n"), BulkString("ceci n'est pas un string"), nil},
	4:  {[]byte("$6\r\nc\r\n"), BulkString(nil), ErrInvalidBulkString},
	5:  {[]byte("$6\r\nabc\r\n"), BulkString(nil), ErrInvalidBulkString},
	6:  {[]byte("$6\nabc\r\n"), BulkString(nil), ErrMissingCRLF},
	7:  {[]byte("$4\r\nabc\r\n"), BulkString(nil), ErrInvalidBulkString},
	8:  {[]byte("$51\r\nceci n'est pas un string\r\navec\rdes\nsauts\r\nde\x00ligne.\r\n"), BulkString("ceci n'est pas un string\r\navec\rdes\nsauts\r\nde\x00ligne."), nil},
	9:  {[]byte("$-1\r\n"), BulkString(nil), nil},
	10: {[]byte("$-3\r\n"), BulkString(nil), ErrInvalidBulkString},
}

var arrays = []struct {
	in  []byte
	out interface{}
	err error
}{
	0: {[]byte("*0\r\n"), Array{}, nil},
	1: {[]byte("*1\r\n:10\r\n"), Array{Integer(10)}, nil},
	2: {[]byte("*1\n:10\r\n"), Array{}, ErrMissingCRLF},
	3: {[]byte("*-1\r\n"), Array(nil), nil},
	4: {[]byte("*-3\r\n"), Array(nil), ErrInvalidArray},
	5: {[]byte("*3\r\n+string\r\n-error\r\n:-2345\r\n"), Array{SimpleString("string"), Error("error"), Integer(-2345)}, nil},
}

func TestSimpleString(t *testing.T) {
	for i, c := range simpleStrings {
		got, _, err := decodeValue(c.in)
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.out == nil {
			continue
		}
		if ss, ok := got.(SimpleString); !ok {
			t.Errorf("%d: expected a simple string, got %T", i, got)
		} else if bytes.Compare(ss, c.out.(SimpleString)) != 0 {
			t.Errorf("%d: expected output %X (%q), got %X (%q)", i, c.out, string(c.out.(SimpleString)), ss, string(ss))
		}
	}
}

func TestError(t *testing.T) {
	for i, c := range errs {
		got, _, err := decodeValue(c.in)
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.out == nil {
			continue
		}
		if er, ok := got.(Error); !ok {
			t.Errorf("%d: expected an error, got %T", i, got)
		} else if bytes.Compare(er, c.out.(Error)) != 0 {
			t.Errorf("%d: expected output %X (%q), got %X (%q)", i, c.out, string(c.out.(Error)), er, string(er))
		}
	}
}

func TestInteger(t *testing.T) {
	for i, c := range integers {
		got, _, err := decodeValue(c.in)
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

func TestBulkString(t *testing.T) {
	for i, c := range bulkStrings {
		got, _, err := decodeValue(c.in)
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.out == nil {
			continue
		}
		if bs, ok := got.(BulkString); !ok {
			t.Errorf("%d: expected a bulk string, got %T", i, got)
		} else if bytes.Compare(bs, c.out.(BulkString)) != 0 {
			t.Errorf("%d: expected output %X (%q), got %X (%q)", i, c.out, string(c.out.(BulkString)), bs, string(bs))
		}
	}
}

func TestArray(t *testing.T) {
	for i, c := range arrays {
		got, _, err := decodeValue(c.in)
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.out == nil {
			continue
		}
		if ar, ok := got.(Array); !ok {
			t.Errorf("%d: expected an array, got %T", i, got)
		} else {
			expar := c.out.(Array)
			if len(ar) != len(expar) {
				t.Errorf("%d: expected an array of %d elements, got %d", i, len(expar), len(ar))
				continue
			}
			for j := range ar {
				art := reflect.TypeOf(ar[j])
				expt := reflect.TypeOf(expar[j])
				if art != expt {
					t.Errorf("%d: at %d, expected type %T, got %T", i, j, expar[j], ar[j])
					continue
				}
				switch ar[j].(type) {
				case SimpleString:
					b1, b2 := []byte(ar[j].(SimpleString)), []byte(expar[j].(SimpleString))
					if bytes.Compare(b1, b2) != 0 {
						t.Errorf("%d: at %d, expected %x (%q), got %x (%q)", i, j, b2, string(b2), b1, string(b1))
					}
				case Error:
					b1, b2 := []byte(ar[j].(Error)), []byte(expar[j].(Error))
					if bytes.Compare(b1, b2) != 0 {
						t.Errorf("%d: at %d, expected %x (%q), got %x (%q)", i, j, b2, string(b2), b1, string(b1))
					}
				case BulkString:
					b1, b2 := []byte(ar[j].(BulkString)), []byte(expar[j].(BulkString))
					if bytes.Compare(b1, b2) != 0 {
						t.Errorf("%d: at %d, expected %x (%q), got %x (%q)", i, j, b2, string(b2), b1, string(b1))
					}
				case Integer:
					i1, i2 := ar[j], expar[j]
					if i1 != i2 {
						t.Errorf("%d: at %d, expected %d, got %d", i, j, i2, i1)
					}
				case Array:
					// TODO: Extract comparison func to call recursively on arrays
				}
			}
		}
	}
}

func BenchmarkSimpleString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decodeValue(simpleStrings[3].in)
	}
}

func BenchmarkError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decodeValue(errs[3].in)
	}
}

func BenchmarkInteger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decodeValue(integers[5].in)
	}
}

func BenchmarkBulkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decodeValue(bulkStrings[8].in)
	}
}

func BenchmarkArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decodeValue(arrays[5].in)
	}
}
