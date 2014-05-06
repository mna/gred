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
	mu  sync.RWMutex
	val string
}

// get returns the value for the key at k.
func (d *Database) get(args ...string) (interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if key, ok := d.keys[args[0]]; ok {
		return key.get(), nil
	}
	return "", errNilSuccess
}

func (k *key) get() string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.val
}

// set stores the value for the key at k.
func (d *Database) set(args ...string) (interface{}, error) {
	d.mu.RLock()
	if ky, ok := d.keys[args[0]]; !ok {
		// Key does not exist yet, must create the key
		d.mu.RUnlock()
		d.mu.Lock()
		defer d.mu.Unlock()
		ky = &key{val: args[1]}
		d.keys[args[0]] = ky

	} else {
		// Key already exists, set the new value
		defer d.mu.RUnlock()
		ky.set(args[1])
	}

	return nil, nil
}

func (k *key) set(v string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.val = v
}

func (d *Database) append(args ...string) (interface{}, error) {
	var ln int64

	d.mu.RLock()
	if ky, ok := d.keys[args[0]]; !ok {
		// Key does not exist yet, must create the key
		d.mu.RUnlock()
		d.mu.Lock()
		defer d.mu.Unlock()
		ky = &key{val: args[1]}
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
