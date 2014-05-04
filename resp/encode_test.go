package resp

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	var buf bytes.Buffer

	for i, c := range validCases {
		buf.Reset()
		err := Encode(&buf, c.val)
		if err != nil {
			t.Errorf("%d: got error %s", i, err)
			continue
		}

		if bytes.Compare(buf.Bytes(), c.enc) != 0 {
			t.Errorf("%d: expected %x (%q), got %x (%q)", i, c.enc, string(c.enc), buf.Bytes(), buf.String())
		}
	}
}
