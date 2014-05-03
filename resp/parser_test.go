package resp

import (
	"bytes"
	"io"
	"testing"
)

var cases = []struct {
	in  []byte
	out interface{}
	err error
}{
	0:  {[]byte{'+', '\r', '\n'}, SimpleString{}, nil},
	1:  {[]byte{'+', 'a', '\r', '\n'}, SimpleString{'a'}, nil},
	2:  {[]byte{'+', 'O', 'K', '\r', '\n'}, SimpleString{'O', 'K'}, nil},
	3:  {[]byte("+ceci n'est pas un string\r\n"), SimpleString("ceci n'est pas un string"), nil},
	4:  {[]byte("+ceci n'est pas un string"), SimpleString(nil), io.EOF},
	5:  {[]byte("+"), SimpleString(nil), io.EOF},
	6:  {[]byte{'-', '\r', '\n'}, Error{}, nil},
	7:  {[]byte{'-', 'a', '\r', '\n'}, Error{'a'}, nil},
	8:  {[]byte{'-', 'K', 'O', '\r', '\n'}, Error{'K', 'O'}, nil},
	9:  {[]byte("-ceci n'est pas un string\r\n"), Error("ceci n'est pas un string"), nil},
	10: {[]byte("-ceci n'est pas un string"), Error(nil), io.EOF},
	11: {[]byte("-"), Error(nil), io.EOF},
	12: {[]byte(":\r\n"), Integer(0), nil},
	13: {[]byte(":1\r\n"), Integer(1), nil},
	14: {[]byte(":123\r\n"), Integer(123), nil},
	15: {[]byte(":123\n"), Integer(0), ErrMissingCRLF},
	16: {[]byte(":123a\r\n"), Integer(0), ErrInvalidInteger},
	17: {[]byte(":-123\r\n"), Integer(-123), nil},
	18: {[]byte(":123"), Integer(0), io.EOF},
	19: {[]byte(":-1-3\r\n"), Integer(0), ErrInvalidInteger},
	20: {[]byte(":"), Integer(0), io.EOF},
	21: {[]byte("$0\r\n\r\n"), BulkString(""), nil},
	22: {[]byte("$"), BulkString(nil), io.EOF},
	23: {[]byte("$\r\n\r\n"), BulkString(""), nil},
	24: {[]byte("$24\r\nceci n'est pas un string\r\n"), BulkString("ceci n'est pas un string"), nil},
	25: {[]byte("$6\r\nc\r\n"), BulkString(nil), ErrInvalidBulkString},
	26: {[]byte("$6\r\nabc\r\n"), BulkString(nil), ErrInvalidBulkString},
	27: {[]byte("$6\nabc\r\n"), BulkString(nil), ErrMissingCRLF},
	28: {[]byte("$4\r\nabc\r\n"), BulkString(nil), ErrInvalidBulkString},
	29: {[]byte("$51\r\nceci n'est pas un string\r\navec\rdes\nsauts\r\nde\x00ligne.\r\n"), BulkString("ceci n'est pas un string\r\navec\rdes\nsauts\r\nde\x00ligne."), nil},
	30: {[]byte("$-1\r\n"), BulkString(nil), nil},
	31: {[]byte("$-3\r\n"), BulkString(nil), ErrInvalidBulkString},
	32: {[]byte("*0\r\n"), Array{}, nil},
	33: {[]byte("*1\r\n:10\r\n"), Array{Integer(10)}, nil},
	34: {[]byte("*1\n:10\r\n"), Array{}, ErrMissingCRLF},
	35: {[]byte("*-1\r\n"), Array(nil), nil},
	36: {[]byte("*-3\r\n"), Array(nil), ErrInvalidArray},
	37: {[]byte("*3\r\n+string\r\n-error\r\n:-2345\r\n"),
		Array{SimpleString("string"), Error("error"), Integer(-2345)}, nil},
	38: {[]byte("*5\r\n+string\r\n-error\r\n:-2345\r\n$4\r\nallo\r\n*2\r\n$0\r\n\r\n$-1\r\n"),
		Array{SimpleString("string"), Error("error"), Integer(-2345), BulkString("allo"),
			Array{BulkString(""), BulkString(nil)}}, nil},
}

func TestValues(t *testing.T) {
	for i, c := range cases {
		got, err := Decode(bytes.NewReader(c.in))
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.out == nil {
			continue
		}
		assertValue(t, i, got, c.out)
	}
}

func assertValue(t *testing.T, i int, got, exp interface{}) {
	switch tt := got.(type) {
	case SimpleString:
		assertSimpleString(t, i, tt, exp)
	case Error:
		assertError(t, i, tt, exp)
	case Integer:
		assertInteger(t, i, tt, exp)
	case BulkString:
		assertBulkString(t, i, tt, exp)
	case Array:
		assertArray(t, i, tt, exp)
	}
}

func assertSimpleString(t *testing.T, i int, got SimpleString, exp interface{}) {
	expv, ok := exp.(SimpleString)
	if !ok {
		t.Errorf("%d: expected a %T, got %T", i, exp, got)
		return
	}
	if bytes.Compare(got, expv) != 0 {
		t.Errorf("%d: expected output %X (%q), got %X (%q)", i, expv, string(expv), got, string(got))
	}
}

func assertError(t *testing.T, i int, got Error, exp interface{}) {
	expv, ok := exp.(Error)
	if !ok {
		t.Errorf("%d: expected a %T, got %T", i, exp, got)
		return
	}
	if bytes.Compare(got, expv) != 0 {
		t.Errorf("%d: expected output %X (%q), got %X (%q)", i, expv, string(expv), got, string(got))
	}
}

func assertInteger(t *testing.T, i int, got Integer, exp interface{}) {
	expv, ok := exp.(Integer)
	if !ok {
		t.Errorf("%d: expected a %T, got %T", i, exp, got)
		return
	}
	if expv != got {
		t.Errorf("%d: expected output %d, got %d", i, expv, got)
	}
}

func assertBulkString(t *testing.T, i int, got BulkString, exp interface{}) {
	expv, ok := exp.(BulkString)
	if !ok {
		t.Errorf("%d: expected a %T, got %T", i, exp, got)
		return
	}
	if bytes.Compare(got, expv) != 0 {
		t.Errorf("%d: expected output %X (%q), got %X (%q)", i, expv, string(expv), got, string(got))
	}
}

func assertArray(t *testing.T, i int, got Array, exp interface{}) {
	expv, ok := exp.(Array)
	if !ok {
		t.Errorf("%d: expected a %T, got %T", i, exp, got)
		return
	}
	if len(got) != len(expv) {
		t.Errorf("%d: expected an array of %d elements, got %d", i, len(expv), len(got))
		return
	}
	for j := 0; j < len(got); j++ {
		assertValue(t, i, got[j], expv[j])
	}
}

func BenchmarkSimpleString(b *testing.B) {
	r := bytes.NewReader(cases[3].in)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkError(b *testing.B) {
	r := bytes.NewReader(cases[9].in)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkInteger(b *testing.B) {
	r := bytes.NewReader(cases[17].in)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkBulkString(b *testing.B) {
	r := bytes.NewReader(cases[29].in)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkArray(b *testing.B) {
	r := bytes.NewReader(cases[37].in)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}
