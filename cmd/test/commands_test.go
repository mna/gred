package test

import (
	"reflect"
	"testing"

	"github.com/PuerkitoBio/gred/cmd"
	_ "github.com/PuerkitoBio/gred/cmd/hashes"
	_ "github.com/PuerkitoBio/gred/cmd/lists"
	_ "github.com/PuerkitoBio/gred/cmd/sets"
	_ "github.com/PuerkitoBio/gred/cmd/strings"
	"github.com/PuerkitoBio/gred/srv"
)

func TestCommand(t *testing.T) {
	cases := []struct {
		name string
		args []string
		res  interface{}
		err  error
	}{
		// First create a key for all types
		{"set", []string{"s", "val"}, cmd.OKVal, nil},
		{"hset", []string{"h", "f1", "v1"}, true, nil},
		{"lpush", []string{"l", "v1"}, int64(1), nil},
		{"sadd", []string{"t", "v1"}, int64(1), nil},

		// Then test all commands
		{"append", []string{"k", "a"}, int64(1), nil},
		{"append", []string{"k", "bcd"}, int64(4), nil},
		{"append", []string{"t", "bcd"}, nil, cmd.ErrInvalidValType},
		{"get", []string{"k"}, "abcd", nil},
		{"get", []string{"z"}, nil, nil},
		{"get", []string{"l"}, nil, cmd.ErrInvalidValType},
		{"getrange", []string{"k", "2", "4"}, "cd", nil},
		{"getrange", []string{"k", "-1", "4"}, "d", nil},
		{"getrange", []string{"k", "-10", "-40"}, "a", nil},
		{"getrange", []string{"k", "10", "40"}, "", nil},
		{"getrange", []string{"l", "10", "40"}, nil, cmd.ErrInvalidValType},
		{"getset", []string{"k", "efg"}, "abcd", nil},
		{"getset", []string{"z", "efg"}, "", nil},
		{"getset", []string{"h", "efg"}, nil, cmd.ErrInvalidValType},
	}

	var got interface{}
	var gotErr error
	var conn mockConn
	for i, c := range cases {
		// Get the command
		cd := cmd.Commands[c.name]

		// Parse the arguments
		args, ints, floats, err := cd.Parse(c.name, c.args)
		if err != nil {
			t.Fatal(err)
		}
		switch cd := cd.(type) {
		case cmd.DBCmd:
			db, _ := srv.DefaultServer.GetDB(0)
			got, gotErr = cd.ExecWithDB(db, args, ints, floats)
		case cmd.SrvCmd:
			got, gotErr = cd.Exec(args, ints, floats)
		case cmd.ConnCmd:
			got, gotErr = cd.ExecWithConn(&conn, args, ints, floats)
		}

		// Assert the results
		if !reflect.DeepEqual(got, c.res) {
			t.Errorf("%d: expected %v, got %v", i, c.res, got)
		}
		if c.err != gotErr {
			t.Errorf("%d: expected error %v, got %v", i, c.err, gotErr)
		}
	}
}

type mockConn struct {
	ix int
}

func (mc *mockConn) Select(ix int) {
	mc.ix = ix
}
