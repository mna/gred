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

var cmdGet = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.mu.RLock()
			defer ctx.key.mu.RUnlock()

			return ctx.key.val, nil
		}, "", errNilSuccess), 1, 1)

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

var cmdGetSet = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			CreateKey(ctx)
			return nil, errNilSuccess
		},
		func(ctx *Ctx) (interface{}, error) {
			// TODO : Remove expiration on GETSET
			ctx.key.mu.Lock()
			defer ctx.key.mu.Unlock()
			old := ctx.key.val
			ctx.key.val = ctx.s1
			return old, nil
		}), 2, 2)

var cmdSet = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			CreateKey(ctx)
			return nil, nil
		},
		func(ctx *Ctx) (interface{}, error) {
			// TODO : Remove expiration on SET
			ctx.key.mu.Lock()
			defer ctx.key.mu.Unlock()
			ctx.key.val = ctx.s1
			return nil, nil
		}), 2, 2)

var cmdStrLen = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.mu.RLock()
			defer ctx.key.mu.RUnlock()

			return int64(len(ctx.key.val)), nil
		}, int64(0), nil), 1, 1)
