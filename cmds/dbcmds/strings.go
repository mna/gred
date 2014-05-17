package dbcmd

import (
	"github.com/PuerkitoBio/gred/cmds"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
)

func init() {
	cmds.Register("append", append_)
	cmds.Register("get", get)
	cmds.Register("set", set)
}

var append_ = &dbCmd{
	noKey:   srv.NoKeyCreateString,
	minArgs: 2,
	maxArgs: 2,
	fn:      appendFn,
}

func appendFn(k srv.Key, args []string, ints []int, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.Append(args[1]), nil
	}
	return nil, cmds.ErrInvalidValType
}

var get = &dbCmd{
	noKey:   srv.NoKeyNone,
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
