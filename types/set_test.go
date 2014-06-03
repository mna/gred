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
	vals := s.SMembers()
	return setFromStrings(vals)
}

func setFromStrings(s []string) Set {
	newset := NewSet()
	newset.SAdd(s...)
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

func TestSetSDiff(t *testing.T) {
	cases := []struct {
		s     Set
		diffs [][]string
		exp   []string
	}{
		0: {setempty, [][]string{}, []string{}},
		1: {setempty, [][]string{
			{"a", "b", "c"},
			{"c", "d"},
		}, []string{}},
		2: {setcase, [][]string{
			{"a", "b", "c"},
			{"c", "d"},
		}, []string{}},
		3: {setcase, [][]string{
			{"a"},
			{"d"},
		}, []string{"b", "c"}},
		4: {setcase, [][]string{
			{"e"},
			{"d"},
		}, []string{"a", "b", "c"}},
	}
	for i, c := range cases {
		sets := make([]Set, len(c.diffs))
		for j, vals := range c.diffs {
			sets[j] = setFromStrings(vals)
		}
		got := c.s.SDiff(sets...)
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
	}
}

func TestSetSInter(t *testing.T) {
	cases := []struct {
		s      Set
		inters [][]string
		exp    []string
	}{
		0: {setempty, [][]string{}, []string{}},
		1: {setempty, [][]string{
			{"a", "b", "c"},
			{"c", "d"},
		}, []string{}},
		2: {setcase, [][]string{
			{"a", "b", "c"},
			{"c", "b"},
		}, []string{"b", "c"}},
		3: {setcase, [][]string{
			{"a"},
			{"d"},
		}, []string{}},
		4: {setcase, [][]string{
			{"e"},
			{"d"},
		}, []string{}},
		5: {setcase, [][]string{
			{"a", "b"},
			{"b", "c"},
		}, []string{"b"}},
	}
	for i, c := range cases {
		sets := make([]Set, len(c.inters))
		for j, vals := range c.inters {
			sets[j] = setFromStrings(vals)
		}
		got := c.s.SInter(sets...)
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
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

func TestSetSRem(t *testing.T) {
	empty := cloneSet(setempty)
	set := cloneSet(setcase)
	cases := []struct {
		s    Set
		vals []string
		exp  int64
		res  []string
	}{
		0: {empty, []string{}, 0, []string{}},
		1: {empty, []string{"a"}, 0, []string{}},
		2: {set, []string{"a", "d", "a", "e"}, 1, []string{"b", "c"}},
		3: {set, []string{"b", "c"}, 2, []string{}},
	}
	for i, c := range cases {
		got := c.s.SRem(c.vals...)
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		vals := c.s.SMembers()
		if !reflect.DeepEqual(vals, c.res) {
			t.Errorf("%d: expected %v, got %v", i, c.res, vals)
		}
	}
}

func TestSetSUnion(t *testing.T) {
	cases := []struct {
		s      Set
		unions [][]string
		exp    []string
	}{
		0: {setempty, [][]string{}, []string{}},
		1: {setempty, [][]string{
			{"a", "b", "c"},
			{"c", "d"},
		}, []string{"a", "b", "c", "d"}},
		2: {setcase, [][]string{
			{"a", "b", "c"},
			{"c", "b"},
		}, []string{"a", "b", "c"}},
		3: {setcase, [][]string{
			{"a"},
			{"d"},
		}, []string{"a", "b", "c", "d"}},
		4: {setcase, [][]string{
			{"e"},
			{"d"},
		}, []string{"a", "b", "c", "e", "d"}},
		5: {setcase, [][]string{
			{"a", "b"},
			{"b", "c"},
		}, []string{"a", "b", "c"}},
	}
	for i, c := range cases {
		sets := make([]Set, len(c.unions))
		for j, vals := range c.unions {
			sets[j] = setFromStrings(vals)
		}
		got := c.s.SUnion(sets...)
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
	}
}
