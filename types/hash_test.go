package types

import (
	"reflect"
	"testing"
)

var hcase = hash{
	"a": "v1",
	"b": "v2",
	"c": "v3",
}

var hempty = hash{}

func cloneHash(h hash) hash {
	h2 := hash{}
	for k, v := range h {
		h2[k] = v
	}
	return h2
}

func TestHashType(t *testing.T) {
	if got := hcase.Type(); got != "hash" {
		t.Errorf("expected %q, got %q", "hash", got)
	}
}

func TestHashHDel(t *testing.T) {
	// copy the hcase
	hdelCase := cloneHash(hcase)
	cases := []struct {
		h      Hash
		fields []string
		exp    int64
	}{
		0: {hdelCase, []string{"z", "x", "y"}, 0},
		1: {hdelCase, []string{"z", "a", "y"}, 1},
		2: {hdelCase, []string{"a", "a", "a"}, 0},
		3: {hdelCase, []string{"b", "c", "d"}, 2},
		4: {hempty, []string{"b", "c", "d"}, 0},
	}
	for i, c := range cases {
		got := c.h.HDel(c.fields...)
		if c.exp != got {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
	}
}

func TestHashExists(t *testing.T) {
	cases := []struct {
		h     Hash
		field string
		exp   bool
	}{
		0: {hcase, "y", false},
		1: {hcase, "a", true},
		2: {hcase, "b", true},
		3: {hcase, "c", true},
		4: {hcase, "d", false},
		5: {hempty, "a", false},
	}
	for i, c := range cases {
		got := c.h.HExists(c.field)
		if c.exp != got {
			t.Errorf("%d: expected %t, got %t", i, c.exp, got)
		}
	}
}

func TestHashGet(t *testing.T) {
	cases := []struct {
		h     Hash
		field string
		exp   string
		ok    bool
	}{
		0: {hcase, "a", "v1", true},
		1: {hcase, "z", "", false},
		2: {hcase, "b", "v2", true},
		3: {hcase, "c", "v3", true},
		4: {hempty, "a", "", false},
	}
	for i, c := range cases {
		got, ok := c.h.HGet(c.field)
		if got != c.exp {
			t.Errorf("%d: expected %q, got %q", i, c.exp, got)
		}
		if ok != c.ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}

func TestHashGetAll(t *testing.T) {
	cases := []struct {
		h   Hash
		exp []string
	}{
		0: {hcase, []string{"a", "v1", "b", "v2", "c", "v3"}},
		1: {hempty, []string{}},
	}
	for i, c := range cases {
		got := c.h.HGetAll()
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
	}
}

func TestHashHKeys(t *testing.T) {
	cases := []struct {
		h   Hash
		exp []string
	}{
		0: {hcase, []string{"a", "b", "c"}},
		1: {hempty, []string{}},
	}
	for i, c := range cases {
		got := c.h.HKeys()
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
	}
}

func TestHashHLen(t *testing.T) {
	cases := []struct {
		h   Hash
		exp int64
	}{
		0: {hcase, 3},
		1: {hempty, 0},
	}
	for i, c := range cases {
		got := c.h.HLen()
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
	}
}

func TestHashHMGet(t *testing.T) {
	cases := []struct {
		h      Hash
		fields []string
		exp    []interface{}
	}{
		0: {hcase, []string{}, []interface{}{}},
		1: {hcase, []string{"a", "a", "a"}, []interface{}{"v1", "v1", "v1"}},
		2: {hcase, []string{"a", "b", "c"}, []interface{}{"v1", "v2", "v3"}},
		3: {hcase, []string{"a", "z", "c", "e"}, []interface{}{"v1", nil, "v3", nil}},
		4: {hempty, []string{"a", "b", "c"}, []interface{}{nil, nil, nil}},
	}
	for i, c := range cases {
		got := c.h.HMGet(c.fields...)
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
	}
}

func TestHashHMSet(t *testing.T) {
	hmset := cloneHash(hcase)
	hmempty := cloneHash(hempty)
	cases := []struct {
		h      Hash
		fields []string
		ln     int64
	}{
		0: {hmset, []string{"a", "v4", "d", "v5"}, 4},
		1: {hmset, []string{"e", "v6"}, 5},
		2: {hmempty, []string{"a", "v1", "b", "v2", "c", "v3", "d", "v4", "a", "v5"}, 4},
	}
	for i, c := range cases {
		c.h.HMSet(c.fields...)
		if ln := c.h.HLen(); ln != c.ln {
			t.Errorf("%d: expected length of %d, got %d", i, c.ln, ln)
		}
	}
}

func TestHashHSet(t *testing.T) {
	hset := cloneHash(hcase)
	hempty := cloneHash(hempty)
	cases := []struct {
		h     Hash
		field string
		val   string
		exp   bool
		ln    int64
	}{
		0: {hset, "a", "v0", false, 3},
		1: {hset, "d", "v4", true, 4},
		2: {hempty, "a", "v1", true, 1},
	}
	for i, c := range cases {
		got := c.h.HSet(c.field, c.val)
		if got != c.exp {
			t.Errorf("%d: expected %t, got %t", i, c.exp, got)
		}
		if ln := c.h.HLen(); ln != c.ln {
			t.Errorf("%d: expected length of %d, got %d", i, c.ln, ln)
		}
		if val, _ := c.h.HGet(c.field); val != c.val {
			t.Errorf("%d: expected field to be %q, got %q", i, c.val, val)
		}
	}
}

func TestHashHSetNx(t *testing.T) {
	hset := cloneHash(hcase)
	hempty := cloneHash(hempty)
	cases := []struct {
		h     Hash
		field string
		val   string
		exp   bool
		ln    int64
	}{
		0: {hset, "a", "v0", false, 3},
		1: {hset, "d", "v4", true, 4},
		2: {hempty, "a", "v1", true, 1},
	}
	for i, c := range cases {
		ori, _ := c.h.HGet(c.field)
		got := c.h.HSetNx(c.field, c.val)
		if got != c.exp {
			t.Errorf("%d: expected %t, got %t", i, c.exp, got)
		}
		if ln := c.h.HLen(); ln != c.ln {
			t.Errorf("%d: expected length of %d, got %d", i, c.ln, ln)
		}
		new, _ := c.h.HGet(c.field)
		if c.exp && c.val != new {
			t.Errorf("%d: expected field to be %q, got %q", i, c.val, new)
		}
		if !c.exp && new != ori {
			t.Errorf("%d: expected field to be %q, got %q", i, ori, new)
		}
	}
}

func TestHashHVals(t *testing.T) {
	cases := []struct {
		h   Hash
		exp []string
	}{
		0: {hcase, []string{"v1", "v2", "v3"}},
		1: {hempty, []string{}},
	}
	for i, c := range cases {
		got := c.h.HVals()
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
	}
}
