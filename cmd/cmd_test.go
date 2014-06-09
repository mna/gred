package cmd

import (
	"reflect"
	"testing"
)

func TestArgDefParse(t *testing.T) {
	cases := []struct {
		ad     *ArgDef
		args   []string
		ints   []int64
		floats []float64
		err    bool
	}{
		0: {&ArgDef{
			MinArgs: 0,
			MaxArgs: 0,
		}, []string{"a", "b"}, nil, nil, true},
		1: {&ArgDef{
			MinArgs: 0,
			MaxArgs: 0,
		}, []string{}, []int64{}, []float64{}, false},
		2: {&ArgDef{
			MinArgs: 1,
			MaxArgs: 1,
		}, []string{"a", "b"}, nil, nil, true},
		3: {&ArgDef{
			MinArgs: 1,
			MaxArgs: 1,
		}, []string{"a"}, []int64{}, []float64{}, false},
		4: {&ArgDef{
			MinArgs: 1,
			MaxArgs: 3,
		}, []string{"a"}, []int64{}, []float64{}, false},
		5: {&ArgDef{
			MinArgs: 1,
			MaxArgs: 3,
		}, []string{"a", "b"}, []int64{}, []float64{}, false},
		6: {&ArgDef{
			MinArgs: 1,
			MaxArgs: 3,
		}, []string{"a", "b", "c"}, []int64{}, []float64{}, false},
		7: {&ArgDef{
			MinArgs: 1,
			MaxArgs: 3,
		}, []string{"a", "b", "c", "d"}, []int64{}, []float64{}, true},
		8: {&ArgDef{
			MinArgs: 1,
			MaxArgs: -1,
		}, []string{"a", "b", "c", "d"}, []int64{}, []float64{}, false},
		9: {&ArgDef{
			MinArgs:    1,
			MaxArgs:    2,
			IntIndices: []int{1},
		}, []string{"a", "b"}, []int64{}, []float64{}, true},
		10: {&ArgDef{
			MinArgs:    1,
			MaxArgs:    2,
			IntIndices: []int{1},
		}, []string{"a", "5"}, []int64{5}, []float64{}, false},
		11: {&ArgDef{
			MinArgs:    1,
			MaxArgs:    2,
			IntIndices: []int{0, 1},
		}, []string{"0", "5"}, []int64{0, 5}, []float64{}, false},
		12: {&ArgDef{
			MinArgs:      1,
			MaxArgs:      2,
			FloatIndices: []int{1},
		}, []string{"a", "5.1"}, []int64{}, []float64{5.1}, false},
		13: {&ArgDef{
			MinArgs:      1,
			MaxArgs:      2,
			FloatIndices: []int{0, 1},
		}, []string{"a", "5.1"}, []int64{}, []float64{}, true},
		14: {&ArgDef{
			MinArgs:      1,
			MaxArgs:      2,
			FloatIndices: []int{0, 1},
		}, []string{"-1.2", "5.1"}, []int64{}, []float64{-1.2, 5.1}, false},
		15: {&ArgDef{
			MinArgs:      1,
			MaxArgs:      3,
			FloatIndices: []int{1},
			IntIndices:   []int{2},
		}, []string{"a", "5.1", "2"}, []int64{2}, []float64{5.1}, false},
	}
	for i, c := range cases {
		gots, goti, gotf, err := c.ad.Parse("", c.args)
		if (err != nil) != c.err {
			t.Errorf("%d: error should be %t, got %t", i, c.err, err != nil)
		}
		if c.err {
			t.Logf("%d: got error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(gots, c.args) {
			t.Errorf("%d: expected args to be %v, got %v", i, c.args, gots)
		}
		if !reflect.DeepEqual(goti, c.ints) {
			t.Errorf("%d: expected ints to be %v, got %v", i, c.ints, goti)
		}
		if !reflect.DeepEqual(gotf, c.floats) {
			t.Errorf("%d: expected floats to be %v, got %v", i, c.floats, gotf)
		}
	}
}
