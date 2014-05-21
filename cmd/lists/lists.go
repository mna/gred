package lists

import (
	"strings"

	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
)

func init() {
	cmd.Register("lindex", lindex)
	cmd.Register("linsert", linsert)
	cmd.Register("llen", llen)
	cmd.Register("lpop", lpop)
	cmd.Register("lpush", lpush)
	cmd.Register("lpushx", lpushx)
}

var lindex = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    2,
		IntIndices: []int{1},
	},
	srv.NoKeyDefaultVal,
	lindexFn)

func lindexFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		val, ok := v.LIndex(ints[0])
		if ok {
			return val, nil
		}
		return nil, nil
	}
	return nil, cmd.ErrInvalidValType
}

var linsert = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 4,
		MaxArgs: 4,
		ValidateFn: func(args []string, ints []int64, floats []float64) error {
			ba := strings.ToLower(args[1])
			if ba != "before" && ba != "after" {
				return cmd.ErrSyntax
			}
			args[1] = ba
			return nil
		},
	},
	srv.NoKeyDefaultVal,
	linsertFn)

func linsertFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		if args[0] == "before" {
			return v.LInsertBefore(args[2], args[3]), nil
		}
		return v.LInsertAfter(args[2], args[3]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var llen = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	llenFn)

func llenFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		return v.LLen(), nil
	}
	return nil, cmd.ErrInvalidValType
}

var lpop = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	lpopFn)

func lpopFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		val, ok := v.LPop()
		if ok {
			// TODO : Remove key if llen = 0
			return val, nil
		}
		return nil, nil
	}
	return nil, cmd.ErrInvalidValType
}

var lpush = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: -1,
	},
	srv.NoKeyCreateList,
	lpushFn)

func lpushFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		return v.LPush(args[1:]...), nil
	}
	return nil, cmd.ErrInvalidValType
}

var lpushx = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	srv.NoKeyNone,
	lpushxFn)

func lpushxFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	if k == nil {
		return int64(0), nil
	}

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		return v.LPush(args[1:]...), nil
	}
	return nil, cmd.ErrInvalidValType
}
