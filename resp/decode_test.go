package resp

import (
	"bytes"
	"io"
	"testing"
)

var decodeErrCases = []struct {
	enc []byte
	val interface{}
	err error
}{
	0:  {[]byte("+ceci n'est pas un string"), nil, io.EOF},
	1:  {[]byte("+"), nil, io.EOF},
	2:  {[]byte("-ceci n'est pas un string"), nil, io.EOF},
	3:  {[]byte("-"), nil, io.EOF},
	4:  {[]byte(":123\n"), Integer(0), ErrMissingCRLF},
	5:  {[]byte(":123a\r\n"), Integer(0), ErrInvalidInteger},
	6:  {[]byte(":123"), Integer(0), io.EOF},
	7:  {[]byte(":-1-3\r\n"), Integer(0), ErrInvalidInteger},
	8:  {[]byte(":"), Integer(0), io.EOF},
	9:  {[]byte("$"), nil, io.EOF},
	10: {[]byte("$6\r\nc\r\n"), nil, ErrInvalidBulkString},
	11: {[]byte("$6\r\nabc\r\n"), nil, ErrInvalidBulkString},
	12: {[]byte("$6\nabc\r\n"), nil, ErrMissingCRLF},
	13: {[]byte("$4\r\nabc\r\n"), nil, ErrInvalidBulkString},
	14: {[]byte("$-3\r\n"), nil, ErrInvalidBulkString},
	15: {[]byte("*1\n:10\r\n"), Array{}, ErrMissingCRLF},
	16: {[]byte("*-3\r\n"), Array(nil), ErrInvalidArray},
	17: {[]byte(":\r\n"), Integer(0), nil},
	18: {[]byte("$\r\n\r\n"), BulkString(""), nil},
}

var validCases = []struct {
	enc []byte
	val interface{}
	err error
}{
	0:  {[]byte{'+', '\r', '\n'}, SimpleString(""), nil},
	1:  {[]byte{'+', 'a', '\r', '\n'}, SimpleString("a"), nil},
	2:  {[]byte{'+', 'O', 'K', '\r', '\n'}, SimpleString("OK"), nil},
	3:  {[]byte("+ceci n'est pas un string\r\n"), SimpleString("ceci n'est pas un string"), nil},
	4:  {[]byte{'-', '\r', '\n'}, Error(""), nil},
	5:  {[]byte{'-', 'a', '\r', '\n'}, Error("a"), nil},
	6:  {[]byte{'-', 'K', 'O', '\r', '\n'}, Error("KO"), nil},
	7:  {[]byte("-ceci n'est pas un string\r\n"), Error("ceci n'est pas un string"), nil},
	8:  {[]byte(":1\r\n"), Integer(1), nil},
	9:  {[]byte(":123\r\n"), Integer(123), nil},
	10: {[]byte(":-123\r\n"), Integer(-123), nil},
	11: {[]byte("$0\r\n\r\n"), BulkString(""), nil},
	12: {[]byte("$24\r\nceci n'est pas un string\r\n"), BulkString("ceci n'est pas un string"), nil},
	13: {[]byte("$51\r\nceci n'est pas un string\r\navec\rdes\nsauts\r\nde\x00ligne.\r\n"), BulkString("ceci n'est pas un string\r\navec\rdes\nsauts\r\nde\x00ligne."), nil},
	14: {[]byte("$-1\r\n"), nil, nil},
	15: {[]byte("*0\r\n"), Array{}, nil},
	16: {[]byte("*1\r\n:10\r\n"), Array{Integer(10)}, nil},
	17: {[]byte("*-1\r\n"), Array(nil), nil},
	18: {[]byte("*3\r\n+string\r\n-error\r\n:-2345\r\n"),
		Array{SimpleString("string"), Error("error"), Integer(-2345)}, nil},
	19: {[]byte("*5\r\n+string\r\n-error\r\n:-2345\r\n$4\r\nallo\r\n*2\r\n$0\r\n\r\n$-1\r\n"),
		Array{SimpleString("string"), Error("error"), Integer(-2345), BulkString("allo"),
			Array{BulkString(""), nil}}, nil},
}

func TestDecode(t *testing.T) {
	for i, c := range append(validCases, decodeErrCases...) {
		got, err := Decode(bytes.NewReader(c.enc))
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.val == nil {
			continue
		}
		assertValue(t, i, got, c.val)
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
	if got != expv {
		t.Errorf("%d: expected output %X (%q), got %X (%q)", i, expv, string(expv), got, string(got))
	}
}

func assertError(t *testing.T, i int, got Error, exp interface{}) {
	expv, ok := exp.(Error)
	if !ok {
		t.Errorf("%d: expected a %T, got %T", i, exp, got)
		return
	}
	if got != expv {
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
	if got != expv {
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

func BenchmarkDecodeSimpleString(b *testing.B) {
	r := bytes.NewReader(validCases[3].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkDecodeError(b *testing.B) {
	r := bytes.NewReader(validCases[7].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkDecodeInteger(b *testing.B) {
	r := bytes.NewReader(validCases[10].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkDecodeBulkString(b *testing.B) {
	r := bytes.NewReader(validCases[13].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkDecodeArray(b *testing.B) {
	r := bytes.NewReader(validCases[19].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}
