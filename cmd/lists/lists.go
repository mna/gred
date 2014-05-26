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
	cmd.Register("lrange", lrange)
	cmd.Register("lrem", lrem)
	cmd.Register("lset", lset)
	cmd.Register("ltrim", ltrim)
	cmd.Register("rpop", rpop)
	cmd.Register("rpoplpush", rpoplpush)
	cmd.Register("rpush", rpush)
	cmd.Register("rpushx", rpushx)
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

var lpop = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	lpopFn)

func lpopFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()

	// Get the key
	k, ok := db.Keys()[args[0]]
	// If the key does not exist, return nil
	if !ok {
		db.RUnlock()
		return nil, nil
	}

	// Pop the value
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		val, ok := v.LPop()
		if ok {
			// If the list is now empty, delete the key
			if v.LLen() == 0 {
				// Upgrade the db lock. Since the key lock is maintained and exclusive,
				// it cannot change during the db key upgrade.
				db.RUnlock()
				db.Lock()
				// Get the keys again
				keys := db.Keys()
				k.Abort()
				delete(keys, k.Name())
				db.Unlock()
				return val, nil
			}
			db.RUnlock()
			return val, nil
		}
		db.RUnlock()
		return nil, nil
	}
	db.RUnlock()
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

var lrange = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{1, 2},
	},
	srv.NoKeyDefaultVal,
	lrangeFn)

func lrangeFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		return v.LRange(ints[0], ints[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var lrem = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{1},
	},
	srv.NoKeyDefaultVal,
	lremFn)

func lremFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		// TODO : Delete key if no more elements
		return v.LRem(ints[0], args[2]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var lset = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{1},
	},
	srv.NoKeyNone,
	lsetFn)

func lsetFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	if k == nil {
		return nil, cmd.ErrNoSuchKey
	}

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		ok := v.LSet(ints[0], args[2])
		if ok {
			return cmd.OKVal, nil
		}
		return nil, cmd.ErrOutOfRange
	}
	return nil, cmd.ErrInvalidValType
}

var ltrim = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{1, 2},
	},
	srv.NoKeyDefaultVal,
	ltrimFn)

func ltrimFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		v.LTrim(ints[0], ints[1])
		return cmd.OKVal, nil
	}
	return nil, cmd.ErrInvalidValType
}

var rpop = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	rpopFn)

func rpopFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		val, ok := v.RPop()
		if ok {
			return val, nil
		}
		return nil, nil
	}
	return nil, cmd.ErrInvalidValType
}

var rpoplpush = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	rpoplpushFn)

func rpoplpushFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()

	// Get the source key
	keys := db.Keys()
	src, ok := keys[args[0]]
	if !ok {
		// Source key does not exist, return nil
		db.RUnlock()
		return nil, nil
	}

	// If source exists, and source and destination are the same
	if args[0] == args[1] {
		defer db.RUnlock()
		src.Lock()
		defer src.Unlock()

		// Simply rotate the value (pop at tail, push at head)
		v := src.Val()
		if v, ok := v.(vals.List); ok {
			val, ok := v.RPop()
			if ok {
				v.LPush(val)
				return val, nil
			}
			return nil, nil
		}
		return nil, cmd.ErrInvalidValType
	}

	// Otherwise get the destination key, and create it if it doesn't exist
	dst, ok := keys[args[1]]
	def := db.RUnlock
	if !ok {
		// Destination does not exist, upgrade db lock
		db.RUnlock()
		db.Lock()
		def = db.Unlock

		// Re-read the destination, as it may have changed during db lock upgrade
		dst, ok = keys[args[1]]
		if !ok {
			// Create the destination key
			dst = srv.NewKey(args[1], vals.NewList())
			keys[args[1]] = dst
		}
		// Re-read the source key, as it may have changed during db lock upgrade
		src, ok = keys[args[0]]
		if !ok {
			// src does not exist anymore, return nil
			db.Unlock()
			return nil, nil
		}
	}
	defer def()

	// At this point, both src and dst exist and are separate keys
	src.Lock()
	defer src.Unlock()
	dst.Lock()
	defer dst.Unlock()

	// Get the values, make sure both are Lists
	vs, vd := src.Val(), dst.Val()
	vsrc, oksrc := vs.(vals.List)
	vdst, okdst := vd.(vals.List)
	if !oksrc || !okdst {
		return nil, cmd.ErrInvalidValType
	}

	// Both are lists, proceed
	val, ok := vsrc.RPop()
	if ok {
		vdst.LPush(val)
		return val, nil
	}
	return nil, nil
}

var rpush = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: -1,
	},
	srv.NoKeyCreateList,
	rpushFn)

func rpushFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		return v.RPush(args[1:]...), nil
	}
	return nil, cmd.ErrInvalidValType
}

var rpushx = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	srv.NoKeyNone,
	rpushxFn)

func rpushxFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	if k == nil {
		return int64(0), nil
	}

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
		return v.RPush(args[1:]...), nil
	}
	return nil, cmd.ErrInvalidValType
}
