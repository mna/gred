package vals

import "testing"

func TestLInsertBefore(t *testing.T) {
	cases := []struct {
		l        []string
		piv, val string
		exp, at  int64
	}{
		0: {nil, "", "", -1, -1},
		1: {[]string{}, "", "", -1, -1},
		2: {[]string{"a"}, "a", "z", 2, 0},
		3: {[]string{"a"}, "x", "z", -1, -1},
		4: {[]string{"a", "b", "c"}, "a", "z", 4, 0},
		5: {[]string{"a", "b", "c"}, "b", "z", 4, 1},
		6: {[]string{"a", "b", "c"}, "c", "z", 4, 2},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LInsertBefore(c.piv, c.val)
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.at >= 0 {
			if l[c.at] != c.val {
				t.Errorf("%d: value %q should be at index %d, got %q", i, c.val, c.at, l[c.at])
			}
		}
		t.Logf("%d: %v", i, l)
	}
}

func TestLInsertAfter(t *testing.T) {
	cases := []struct {
		l        []string
		piv, val string
		exp, at  int64
	}{
		0: {nil, "", "", -1, -1},
		1: {[]string{}, "", "", -1, -1},
		2: {[]string{"a"}, "a", "z", 2, 1},
		3: {[]string{"a"}, "x", "z", -1, -1},
		4: {[]string{"a", "b", "c"}, "a", "z", 4, 1},
		5: {[]string{"a", "b", "c"}, "b", "z", 4, 2},
		6: {[]string{"a", "b", "c"}, "c", "z", 4, 3},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LInsertAfter(c.piv, c.val)
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.at >= 0 {
			if l[c.at] != c.val {
				t.Errorf("%d: value %q should be at index %d, got %q", i, c.val, c.at, l[c.at])
			}
		}
		t.Logf("%d: %v", i, l)
	}
}
