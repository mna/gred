package test

import (
	"reflect"
	"testing"

	"github.com/PuerkitoBio/gred/cmd"
	_ "github.com/PuerkitoBio/gred/cmd/hashes"
	_ "github.com/PuerkitoBio/gred/cmd/keys"
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

		// Strings
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
		{"del", []string{"z"}, int64(1), nil},
		{"getset", []string{"h", "efg"}, nil, cmd.ErrInvalidValType},
		{"set", []string{"h", "efg"}, nil, cmd.ErrInvalidValType},
		{"setrange", []string{"k", "1", "zzzz"}, int64(5), nil},
		{"setrange", []string{"k", "10", "aa"}, int64(12), nil},
		{"setrange", []string{"t", "10", "aa"}, nil, cmd.ErrInvalidValType},
		{"strlen", []string{"k"}, int64(12), nil},
		{"strlen", []string{"z"}, int64(0), nil},
		{"strlen", []string{"l"}, nil, cmd.ErrInvalidValType},
		{"del", []string{"k"}, int64(1), nil},
		{"incr", []string{"k"}, int64(1), nil},
		{"incr", []string{"k"}, int64(2), nil},
		{"set", []string{"z", "not an int"}, cmd.OKVal, nil},
		{"incr", []string{"z"}, nil, cmd.ErrNotInteger},
		{"incr", []string{"t"}, nil, cmd.ErrInvalidValType},
		{"incrby", []string{"k", "3"}, int64(5), nil},
		{"incrby", []string{"k", "-12"}, int64(-7), nil},
		{"del", []string{"k"}, int64(1), nil},
		{"incrby", []string{"k", "-12"}, int64(-12), nil},
		{"incrby", []string{"l", "-12"}, nil, cmd.ErrInvalidValType},
		{"incrbyfloat", []string{"k", "3.1"}, "-8.9", nil},
		{"incrbyfloat", []string{"k", "-0.2"}, "-9.1", nil},
		{"del", []string{"k"}, int64(1), nil},
		{"incrbyfloat", []string{"k", "0.2"}, "0.2", nil},
		{"incrbyfloat", []string{"l", "0.2"}, nil, cmd.ErrInvalidValType},

		// Hashes
		{"del", []string{"k"}, int64(1), nil},
		{"del", []string{"z"}, int64(1), nil},
		{"hset", []string{"k", "f1", "v1"}, true, nil},
		{"hset", []string{"k", "f2", "v2"}, true, nil},
		{"hdel", []string{"z", "f1"}, int64(0), nil},
		{"hdel", []string{"k", "f1", "f1", "f2", "f3"}, int64(2), nil},
		{"hdel", []string{"s", "f1"}, nil, cmd.ErrInvalidValType},
		{"hset", []string{"k", "f1", "v1"}, true, nil},
		{"hset", []string{"k", "f2", "v2"}, true, nil},
		{"hexists", []string{"z", "f1"}, false, nil},
		{"hexists", []string{"k", "f1"}, true, nil},
		{"hexists", []string{"k", "f3"}, false, nil},
		{"hexists", []string{"l", "f3"}, nil, cmd.ErrInvalidValType},
		{"hget", []string{"z", "f1"}, nil, nil},
		{"hget", []string{"k", "f1"}, "v1", nil},
		{"hget", []string{"k", "f2"}, "v2", nil},
		{"hget", []string{"k", "f3"}, nil, nil},
		{"hget", []string{"t", "f3"}, nil, cmd.ErrInvalidValType},
		{"hgetall", []string{"z"}, []string{}, nil},
		{"hgetall", []string{"k"}, []string{"f1", "v1", "f2", "v2"}, nil},
		{"hgetall", []string{"s"}, nil, cmd.ErrInvalidValType},
		{"hkeys", []string{"z"}, []string{}, nil},
		{"hkeys", []string{"k"}, []string{"f1", "f2"}, nil},
		{"hkeys", []string{"l"}, nil, cmd.ErrInvalidValType},
		{"hlen", []string{"z"}, int64(0), nil},
		{"hlen", []string{"k"}, int64(2), nil},
		{"hlen", []string{"l"}, nil, cmd.ErrInvalidValType},
		{"hmget", []string{"z", "f1", "f2", "f3"}, []interface{}{nil, nil, nil}, nil},
		{"hmget", []string{"k", "f1", "f2", "f3"}, []interface{}{"v1", "v2", nil}, nil},
		{"hmget", []string{"s", "f1", "f2", "f3"}, nil, cmd.ErrInvalidValType},
		{"hmset", []string{"z", "f1", "v1", "f2", "v2"}, cmd.OKVal, nil},
		{"hmset", []string{"k", "f1", "x1", "f3", "v3"}, cmd.OKVal, nil},
		{"hgetall", []string{"k"}, []string{"f1", "x1", "f2", "v2", "f3", "v3"}, nil},
		{"hmset", []string{"t", "f1", "x1"}, nil, cmd.ErrInvalidValType},
		{"del", []string{"z"}, int64(1), nil},
		{"hset", []string{"z", "f1", "v1"}, true, nil},
		{"hset", []string{"k", "f1", "v1"}, false, nil},
		{"hget", []string{"k", "f1"}, "v1", nil},
		{"hset", []string{"k", "f4", "v4"}, true, nil},
		{"hset", []string{"l", "f4", "v4"}, nil, cmd.ErrInvalidValType},
		{"del", []string{"z"}, int64(1), nil},
		{"hsetnx", []string{"z", "f1", "v1"}, true, nil},
		{"hsetnx", []string{"k", "f1", "x1"}, false, nil},
		{"hget", []string{"k", "f1"}, "v1", nil},
		{"hsetnx", []string{"k", "f5", "v5"}, true, nil},
		{"hsetnx", []string{"l", "f4", "v4"}, nil, cmd.ErrInvalidValType},
		{"del", []string{"z"}, int64(1), nil},
		{"hvals", []string{"z"}, []string{}, nil},
		{"hvals", []string{"k"}, []string{"v1", "v2", "v3", "v4", "v5"}, nil},
		{"hvals", []string{"s"}, nil, cmd.ErrInvalidValType},
		{"hincrby", []string{"z", "i1", "3"}, int64(3), nil},
		{"hincrby", []string{"k", "i1", "3"}, int64(3), nil},
		{"hincrby", []string{"k", "i1", "-7"}, int64(-4), nil},
		{"hincrby", []string{"k", "f1", "7"}, nil, cmd.ErrHashFieldNotInt},
		{"hincrby", []string{"l", "f1", "7"}, nil, cmd.ErrInvalidValType},
		{"del", []string{"z"}, int64(1), nil},
		{"hincrbyfloat", []string{"z", "ff", "3.1"}, "3.1", nil},
		{"hincrbyfloat", []string{"k", "ff", "1.2"}, "1.2", nil},
		{"hincrbyfloat", []string{"k", "ff", "-7.4"}, "-6.2", nil},
		{"hincrbyfloat", []string{"k", "f1", "1.4"}, nil, cmd.ErrHashFieldNotFloat},
		{"hincrbyfloat", []string{"l", "f1", "7"}, nil, cmd.ErrInvalidValType},
		{"hdel", []string{"k", "f1", "f2", "f3", "f4", "f5", "i1", "ff"}, int64(7), nil},
		{"exists", []string{"k"}, false, nil},

		// Lists
		{"del", []string{"k"}, int64(0), nil}, // because deleted by previous hdel
		{"del", []string{"z"}, int64(1), nil},
		{"lpush", []string{"k", "a", "b", "c", "a"}, int64(4), nil}, // a c b a
		{"lrange", []string{"k", "0", "-1"}, []string{"a", "c", "b", "a"}, nil},
		{"lpush", []string{"k", "d", "e"}, int64(6), nil}, // e d a c b a
		{"lrange", []string{"k", "0", "-1"}, []string{"e", "d", "a", "c", "b", "a"}, nil},
		{"lpush", []string{"h", "e"}, nil, cmd.ErrInvalidValType},
		{"lindex", []string{"z", "1"}, nil, nil},
		{"lindex", []string{"k", "0"}, "e", nil},
		{"lindex", []string{"k", "-3"}, "c", nil},
		{"lindex", []string{"k", "-30"}, nil, nil},
		{"lindex", []string{"s", "0"}, nil, cmd.ErrInvalidValType},
		{"linsert", []string{"z", "before", "a", "z"}, int64(0), nil},
		{"linsert", []string{"k", "before", "a", "z"}, int64(7), nil}, // e d z a c b a
		{"lrange", []string{"k", "0", "-1"}, []string{"e", "d", "z", "a", "c", "b", "a"}, nil},
		{"linsert", []string{"k", "before", "x", "z"}, int64(-1), nil},
		{"linsert", []string{"z", "after", "a", "y"}, int64(0), nil},
		{"linsert", []string{"k", "after", "a", "y"}, int64(8), nil}, // e d z a y c b a
		{"lrange", []string{"k", "0", "-1"}, []string{"e", "d", "z", "a", "y", "c", "b", "a"}, nil},
		{"linsert", []string{"k", "before", "x", "z"}, int64(-1), nil},
		{"linsert", []string{"k", "after", "x", "y"}, int64(-1), nil},
		{"linsert", []string{"t", "after", "x", "y"}, nil, cmd.ErrInvalidValType},
		{"llen", []string{"z"}, int64(0), nil},
		{"llen", []string{"k"}, int64(8), nil},
		{"llen", []string{"h"}, nil, cmd.ErrInvalidValType},
		{"lpop", []string{"z"}, nil, nil},
		{"lpop", []string{"k"}, "e", nil}, // d z a y c b a
		{"lrange", []string{"k", "0", "-1"}, []string{"d", "z", "a", "y", "c", "b", "a"}, nil},
		{"linsert", []string{"k", "before", "x", "z"}, int64(-1), nil},
		{"lpop", []string{"k"}, "d", nil}, // z a y c b a
		{"lrange", []string{"k", "0", "-1"}, []string{"z", "a", "y", "c", "b", "a"}, nil},
		{"lpop", []string{"s"}, nil, cmd.ErrInvalidValType},
		{"lpush", []string{"j", "a"}, int64(1), nil},
		{"lpop", []string{"j"}, "a", nil},
		{"exists", []string{"j"}, false, nil},
		{"lpushx", []string{"z", "a"}, int64(0), nil},
		{"lpushx", []string{"k", "f"}, int64(7), nil}, // f z a y c b a
		{"lrange", []string{"k", "0", "-1"}, []string{"f", "z", "a", "y", "c", "b", "a"}, nil},
		{"lpushx", []string{"t", "f"}, nil, cmd.ErrInvalidValType},
		{"lrange", []string{"z", "1", "5"}, []string{}, nil},
		{"lrange", []string{"k", "2", "5"}, []string{"a", "y", "c", "b"}, nil},
		{"lrange", []string{"k", "6", "9"}, []string{"a"}, nil},
		{"lrange", []string{"k", "-4", "5"}, []string{"y", "c", "b"}, nil},
		{"lrange", []string{"k", "-4", "-10"}, []string{}, nil},
		{"lrange", []string{"k", "-10", "-5"}, []string{"f", "z", "a"}, nil},
		{"lrange", []string{"h", "0", "5"}, nil, cmd.ErrInvalidValType},
		{"lpush", []string{"k", "a", "b"}, int64(9), nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"b", "a", "f", "z", "a", "y", "c", "b", "a"}, nil},
		{"lrem", []string{"z", "3", "a"}, int64(0), nil},
		{"lrem", []string{"k", "3", "f"}, int64(1), nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"b", "a", "z", "a", "y", "c", "b", "a"}, nil},
		{"lrem", []string{"k", "-1", "b"}, int64(1), nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"b", "a", "z", "a", "y", "c", "a"}, nil},
		{"lrem", []string{"k", "2", "a"}, int64(2), nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"b", "z", "y", "c", "a"}, nil},
		{"lrem", []string{"s", "2", "a"}, nil, cmd.ErrInvalidValType},
		{"lpush", []string{"j", "a"}, int64(1), nil},
		{"lrem", []string{"j", "1", "a"}, int64(1), nil},
		{"exists", []string{"j"}, false, nil},
		{"lset", []string{"z", "2", "a"}, nil, cmd.ErrNoSuchKey},
		{"lset", []string{"k", "2", "a"}, cmd.OKVal, nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"b", "z", "a", "c", "a"}, nil},
		{"lset", []string{"k", "-1", "z"}, cmd.OKVal, nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"b", "z", "a", "c", "z"}, nil},
		{"lset", []string{"k", "-9", "z"}, nil, cmd.ErrOutOfRange},
		{"lset", []string{"k", "9", "z"}, nil, cmd.ErrOutOfRange},
		{"lset", []string{"t", "9", "z"}, nil, cmd.ErrInvalidValType},
		{"ltrim", []string{"z", "1", "3"}, cmd.OKVal, nil},
		{"llen", []string{"z"}, int64(0), nil},
		{"ltrim", []string{"k", "0", "3"}, cmd.OKVal, nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"b", "z", "a", "c"}, nil},
		{"ltrim", []string{"k", "-3", "-1"}, cmd.OKVal, nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"z", "a", "c"}, nil},
		{"ltrim", []string{"k", "1", "1"}, cmd.OKVal, nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"a"}, nil},
		{"ltrim", []string{"s", "1", "1"}, nil, cmd.ErrInvalidValType},
		{"ltrim", []string{"k", "2", "4"}, cmd.OKVal, nil},
		{"exists", []string{"k"}, false, nil},
		{"lpush", []string{"k", "a", "b", "c", "a", "d"}, int64(5), nil},
		{"rpop", []string{"k"}, "a", nil},
		{"rpop", []string{"k"}, "b", nil},
		{"rpop", []string{"z"}, nil, nil},
		{"rpop", []string{"t"}, nil, cmd.ErrInvalidValType},
		{"rpoplpush", []string{"z", "x"}, nil, nil},
		{"exists", []string{"x"}, false, nil},
		{"rpoplpush", []string{"k", "k"}, "c", nil}, // k is now c d a
		{"lrange", []string{"k", "0", "-1"}, []string{"c", "d", "a"}, nil},
		{"rpoplpush", []string{"k", "z"}, "a", nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"c", "d"}, nil},
		{"lrange", []string{"z", "0", "-1"}, []string{"a"}, nil},
		{"rpoplpush", []string{"s", "z"}, nil, cmd.ErrInvalidValType},
		{"del", []string{"z"}, int64(1), nil},
		{"rpush", []string{"z", "a", "b", "c"}, int64(3), nil},
		{"lrange", []string{"z", "0", "-1"}, []string{"a", "b", "c"}, nil},
		{"rpush", []string{"k", "a", "b", "c"}, int64(5), nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"c", "d", "a", "b", "c"}, nil},
		{"rpush", []string{"t", "a", "b", "c"}, nil, cmd.ErrInvalidValType},
		{"del", []string{"z"}, int64(1), nil},
		{"rpushx", []string{"z", "a"}, int64(0), nil},
		{"exists", []string{"z"}, false, nil},
		{"rpushx", []string{"k", "a"}, int64(6), nil},
		{"lrange", []string{"k", "0", "-1"}, []string{"c", "d", "a", "b", "c", "a"}, nil},
		{"rpushx", []string{"t", "a"}, nil, cmd.ErrInvalidValType},
		{"del", []string{"z"}, int64(1), nil},
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
			t.Errorf("%d [%s %v]: expected %v, got %v", i, c.name, c.args, c.res, got)
		}
		if c.err != gotErr {
			t.Errorf("%d [%s %v]: expected error %v, got %v", i, c.name, c.args, c.err, gotErr)
		}
	}
}

type mockConn struct {
	ix int
}

func (mc *mockConn) Select(ix int) {
	mc.ix = ix
}
