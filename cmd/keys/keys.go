package dbcmds

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
)

func init() {
	cmd.Register("del", del)
	cmd.Register("exists", exists)
	cmd.Register("expire", expire)
	cmd.Register("expireat", expireat)
	cmd.Register("persist", persist)
	cmd.Register("pexpire", pexpire)
	cmd.Register("pexpireat", pexpireat)
	cmd.Register("pttl", pttl)
	cmd.Register("ttl", ttl)
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

func delExpFn(db srv.DB, nm string) {
	db.Lock()
	defer db.Unlock()
	db.Del(nm)
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

	return db.Expire(args[0], ints[0], func() { delExpFn(db, args[0]) }), nil
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

	return db.ExpireAt(args[0], ints[0], func() { delExpFn(db, args[0]) }), nil
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

var pexpire = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    2,
		IntIndices: []int{1},
	},
	pexpireFn)

func pexpireFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	return db.PExpire(args[0], ints[0], func() { delExpFn(db, args[0]) }), nil
}

var pexpireat = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    2,
		IntIndices: []int{1},
	},
	pexpireatFn)

func pexpireatFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	return db.PExpireAt(args[0], ints[0], func() { delExpFn(db, args[0]) }), nil
}

var pttl = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	pttlFn)

func pttlFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	return db.PTTL(args[0]), nil
}

var ttl = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	ttlFn)

func ttlFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	return db.TTL(args[0]), nil
}
