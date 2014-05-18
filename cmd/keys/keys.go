package dbcmds

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
)

func init() {
	cmd.Register("del", del)
	cmd.Register("exists", exists)
	cmd.Register("expire", expire)
	cmd.Register("persist", persist)
}

var del = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: -1,
	},
	delFn)

func delFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.Lock()
	defer db.Unlock()

	return db.Del(args...), nil
}

var exists = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	existsFn)

func existsFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	return db.Exists(args[0]), nil
}

var expire = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    2,
		IntIndices: []int{1},
	},
	expireFn)

func expireFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	return db.Expire(args[0], ints[0], func() {
		db.Lock()
		defer db.Unlock()
		db.Del(args[0])
	}), nil
}

var expireat = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    2,
		IntIndices: []int{1},
	},
	expireatFn)

func expireatFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	return db.ExpireAt(args[0], ints[0], func() {
		db.Lock()
		defer db.Unlock()
		db.Del(args[0])
	}), nil
}

var persist = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	persistFn)

func persistFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	return db.Persist(args[0]), nil
}
