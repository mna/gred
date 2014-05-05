package resp

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	var buf bytes.Buffer
	var err error

	for i, c := range validCases {
		buf.Reset()
		// Convert to correct value, based on expectation
		switch c.enc[0] {
		case '+':
			err = Encode(&buf, SimpleString(c.val.(string)))
		case '-':
			err = Encode(&buf, Error(c.val.(string)))
		case '$':
			if c.val != nil {
				err = Encode(&buf, BulkString(c.val.(string)))
			} else {
				err = Encode(&buf, c.val)
			}
		case ':':
			err = Encode(&buf, Integer(c.val.(int64)))
		default:
			err = Encode(&buf, c.val)
		}
		if err != nil {
			t.Errorf("%d: got error %s", i, err)
			continue
		}

		if bytes.Compare(buf.Bytes(), c.enc) != 0 {
			t.Errorf("%d: expected %x (%q), got %x (%q)", i, c.enc, string(c.enc), buf.Bytes(), buf.String())
		}
	}
}

func BenchmarkEncodeSimpleString(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		Encode(&buf, validCases[3].val)
	}
}

func BenchmarkEncodeError(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		Encode(&buf, validCases[7].val)
	}
}

func BenchmarkEncodeInteger(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		Encode(&buf, validCases[10].val)
	}
}

func BenchmarkEncodeBulkString(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		Encode(&buf, validCases[13].val)
	}
}

func BenchmarkEncodeArray(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		Encode(&buf, validCases[19].val)
	}
}
