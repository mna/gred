package types

import (
	"reflect"
	"testing"
)

func TestListLIndex(t *testing.T) {
	cases := []struct {
		l   []string
		ix  int64
		exp string
		ok  bool
	}{
		0: {nil, 0, "", false},
		1: {[]string{}, 0, "", false},
		2: {[]string{"a", "b", "c"}, 0, "a", true},
		3: {[]string{"a", "b", "c"}, 1, "b", true},
		4: {[]string{"a", "b", "c"}, 2, "c", true},
		5: {[]string{"a", "b", "c"}, 3, "", false},
		6: {[]string{"a", "b", "c"}, -1, "c", true},
		7: {[]string{"a", "b", "c"}, -2, "b", true},
		8: {[]string{"a", "b", "c"}, -3, "a", true},
		9: {[]string{"a", "b", "c"}, -4, "", false},
	}
	for i, c := range cases {
		l := list(c.l)
		got, ok := l.LIndex(c.ix)
		if got != c.exp {
			t.Errorf("%d: expected %q, got %q", i, c.exp, got)
		}
		if ok != c.ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}

func TestListLLen(t *testing.T) {
	cases := []struct {
		l   []string
		exp int64
	}{
		0: {nil, 0},
		1: {[]string{}, 0},
		2: {[]string{"a"}, 1},
		3: {[]string{"a", "b", "c"}, 3},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LLen()
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
	}
}

func TestListLPush(t *testing.T) {
	cases := []struct {
		l    []string
		vals []string
		exp  []string
	}{
		0: {nil, nil, nil},
		1: {[]string{}, []string{}, []string{}},
		2: {[]string{}, []string{"a"}, []string{"a"}},
		3: {[]string{}, []string{"a", "b", "c"}, []string{"c", "b", "a"}},
		4: {[]string{"a", "b"}, []string{"x", "y", "z"}, []string{"z", "y", "x", "a", "b"}},
		5: {[]string{}, []string{"c", "b", "a"}, []string{"a", "b", "c"}},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LPush(c.vals...)
		if got != int64(len(c.exp)) {
			t.Errorf("%d: expected length of %d, got %d", i, len(c.exp), got)
		}
		if !reflect.DeepEqual([]string(l), c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, l)
		}
	}
}

func TestListLInsertBefore(t *testing.T) {
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
		7: {[]string{"e", "d", "a", "c", "b", "a"}, "a", "z", 7, 2},
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
	}
}

func TestListLInsertAfter(t *testing.T) {
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
	}
}

func TestListLPop(t *testing.T) {
	cases := []struct {
		l   []string
		exp string
		ok  bool
	}{
		0: {nil, "", false},
		1: {[]string{}, "", false},
		2: {[]string{"a"}, "a", true},
		3: {[]string{"a", "b", "c"}, "a", true},
	}
	for i, c := range cases {
		l := list(c.l)
		orilen := len(c.l)
		got, ok := l.LPop()
		if got != c.exp {
			t.Errorf("%d: expected %q, got %q", i, c.exp, got)
		}
		if ok != c.ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
		if c.ok {
			if ln := l.LLen(); ln != int64(orilen-1) {
				t.Errorf("%d: expected length to be %d, got %d", i, orilen-1, ln)
			}
		}
	}
}

func TestListLRange(t *testing.T) {
	cases := []struct {
		l           []string
		start, stop int64
		exp         []string
	}{
		0: {nil, 0, 1, []string{}},
		1: {[]string{}, 0, 2, []string{}},
		2: {[]string{"a"}, 0, 2, []string{"a"}},
		3: {[]string{"a", "b", "c"}, 1, 2, []string{"b", "c"}},
		4: {[]string{"a", "b", "c"}, -3, 2, []string{"a", "b", "c"}},
		5: {[]string{"a", "b", "c"}, 1, 222, []string{"b", "c"}},
		6: {[]string{"a", "b", "c"}, -123, -2, []string{"a", "b"}},
		7: {[]string{"a", "b", "c"}, -123, -5, []string{}},
		8: {[]string{"a", "b", "c"}, 17, -1, []string{}},
		9: {[]string{"a", "b", "c"}, 17, -18, []string{}},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LRange(c.start, c.stop)
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
	}
}

func TestListLRem(t *testing.T) {
	cases := []struct {
		l      []string
		val    string
		cnt, n int64
		exp    []string
	}{
		0:  {nil, "", 0, 0, nil},
		1:  {[]string{}, "", 0, 0, []string{}},
		2:  {[]string{"a", "b", "c"}, "z", 0, 0, []string{"a", "b", "c"}},
		3:  {[]string{"a", "b", "c"}, "z", 2, 0, []string{"a", "b", "c"}},
		4:  {[]string{"a", "b", "c"}, "z", -1, 0, []string{"a", "b", "c"}},
		5:  {[]string{"a", "z", "c", "z"}, "z", 0, 2, []string{"a", "c"}},
		6:  {[]string{"a", "z", "c", "z"}, "z", 1, 1, []string{"a", "c", "z"}},
		7:  {[]string{"a", "z", "c", "z"}, "z", 3, 2, []string{"a", "c"}},
		8:  {[]string{"a", "z", "c", "z"}, "z", -1, 1, []string{"a", "z", "c"}},
		9:  {[]string{"a", "z", "c", "z"}, "z", -4, 2, []string{"a", "c"}},
		10: {[]string{"a", "z", "c", "z"}, "a", -4, 1, []string{"z", "c", "z"}},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LRem(c.cnt, c.val)
		if got != c.n {
			t.Errorf("%d: expected %d elements removed, got %d", i, c.n, got)
		}
		if !reflect.DeepEqual([]string(l), c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, l)
		}
	}
}

