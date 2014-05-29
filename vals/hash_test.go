package vals

import "testing"

var hcase = hash{
	"a": "v1",
	"b": "v2",
	"c": "v3",
}

func TestHashType(t *testing.T) {
	if got := hcase.Type(); got != "hash" {
		t.Errorf("expected %q, got %q", "hash", got)
	}
}

func TestHashHDel(t *testing.T) {
	// copy the hcase
	hdelCase := make(hash)
	for k, v := range hcase {
		hdelCase[k] = v
	}
	cases := []struct {
		h      Hash
		fields []string
		exp    int64
	}{
		0: {hdelCase, []string{"z", "x", "y"}, 0},
		1: {hdelCase, []string{"z", "a", "y"}, 1},
		2: {hdelCase, []string{"a", "a", "a"}, 0},
		3: {hdelCase, []string{"b", "c", "d"}, 2},
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
	}
	for i, c := range cases {
		got := c.h.HExists(c.field)
		if c.exp != got {
			t.Errorf("%d: expected %t, got %t", i, c.exp, got)
		}
	}
}
