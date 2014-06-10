package lists

import (
	"strings"

	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/types"
)

func init() {
	cmd.Register("blpop", blpop)
	cmd.Register("brpop", brpop)
	cmd.Register("brpoplpush", brpoplpush)
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

var blpop = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    -1,
		IntIndices: []int{-1},
	},
	blpopFn)

func blpopFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	ar, err := blockPop(db, ints[0], false, args[:len(args)-1]...)
	if ar == nil {
		return nil, err
	}
	return ar, nil
}

var brpop = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    2,
		MaxArgs:    -1,
		IntIndices: []int{-1},
	},
	brpopFn)

func brpopFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	ar, err := blockPop(db, ints[0], true, args[:len(args)-1]...)
	if ar == nil {
		return nil, err
	}
	return ar, nil
}

var brpoplpush = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{2},
	},
	brpoplpushFn)

func brpoplpushFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	// First do the brpop part
	vals, err := blockPop(db, ints[0], true, args[0])
	if vals == nil {
		// Return either an error, or the nil timeout value
		return nil, err
	}

	// Then proceed with lpush
	_, err = lpushFn(db, []string{args[1], vals[1]}, nil, nil)
	if err != nil {
		return nil, err
	}

	// Return the value popped and pushed
	return vals[1], nil
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
	if v, ok := v.(types.List); ok {
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
	if v, ok := v.(types.List); ok {
		if args[1] == "before" {
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
	if v, ok := v.(types.List); ok {
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
	// Since LPOP may delete the key (if the list is empty), must get an exclusive
	// DB lock right away (can't think of a sane way to upgrade the lock without restartint
	// the whole operation).
	k, unl := db.XLockGetKey(args[0], srv.NoKeyDefaultVal)
	defer unl()

	// Lock the key
	k.Lock()
	defer k.Unlock()

	// Pop the value
	v := k.Val()
	if v, ok := v.(types.List); ok {
		val, ok := v.LPop()
		if ok {
			// If the list is now empty, delete the key
			if v.LLen() == 0 {
				db.DelKey(args[0])
			}
			return val, nil
		}
		return nil, nil
	}
	return nil, cmd.ErrInvalidValType
}

var lpush = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: -1,
	},
	lpushFn)

func lpushFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	// DB must be exclusively locked, because of the unblock behaviour, which
	// may result in a delete of the key.
	k, unl := db.XLockGetKey(args[0], srv.NoKeyCreateList)
	defer unl()

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.List); ok {
		val := v.LPush(args[1:]...)
		// Unblock any waiters on this key
		if unblock(db, k, v) > 0 {
			// If the list is now empty, delete the key
			if v.LLen() == 0 {
				db.DelKey(args[0])
			}
		}
		return val, nil
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

// LPUSHX can't have any waiters, because it only pushes if the key already
// exists, and the key is removed if it doesn't have any value.
func lpushxFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	if k == nil {
		return int64(0), nil
	}

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.List); ok {
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
	if v, ok := v.(types.List); ok {
		return v.LRange(ints[0], ints[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var lrem = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{1},
	},
	lremFn)

func lremFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	// Since LREM may delete the key (if the list is empty), must get an exclusive
	// DB lock right away (can't think of a sane way to upgrade the lock without restartint
	// the whole operation).
	k, unl := db.XLockGetKey(args[0], srv.NoKeyDefaultVal)
	defer unl()

	// Lock the key
	k.Lock()
	defer k.Unlock()

	// Remove the value(s)
	v := k.Val()
	if v, ok := v.(types.List); ok {
		val := v.LRem(ints[0], args[2])
		if val > 0 {
			// If the list is now empty, delete the key
			if v.LLen() == 0 {
				db.DelKey(args[0])
			}
		}
		return val, nil
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
	if v, ok := v.(types.List); ok {
		ok := v.LSet(ints[0], args[2])
		if ok {
			return cmd.OKVal, nil
		}
		return nil, cmd.ErrOutOfRange
	}
	return nil, cmd.ErrInvalidValType
}

var ltrim = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{1, 2},
	},
	ltrimFn)

func ltrimFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	// Since LTRIM may delete the key (if the list is empty), must get an exclusive
	// DB lock right away (can't think of a sane way to upgrade the lock without restartint
	// the whole operation).
	k, unl := db.XLockGetKey(args[0], srv.NoKeyDefaultVal)
	defer unl()

	// Lock the key
	k.Lock()
	defer k.Unlock()

	// Trim the value
	v := k.Val()
	if v, ok := v.(types.List); ok {
		v.LTrim(ints[0], ints[1])
		// If the list is now empty, delete the key
		if v.LLen() == 0 {
			db.DelKey(args[0])
		}
		return cmd.OKVal, nil
	}
	return nil, cmd.ErrInvalidValType
}

var rpop = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	rpopFn)

func rpopFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	// Since RPOP may delete the key (if the list is empty), must get an exclusive
	// DB lock right away (can't think of a sane way to upgrade the lock without restartint
	// the whole operation).
	k, unl := db.XLockGetKey(args[0], srv.NoKeyDefaultVal)
	defer unl()

	// Lock the key
	k.Lock()
	defer k.Unlock()

	// Pop the value
	v := k.Val()
	if v, ok := v.(types.List); ok {
		val, ok := v.RPop()
		if ok {
			// If the list is now empty, delete the key
			if v.LLen() == 0 {
				db.DelKey(args[0])
			}
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
	// Since RPOPLPUSH may delete the key (if the src list is empty), must get an exclusive
	// DB lock right away (can't think of a sane way to upgrade the lock without restartint
	// the whole operation).
	db.Lock()
	defer db.Unlock()

	// Get the source key
	keys := db.Keys()
	src, ok := keys[args[0]]
	if !ok {
		// Source key does not exist, return nil
		return nil, nil
	}

	// If source exists, and source and destination are the same
	if args[0] == args[1] {
		src.Lock()
		defer src.Unlock()

		// Simply rotate the value (pop at tail, push at head)
		// Do not check if list is empty to delete it, since we push back
		// to the same list.
		v := src.Val()
		if v, ok := v.(types.List); ok {
			val, ok := v.RPop()
			if ok {
				v.LPush(val)
				return val, nil
			}
			return nil, nil
		}
		return nil, cmd.ErrInvalidValType
	}

	src.Lock()
	defer src.Unlock()
	// Exit early if the source key does not hold a List so that the dst
	// is not created.
	vs := src.Val()
	vsrc, ok := vs.(types.List)
	if !ok {
		return nil, cmd.ErrInvalidValType
	}

	// Otherwise get the destination key, and create it if it doesn't exist
	dst, ok := keys[args[1]]
	if !ok {
		// Destination does not exist, create it
		dst = srv.NewKey(args[1], types.NewList())
		keys[args[1]] = dst
	}

	dst.Lock()
	defer dst.Unlock()
	// Get the dst value, make sure it is a List
	vd := dst.Val()
	vdst, ok := vd.(types.List)
	if !ok {
		return nil, cmd.ErrInvalidValType
	}

	// Both are lists, proceed
	val, ok := vsrc.RPop()
	if ok {
		// Check if the src is now empty, if so delete the key
		if vsrc.LLen() == 0 {
			db.DelKey(args[0])
		}
		vdst.LPush(val)
		// Unblock any waiters on the dst key
		if unblock(db, dst, vdst) > 0 {
			// If the list is now empty, delete the key
			if vdst.LLen() == 0 {
				db.DelKey(args[1])
			}
		}
		return val, nil
	}
	return nil, nil
}

var rpush = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: -1,
	},
	rpushFn)

func rpushFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	// DB must be exclusively locked, because of the unblock behaviour, which
	// may result in a delete of the key.
	k, unl := db.XLockGetKey(args[0], srv.NoKeyCreateList)
	defer unl()

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.List); ok {
		val := v.RPush(args[1:]...)
		// Unblock any waiters on this key
		if unblock(db, k, v) > 0 {
			// If the list is now empty, delete the key
			if v.LLen() == 0 {
				db.DelKey(args[0])
			}
		}
		return val, nil
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

// RPUSHX can't have any waiters, because it only pushes if the key already
// exists, and the key is removed if it doesn't have any value.
func rpushxFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	if k == nil {
		return int64(0), nil
	}

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.List); ok {
		return v.RPush(args[1:]...), nil
	}
	return nil, cmd.ErrInvalidValType
}
