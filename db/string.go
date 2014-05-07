package db

import (
	"errors"
	"strconv"
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

func (d *Database) append(args ...string) (interface{}, error) {
	var ln int64

	d.mu.RLock()
	if ky, ok := d.keys[args[0]]; !ok {
		// Key does not exist yet, must create the key
		d.mu.RUnlock()
		d.mu.Lock()
		defer d.mu.Unlock()
		ky = &key{name: args[0], val: args[1]}
		ln = int64(len(args[1]))
		d.keys[args[0]] = ky

	} else {
		// Key already exists, set the new value
		defer d.mu.RUnlock()
		ln = ky.append(args[1])
	}

	return ln, nil
}

func (k *key) append(v string) int64 {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.val += v
	return int64(len(k.val))
}

func (d *Database) getRange(args ...string) (interface{}, error) {
	st, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, ErrNotAnInt
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, ErrNotAnInt
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	if ky, ok := d.keys[args[0]]; ok {
		return ky.getRange(st, end), nil
	}
	return "", nil
}

func (k *key) getRange(st, end int) string {
	k.mu.RLock()
	val := k.val
	k.mu.RUnlock()

	if st < 0 {
		st = len(val) + st
		if st < 0 {
			st = 0
		}
	}
	if st >= len(val) {
		return ""
	}
	if end < 0 {
		end = len(val) + end
	}
	if end < 0 || end < st {
		return ""
	}
	if end >= len(val) {
		end = len(val) - 1
	}
	return val[st : end+1]
}

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

func (d *Database) expire(args ...string) (interface{}, error) {
	secs, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	if k, ok := d.keys[args[0]]; !ok {
		return int64(0), nil
	} else {
		if secs <= 0 {
			// Remove immediately
			// TODO : Call an impl that doesn't set the mutex
		}
	}
}
