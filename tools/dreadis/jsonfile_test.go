package main

import "testing"

func TestLoadJSONFile(t *testing.T) {
	cases := []struct {
		file string
		cnt  int
		err  string
	}{
		{"../../fixtures/cmd-single-get.json", 1, ""},
	}
	for i, c := range cases {
		jf, err := loadJSONFile(c.file)
		if (err == nil) != (c.err == "") {
			t.Errorf("%d: expected err? to be %t, got %v", i, (c.err != ""), err)
			continue
		}
		if err != nil && err.Error() != c.err {
			t.Errorf("%d: expected error %q, got %q", i, c.err, err)
			continue
		}
		if jf.cmds != c.cnt {
			t.Errorf("%d: expected %d commands, got %d", i, c.cnt, jf.cmds)
		}
	}
}
