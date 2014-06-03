package types

import "testing"

func TestStringType(t *testing.T) {
	s := NewString("")
	tp := s.Type()
	if tp != "string" {
		t.Errorf("expected %q, got %q", "string", tp)
	}
}

func TestStringAppend(t *testing.T) {
	cases := []struct {
		s   string
		v   string
		exp int64
	}{
		0: {"", "", 0},
		1: {"", "a", 1},
		2: {"a", "b", 2},
		3: {"a", "", 1},
		4: {"abc", "def", 6},
	}
	for i, c := range cases {
		s := NewString(c.s)
		got := s.Append(c.v)
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if news := c.s + c.v; news != s.Get() {
			t.Errorf("%d: expected %q, got %q", i, news, s.Get())
		}
	}
}

func TestStringGet(t *testing.T) {
	cases := []string{
		0: "",
		1: "a",
		2: "abc",
	}
	for i, c := range cases {
		s := NewString(c)
		got := s.Get()
		if got != c {
			t.Errorf("%d: expected %q, got %q", i, c, got)
		}
	}
}

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
