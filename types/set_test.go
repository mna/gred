package types

import (
	"reflect"
	"testing"
)

var setcase Set
var setempty = NewSet()

func init() {
	setcase = NewSet()
	setcase.SAdd("a", "b", "c")
}

func cloneSet(s Set) Set {
	newset := NewSet()
	vals := setcase.SMembers()
	newset.SAdd(vals...)
	return newset
}

func TestSetType(t *testing.T) {
	tp := setcase.Type()
	if tp != "set" {
		t.Errorf("expected %q, got %q", "set", tp)
	}
}

func TestSetSAdd(t *testing.T) {
	empty := NewSet()
	set := cloneSet(setcase)
	cases := []struct {
		s    Set
		vals []string
		exp  int64
	}{
		0: {empty, nil, 0},
		1: {empty, []string{"a"}, 1},
		2: {empty, []string{"a", "b", "c"}, 2},
		3: {set, []string{"a", "b", "c"}, 0},
		4: {set, []string{"d", "d", "d"}, 1},
	}
	for i, c := range cases {
		orilen := c.s.SCard()
		got := c.s.SAdd(c.vals...)
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if ln := c.s.SCard(); ln != orilen+c.exp {
			t.Errorf("%d: expected length to be %d, got %d", i, orilen+c.exp, ln)
		}
	}
}

func TestSetSCard(t *testing.T) {
	cases := []struct {
		s   Set
		exp int64
	}{
		0: {setempty, 0},
		1: {setcase, 3},
	}
	for i, c := range cases {
		got := c.s.SCard()
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
	}
}

func TestSetSIsMember(t *testing.T) {
	cases := []struct {
		s   Set
		v   string
		exp bool
	}{
		0: {setempty, "a", false},
		1: {setcase, "a", true},
		2: {setcase, "b", true},
		3: {setcase, "c", true},
		4: {setcase, "d", false},
	}
	for i, c := range cases {
		got := c.s.SIsMember(c.v)
		if got != c.exp {
			t.Errorf("%d: expected %t, got %t", i, c.exp, got)
		}
	}
}

func TestSetSMembers(t *testing.T) {
	cases := []struct {
		s   Set
		exp []string
	}{
		0: {setempty, []string{}},
		1: {setcase, []string{"a", "b", "c"}},
	}
	for i, c := range cases {
		got := c.s.SMembers()
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
	}
}
