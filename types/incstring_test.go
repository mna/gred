package types

import "testing"

func TestIncStringIncr(t *testing.T) {
	cases := []struct {
		s   string
		exp int64
		ok  bool
	}{
		0: {"", 0, false},
		1: {"0", 1, true},
		2: {"5", 6, true},
		3: {"-6.34", 0, false},
		4: {"-6", -5, true},
		5: {"abc", 0, false},
	}
	for i, c := range cases {
		is := NewIncString(c.s)
		got, ok := is.Incr()
		if c.exp != got {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.ok != ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}

func TestIncStringIncrBy(t *testing.T) {
	cases := []struct {
		s   string
		inc int64
		exp int64
		ok  bool
	}{
		0: {"", 1, 0, false},
		1: {"3", 45, 48, true},
		2: {"3", -45, -42, true},
		3: {"abc", -45, 0, false},
	}
	for i, c := range cases {
		is := NewIncString(c.s)
		got, ok := is.IncrBy(c.inc)
		if c.exp != got {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.ok != ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}

func TestIncStringIncrByFloat(t *testing.T) {
	cases := []struct {
		s   string
		inc float64
		exp string
		ok  bool
	}{
		0: {"", 1.0, "", false},
		1: {"a", 1.3, "", false},
		2: {"3", 1.3, "4.3", true},
		3: {"3.4", 1.3, "4.7", true},
		4: {"-3.45", -2.3, "-5.75", true},
	}
	for i, c := range cases {
		is := NewIncString(c.s)
		got, ok := is.IncrByFloat(c.inc)
		if c.exp != got {
			t.Errorf("%d: expected %s, got %s", i, c.exp, got)
		}
		if c.ok != ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}

func TestIncStringDecr(t *testing.T) {
	cases := []struct {
		s   string
		exp int64
		ok  bool
	}{
		0: {"", 0, false},
		1: {"0", -1, true},
		2: {"5", 4, true},
		3: {"-6.34", 0, false},
		4: {"-6", -7, true},
		5: {"abc", 0, false},
	}
	for i, c := range cases {
		is := NewIncString(c.s)
		got, ok := is.Decr()
		if c.exp != got {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.ok != ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}

func TestIncStringDecrBy(t *testing.T) {
	cases := []struct {
		s   string
		dec int64
		exp int64
		ok  bool
	}{
		0: {"", 1, 0, false},
		1: {"3", 45, -42, true},
		2: {"3", -45, 48, true},
		3: {"abc", 45, 0, false},
	}
	for i, c := range cases {
		is := NewIncString(c.s)
		got, ok := is.DecrBy(c.dec)
		if c.exp != got {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.ok != ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}
