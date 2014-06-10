package strings

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/types"
)

func init() {
	cmd.Register("append", appendƒ)
	cmd.Register("decr", decr)
	cmd.Register("decrby", decrby)
	cmd.Register("get", get)
	cmd.Register("getrange", getrange)
	cmd.Register("getset", getset)
	cmd.Register("incr", incr)
	cmd.Register("incrby", incrby)
	cmd.Register("incrbyfloat", incrbyfloat)
	cmd.Register("set", set)
	cmd.Register("setrange", setrange)
	cmd.Register("strlen", strlen)
}

var appendƒ = cmd.NewSingleKeyCmd(
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
	if v, ok := v.(types.String); ok {
		return v.Append(args[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var decr = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyCreateStringInt,
	decrFn)

func decrFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.IncString); ok {
		val, ok := v.Decr()
		if ok {
			return val, nil
		}
		return nil, cmd.ErrNotInteger
	}
	return nil, cmd.ErrInvalidValType
}

var decrby = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    2,
		IntIndices: []int{1},
	},
	srv.NoKeyCreateStringInt,
	decrbyFn)

func decrbyFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.IncString); ok {
		val, ok := v.DecrBy(ints[0])
		if ok {
			return val, nil
		}
		return nil, cmd.ErrNotInteger
	}
	return nil, cmd.ErrInvalidValType
}

var get = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyNone,
	getFn)

func getFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	if k == nil {
		return nil, nil
	}

	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.String); ok {
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
	if v, ok := v.(types.String); ok {
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
	if v, ok := v.(types.String); ok {
		k.Abort()
		return v.GetSet(args[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var incr = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyCreateStringInt,
	incrFn)

func incrFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.IncString); ok {
		val, ok := v.Incr()
		if ok {
			return val, nil
		}
		return nil, cmd.ErrNotInteger
	}
	return nil, cmd.ErrInvalidValType
}

var incrby = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    2,
		IntIndices: []int{1},
	},
	srv.NoKeyCreateStringInt,
	incrbyFn)

func incrbyFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.IncString); ok {
		val, ok := v.IncrBy(ints[0])
		if ok {
			return val, nil
		}
		return nil, cmd.ErrNotInteger
	}
	return nil, cmd.ErrInvalidValType
}

var incrbyfloat = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:      2,
		MaxArgs:      2,
		FloatIndices: []int{1},
	},
	srv.NoKeyCreateStringInt,
	incrbyfloatFn)

func incrbyfloatFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.IncString); ok {
		val, ok := v.IncrByFloat(floats[0])
		if ok {
			return val, nil
		}
		return nil, cmd.ErrNotFloat
	}
	return nil, cmd.ErrInvalidValType
}

// TODO : This doesn't handle the extra optional args (EX, NX, PX)
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
	if v, ok := v.(types.String); ok {
		k.Abort()
		v.Set(args[1])
		return cmd.OKVal, nil
	}
	return nil, cmd.ErrInvalidValType
}

var setrange = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{1},
		ValidateFn: func(args []string, ints []int64, floats []float64) error {
			if ints[0] < 0 {
				return cmd.ErrOfsOutOfRange
			}
			return nil
		},
	},
	srv.NoKeyDefaultVal,
	setrangeFn)

func setrangeFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.String); ok {
		return v.SetRange(ints[0], args[2]), nil
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
	if v, ok := v.(types.String); ok {
		return v.StrLen(), nil
	}
	return nil, cmd.ErrInvalidValType
}