func TestListLSet(t *testing.T) {
	cases := []struct {
		l   []string
		val string
		ix  int64
		exp []string
		res bool
	}{
		0:  {nil, "", 0, nil, false},
		1:  {[]string{}, "", 0, []string{}, false},
		2:  {[]string{"a"}, "b", 0, []string{"b"}, true},
		3:  {[]string{"a", "b", "c"}, "z", 0, []string{"z", "b", "c"}, true},
		4:  {[]string{"a", "b", "c"}, "z", 1, []string{"a", "z", "c"}, true},
		5:  {[]string{"a", "b", "c"}, "z", 2, []string{"a", "b", "z"}, true},
		6:  {[]string{"a", "b", "c"}, "z", 3, []string{"a", "b", "c"}, false},
		7:  {[]string{"a", "b", "c"}, "z", -1, []string{"a", "b", "z"}, true},
		8:  {[]string{"a", "b", "c"}, "z", -2, []string{"a", "z", "c"}, true},
		9:  {[]string{"a", "b", "c"}, "z", -3, []string{"z", "b", "c"}, true},
		10: {[]string{"a", "b", "c"}, "z", -4, []string{"a", "b", "c"}, false},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LSet(c.ix, c.val)
		if got != c.res {
			t.Errorf("%d: expected %v, got %v", i, c.res, got)
		}
		if !reflect.DeepEqual([]string(l), c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, l)
		}
	}
}

func TestListLTrim(t *testing.T) {
	cases := []struct {
		l           []string
		start, stop int64
		exp         []string
	}{
		0:  {nil, 0, 0, nil},
		1:  {[]string{}, 0, 0, []string{}},
		2:  {[]string{"a"}, 0, 0, []string{"a"}},
		3:  {[]string{"a", "b", "c"}, 0, 0, []string{"a"}},
		4:  {[]string{"a", "b", "c"}, 0, 1, []string{"a", "b"}},
		5:  {[]string{"a", "b", "c"}, 0, 2, []string{"a", "b", "c"}},
		6:  {[]string{"a", "b", "c"}, 0, 3, []string{"a", "b", "c"}},
		7:  {[]string{"a", "b", "c"}, 2, 3, []string{"c"}},
		8:  {[]string{"a", "b", "c"}, -2, 1, []string{"b"}},
		9:  {[]string{"a", "b", "c"}, -2, 0, []string{}},
		10: {[]string{"a", "b", "c"}, -1, -3, []string{}},
		11: {[]string{"a", "b", "c"}, -5, -3, []string{"a"}},
		12: {[]string{"a", "b", "c"}, -15, -13, []string{}},
	}
	for i, c := range cases {
		l := list(c.l)
		l.LTrim(c.start, c.stop)
		if !reflect.DeepEqual([]string(l), c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, l)
		}
	}
}

func TestListRPush(t *testing.T) {
	cases := []struct {
		l    []string
		vals []string
		exp  []string
	}{
		0: {nil, nil, nil},
		1: {[]string{}, []string{}, []string{}},
		2: {[]string{}, []string{"a"}, []string{"a"}},
		3: {[]string{}, []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		4: {[]string{"a", "b"}, []string{"x", "y", "z"}, []string{"a", "b", "x", "y", "z"}},
		5: {[]string{}, []string{"c", "b", "a"}, []string{"c", "b", "a"}},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.RPush(c.vals...)
		if got != int64(len(c.exp)) {
			t.Errorf("%d: expected length of %d, got %d", i, len(c.exp), got)
		}
		if !reflect.DeepEqual([]string(l), c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, l)
		}
	}
}

func TestListRPop(t *testing.T) {
	cases := []struct {
		l   []string
		exp string
		ok  bool
	}{
		0: {nil, "", false},
		1: {[]string{}, "", false},
		2: {[]string{"a"}, "a", true},
		3: {[]string{"a", "b", "c"}, "c", true},
	}
	for i, c := range cases {
		l := list(c.l)
		orilen := len(c.l)
		got, ok := l.RPop()
		if got != c.exp {
			t.Errorf("%d: expected %q, got %q", i, c.exp, got)
		}
		if ok != c.ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
		if c.ok {
			if ln := l.LLen(); ln != int64(orilen-1) {
				t.Errorf("%d: expected length to be %d, got %d", i, orilen-1, ln)
			}
		}
	}
}

func TestListType(t *testing.T) {
	l := NewList()
	tp := l.Type()
	if tp != "list" {
		t.Errorf("expected %q, got %q", "list", tp)
	}
}
