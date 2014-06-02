package types

import "testing"

func TestStringSetRange(t *testing.T) {
	cases := []struct {
		src string
		ofs int64
		set string
		exp string
	}{
		0: {"", 0, "", ""},
		1: {"", 19, "", ""},
		2: {"", 3, "abc", "\x00\x00\x00abc"},
		3: {"abc", 3, "def", "abcdef"},
		4: {"abc", 1, "def", "adef"},
		5: {"abc", 0, "def", "def"},
		6: {"abc", 0, "de", "dec"},
		7: {"abc", 10, "def", "abc\x00\x00\x00\x00\x00\x00\x00def"},
		8: {"abcdef", 3, "xy", "abcxyf"},
	}
	for i, c := range cases {
		sv := NewString(c.src)
		ln := sv.SetRange(c.ofs, c.set)
		got := sv.Get()

		if got != c.exp {
			t.Errorf("%d: expected %q, got %q", i, c.exp, got)
		}
		if int(ln) != len(c.exp) {
			t.Errorf("%d: expected length of %d, got %d", i, len(c.exp), ln)
		}
	}
}
