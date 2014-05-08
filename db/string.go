package db

import (
	"errors"
	"sync"
)

// ErrNotAnInt is returned is the value is not an integer when an integer
// argument is expected.
var ErrNotAnInt = errors.New("db: value is not an integer")

type key struct {
	name string

	mu  sync.RWMutex
	val string
	exp *expirer
}

var cmdGet = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.mu.RLock()
			defer ctx.key.mu.RUnlock()

			return ctx.key.val, nil
		}, "", errNilSuccess), 1, 1)

var cmdSet = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			CreateKey(ctx)
			return nil, nil
		},
		func(ctx *Ctx) (interface{}, error) {
			// TODO : Remove expireation on SET
			ctx.key.mu.Lock()
			defer ctx.key.mu.Unlock()
			ctx.key.val = ctx.s1
			return nil, nil
		}), 2, 2)

var cmdAppend = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			val, err := CreateKey(ctx)
			return int64(len(val.(string))), err
		},
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.mu.Lock()
			defer ctx.key.mu.Unlock()

			ctx.key.val += ctx.s1
			return int64(len(ctx.s1)), nil
		}), 2, 2)

var cmdGetRange = CheckArgCount(
	ParseIntArgs(
		RLockExistBranch(
			func(ctx *Ctx) (interface{}, error) {
				ctx.key.mu.RLock()
				val := ctx.key.val
				ctx.key.mu.RUnlock()

				st, end := ctx.i0, ctx.i1
				if st < 0 {
					st = len(val) + st
					if st < 0 {
						st = 0
					}
				}
				if st >= len(val) {
					return "", nil
				}
				if end < 0 {
					end = len(val) + end
				}
				if end < 0 || end < st {
					return "", nil
				}
				if end >= len(val) {
					end = len(val) - 1
				}
				return val[st : end+1], nil
			}, "", nil),
	), 3, 3)

func (d *Database) getset(args ...string) (interface{}, error) {
	d.mu.RLock()
	if ky, ok := d.keys[args[0]]; !ok {
		// Key does not exist yet, must create the key
		d.mu.RUnlock()
		d.mu.Lock()
		defer d.mu.Unlock()
		ky = &key{name: args[0], val: args[1]}
		d.keys[args[0]] = ky

		return nil, errNilSuccess

	} else {
		// Key already exists, set the new value and return the old
		defer d.mu.RUnlock()
		return ky.getset(args[1]), nil
	}
}

func (k *key) getset(val string) string {
	k.mu.Lock()
	defer k.mu.Unlock()

	// getset removes expiration
	if k.exp != nil {
		k.exp.abort()
		k.exp = nil
	}
	old := k.val
	k.val = val
	return old
}

func (d *Database) strlen(args ...string) (interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if k, ok := d.keys[args[0]]; !ok {
		return int64(0), nil
	} else {
		return k.strlen(), nil
	}
}

func (k *key) strlen() int64 {
	k.mu.RLock()
	defer k.mu.RUnlock()

	return int64(len(k.val))
}

func (d *Database) exists(args ...string) (interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	_, ok := d.keys[args[0]]
	if ok {
		return int64(1), nil
	}
	return int64(0), nil
}

func (d *Database) del(args ...string) (interface{}, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var n int64
	for _, arg := range args {
		if k, ok := d.keys[arg]; ok {
			n++
			k.mu.Lock()
			if k.exp != nil {
				k.exp.abort()
			}
			delete(d.keys, arg)
			k.mu.Unlock()
		}
	}
	return n, nil
}
