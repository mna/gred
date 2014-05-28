package lists

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
)

func init() {
	cmd.Register("blpop", blpop)
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
	db.Lock()
	unlocks := make([]func(), 0)
	unlocks = append(unlocks, db.Unlock)

	keys := db.Keys()
	for _, nm := range args[:len(args)-1] { // last arg is the timeout
		k, ok := keys[nm]
		// Ignore non-existing keys in non-blocking portion
		if ok {
			// Lock the key
			k.Lock()
			unlocks = append(unlocks, k.Unlock)

			// Get the value, if possible
			v := k.Val()
			if v, ok := v.(vals.List); ok {
				val, ok := v.LPop()
				if ok {
					// Delete the key if there are no more values
					if v.LLen() == 0 {
						db.DelKey(k.Name())
					}

					// Unlock all keys in reverse order, and return
					for i := len(unlocks) - 1; i >= 0; i-- {
						unlocks[i]()
					}
					return []string{k.Name(), val}, nil
				}
			} else {
				// Unlock all keys in reverse order
				for i := len(unlocks) - 1; i >= 0; i-- {
					unlocks[i]()
				}
				// Return invalid type error
				return nil, cmd.ErrInvalidValType
			}
		}
	}

	// If no value was readily available, now all keys are locked, enter
	// the waiting workflow.
	ch := make(chan chan<- [2]string)
	for _, nm := range args[:len(args)-1] {
		db.WaitLPop(nm, ch)
	}

	// Prepare channels (timeout and receive values)
	var timeoutCh <-chan time.Time
	if ints[0] > 0 {
		timeoutCh = time.After(time.Duration(ints[0]) * time.Second)
	}
	recCh := make(chan [2]string)

	// Unlock all locks so that other connections can proceed
	for i := len(unlocks) - 1; i >= 0; i-- {
		unlocks[i]()
	}

	// Wait for a value
	select {
	case ch <- (chan<- [2]string)(recCh):
		close(ch)
		vals := <-recCh
		return vals[:], nil
	case <-timeoutCh:
		close(ch)
		return nil, nil
	}
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
	if v, ok := v.(vals.List); ok {
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
	// DB must be exclusively locked, because of the unblock behaviour.
	k, unl := db.XLockGetKey(args[0], srv.NoKeyCreateList)
	defer unl()

	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.List); ok {
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
	if v, ok := v.(vals.List); ok {
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
	if v, ok := v.(vals.List); ok {
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
	if v, ok := v.(vals.List); ok {
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
	if v, ok := v.(vals.List); ok {
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

	src.Lock()
	defer src.Unlock()
	// Exit early if the source key does not hold a List so that the dst
	// is not created.
	vs := src.Val()
	vsrc, ok := vs.(vals.List)
	if !ok {
		return nil, cmd.ErrInvalidValType
	}

	// Otherwise get the destination key, and create it if it doesn't exist
	dst, ok := keys[args[1]]
	if !ok {
		// Destination does not exist, create it
		dst = srv.NewKey(args[1], vals.NewList())
		keys[args[1]] = dst
	}

	dst.Lock()
	defer dst.Unlock()
	// Get the dst value, make sure it is a List
	vd := dst.Val()
	vdst, ok := vd.(vals.List)
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
