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
	4:  {[]byte(":123\n"), int64(0), ErrMissingCRLF},
	5:  {[]byte(":123a\r\n"), int64(0), ErrInvalidInteger},
	6:  {[]byte(":123"), int64(0), io.EOF},
	7:  {[]byte(":-1-3\r\n"), int64(0), ErrInvalidInteger},
	8:  {[]byte(":"), int64(0), io.EOF},
	9:  {[]byte("$"), nil, io.EOF},
	10: {[]byte("$6\r\nc\r\n"), nil, ErrInvalidBulkString},
	11: {[]byte("$6\r\nabc\r\n"), nil, ErrInvalidBulkString},
	12: {[]byte("$6\nabc\r\n"), nil, ErrMissingCRLF},
	13: {[]byte("$4\r\nabc\r\n"), nil, ErrInvalidBulkString},
	14: {[]byte("$-3\r\n"), nil, ErrInvalidBulkString},
	15: {[]byte("*1\n:10\r\n"), Array{}, ErrMissingCRLF},
	16: {[]byte("*-3\r\n"), Array(nil), ErrInvalidArray},
	17: {[]byte(":\r\n"), int64(0), nil},
	18: {[]byte("$\r\n\r\n"), "", nil},
}

var validCases = []struct {
	enc []byte
	val interface{}
	err error
}{
	0:  {[]byte{'+', '\r', '\n'}, "", nil},
	1:  {[]byte{'+', 'a', '\r', '\n'}, "a", nil},
	2:  {[]byte{'+', 'O', 'K', '\r', '\n'}, "OK", nil},
	3:  {[]byte("+ceci n'est pas un string\r\n"), "ceci n'est pas un string", nil},
	4:  {[]byte{'-', '\r', '\n'}, "", nil},
	5:  {[]byte{'-', 'a', '\r', '\n'}, "a", nil},
	6:  {[]byte{'-', 'K', 'O', '\r', '\n'}, "KO", nil},
	7:  {[]byte("-ceci n'est pas un string\r\n"), "ceci n'est pas un string", nil},
	8:  {[]byte(":1\r\n"), int64(1), nil},
	9:  {[]byte(":123\r\n"), int64(123), nil},
	10: {[]byte(":-123\r\n"), int64(-123), nil},
	11: {[]byte("$0\r\n\r\n"), "", nil},
	12: {[]byte("$24\r\nceci n'est pas un string\r\n"), "ceci n'est pas un string", nil},
	13: {[]byte("$51\r\nceci n'est pas un string\r\navec\rdes\nsauts\r\nde\x00ligne.\r\n"), "ceci n'est pas un string\r\navec\rdes\nsauts\r\nde\x00ligne.", nil},
	14: {[]byte("$-1\r\n"), nil, nil},
	15: {[]byte("*0\r\n"), Array{}, nil},
	16: {[]byte("*1\r\n:10\r\n"), Array{int64(10)}, nil},
	17: {[]byte("*-1\r\n"), Array(nil), nil},
	18: {[]byte("*3\r\n+string\r\n-error\r\n:-2345\r\n"),
		Array{"string", "error", int64(-2345)}, nil},
	19: {[]byte("*5\r\n+string\r\n-error\r\n:-2345\r\n$4\r\nallo\r\n*2\r\n$0\r\n\r\n$-1\r\n"),
		Array{"string", "error", int64(-2345), "allo",
			Array{"", nil}}, nil},
}

func TestDecode(t *testing.T) {
	for i, c := range append(validCases, decodeErrCases...) {
		got, err := Decode(bytes.NewBuffer(c.enc))
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.val == nil {
			continue
		}
		assertValue(t, i, got, c.val)
	}
}

func TestDecodeRequest(t *testing.T) {
	cases := []struct {
		raw []byte
		exp []string
		err error
	}{
		0: {[]byte("*-1\r\n"), nil, ErrInvalidRequest},
		1: {[]byte(":4\r\n"), nil, ErrNotAnArray},
		2: {[]byte("*0\r\n"), nil, ErrInvalidRequest},
		3: {[]byte("*1\r\n:6\r\n"), nil, ErrInvalidRequest},
		4: {[]byte("*1\r\n$2\r\nab\r\n"), []string{"ab"}, nil},
	}

	for i, c := range cases {
		buf := bytes.NewBuffer(c.raw)
		got, err := DecodeRequest(buf)
		if err != c.err {
			t.Errorf("%d: expected error %v, got %v", i, c.err, err)
		}
		if got == nil && c.exp == nil {
			continue
		}
		assertValue(t, i, got, c.exp)
	}
}

func assertValue(t *testing.T, i int, got, exp interface{}) {
	switch tt := got.(type) {
	case string:
		assertString(t, i, tt, exp)
	case int64:
		assertInteger(t, i, tt, exp)
	case Array:
		assertArray(t, i, tt, exp)
	case nil:
		if exp != nil {
			t.Errorf("%d: expected nil, got %v", i, exp)
		}
	default:
		t.Errorf("%d: unknown value type %T", i, got)
	}
}

func assertString(t *testing.T, i int, got string, exp interface{}) {
	expv, ok := exp.(string)
	if !ok {
		t.Errorf("%d: expected a %T, got %T", i, exp, got)
		return
	}
	if got != expv {
		t.Errorf("%d: expected output %X (%q), got %X (%q)", i, expv, string(expv), got, string(got))
	}
}

func assertInteger(t *testing.T, i int, got int64, exp interface{}) {
	expv, ok := exp.(int64)
	if !ok {
		t.Errorf("%d: expected a %T, got %T", i, exp, got)
		return
	}
	if expv != got {
		t.Errorf("%d: expected output %d, got %d", i, expv, got)
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
	r := bytes.NewBuffer(validCases[3].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkDecodeError(b *testing.B) {
	r := bytes.NewBuffer(validCases[7].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkDecodeInteger(b *testing.B) {
	r := bytes.NewBuffer(validCases[10].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkDecodeBulkString(b *testing.B) {
	r := bytes.NewBuffer(validCases[13].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}

func BenchmarkDecodeArray(b *testing.B) {
	r := bytes.NewBuffer(validCases[19].enc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(r)
	}
}
