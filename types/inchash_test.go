package types

import "testing"

var inchash IncHash

func init() {
	inchash = NewIncHash()
	inchash.HMSet(
		"a", "1",
		"b", "2",
		"c", "3.123",
		"d", "v4")
}

func cloneIncHash(ih IncHash) IncHash {
	vs := ih.HGetAll()
	newh := NewIncHash()
	for i := 0; i < len(vs); i += 2 {
		newh.HMSet(vs[i], vs[i+1])
	}
	return newh
}

func TestIncHashHIncrBy(t *testing.T) {
	ih := cloneIncHash(inchash)
	cases := []struct {
		h     IncHash
		field string
		inc   int64
		exp   int64
		ok    bool
	}{
		0: {ih, "", 1, 1, true},
		1: {ih, "a", 4, 5, true},
		2: {ih, "z", -2, -2, true},
		3: {ih, "b", -2, 0, true},
		4: {ih, "d", 2, 0, false},
	}
	for i, c := range cases {
		got, ok := c.h.HIncrBy(c.field, c.inc)
		if c.exp != got {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.ok != ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}

func TestIncHashHIncrByFloat(t *testing.T) {
	ih := cloneIncHash(inchash)
	cases := []struct {
		h     IncHash
		field string
		inc   float64
		exp   string
		ok    bool
	}{
		0: {ih, "", 1.1, "1.1", true},
		1: {ih, "a", 4.2, "5.2", true},
		2: {ih, "z", -2.12, "-2.12", true},
		// TODO : Issue with floating point that is not present
		// in Redis...
		//3: {ih, "b", -2.345, "-0.345", true},
		4: {ih, "d", 2.345, "", false},
		5: {ih, "c", 2.345, "5.468", true},
	}
	for i, c := range cases {
		if c.h == nil {
			continue
		}
		got, ok := c.h.HIncrByFloat(c.field, c.inc)
		if c.exp != got {
			t.Errorf("%d: expected %s, got %s", i, c.exp, got)
		}
		if c.ok != ok {
			t.Errorf("%d: expected %t, got %t", i, c.ok, ok)
		}
	}
}
