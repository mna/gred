package strings

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
)

func init() {
	cmd.Register("append", append_)
	cmd.Register("get", get)
	cmd.Register("getrange", getrange)
	cmd.Register("substr", getrange) // alias
	cmd.Register("getset", getset)
	cmd.Register("set", set)
	cmd.Register("strlen", strlen)
}

var append_ = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	srv.NoKeyCreateString,
	appendFn)

func appendFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.Append(args[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var get = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	getFn)

func getFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.Get(), nil
	}
	return nil, cmd.ErrInvalidValType
}

var getrange = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{1, 2},
	},
	srv.NoKeyDefaultVal,
	getrangeFn)

func getrangeFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.GetRange(ints[0], ints[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var getset = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	srv.NoKeyCreateString,
	getsetFn)

func getsetFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		k.Abort()
		return v.GetSet(args[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var set = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	srv.NoKeyCreateString,
	setFn)

func setFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		k.Abort()
		v.Set(args[1])
		return nil, nil
	}
	return nil, cmd.ErrInvalidValType
}

var strlen = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	strlenFn)

func strlenFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.StrLen(), nil
	}
	return nil, cmd.ErrInvalidValType
}
