package resp

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	var buf bytes.Buffer

	for i, c := range cases {
		buf.Reset()
		if c.err == nil {
			err := Encode(&buf, c.out)
			if err != nil {
				t.Errorf("%d: got error %s", i, err)
				continue
			}

			if bytes.Compare(buf.Bytes(), c.in) != 0 {
				t.Errorf("%d: expected %x (%q), got %x (%q)", i, c.in, string(c.in), buf.Bytes(), buf.String())
			}
		}
	}
}
