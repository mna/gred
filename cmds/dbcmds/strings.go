package dbcmd

import (
	"github.com/PuerkitoBio/gred/cmds"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
)

func init() {
	cmds.Register("append", append_)
	cmds.Register("get", get)
	cmds.Register("getrange", getrange)
	cmds.Register("substr", getrange) // alias
	cmds.Register("getset", getset)
	cmds.Register("set", set)
	cmds.Register("strlen", strlen)
}

var append_ = &dbCmd{
	noKey:   srv.NoKeyCreateString,
	minArgs: 2,
	maxArgs: 2,
	fn:      appendFn,
}

func appendFn(db srv.DB, args []string, ints []int, floats []float64) (interface{}, error) {
	k, def := db.LockGetKey(args[0], db.noKey)
	defer def()

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.Append(args[1]), nil
	}
	return nil, cmds.ErrInvalidValType
}

var get = &dbCmd{
	noKey:   srv.NoKeyDefaultVal,
	minArgs: 1,
	maxArgs: 1,
	fn:      getFn,
}

func getFn(k srv.Key, args []string, ints []int, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.Get(), nil
	}
	return nil, cmds.ErrInvalidValType
}

var getrange = &dbCmd{
	noKey:      srv.NoKeyDefaultVal,
	minArgs:    3,
	maxArgs:    3,
	intIndices: []int{1, 2},
	fn:         getrangeFn,
}

func getrangeFn(k srv.Key, args []string, ints []int, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.GetRange(ints[0], ints[1]), nil
	}
	return nil, cmds.ErrInvalidValType
}

var getset = &dbCmd{
	noKey:   srv.NoKeyCreateString,
	minArgs: 2,
	maxArgs: 2,
	fn:      getsetFn,
}

func getsetFn(k srv.Key, args []string, ints []int, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		k.Abort()
		return v.GetSet(args[1]), nil
	}
	return nil, cmds.ErrInvalidValType
}

var set = &dbCmd{
	noKey:   srv.NoKeyCreateString,
	minArgs: 2,
	maxArgs: 2,
	fn:      setFn,
}

func setFn(k srv.Key, args []string, ints []int, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		k.Abort()
		v.Set(args[1])
		return nil, nil
	}
	return nil, cmds.ErrInvalidValType
}

var strlen = &dbCmd{
	noKey:   srv.NoKeyDefaultVal,
	minArgs: 1,
	maxArgs: 1,
	fn:      strlenFn,
}

func strlenFn(k srv.Key, args []string, ints []int, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.StrLen(), nil
	}
	return nil, cmds.ErrInvalidValType
}
