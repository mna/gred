package db

import "sync"

var _ StringKey = (*stringKey)(nil)

type stringKey struct {
	sync.RWMutex

	name string
	val  string
}

func (s *stringKey) Name() string {
	return s.name
}

func (s *stringKey) Get() string {
	return s.val
}

func (s *stringKey) Set(v string) {
	s.val = v
}

var cmdAppend = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			val, err := CreateKey(ctx)
			return int64(len(val.(string))), err
		},
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(StringKey); ok {
				v := key.Get()
				v += ctx.s1
				key.Set(v)
				return int64(len(ctx.s1)), nil
			}
			return nil, errInvalidKeyType
		}), 2, 2)

var cmdGet = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(StringKey); ok {
				return key.Get(), nil
			}
			return nil, errInvalidKeyType
		}, "", errNilSuccess), 1, 1)

var cmdGetRange = CheckArgCount(
	ParseIntArgs(
		RLockExistBranch(
			func(ctx *Ctx) (interface{}, error) {
				var val string

				ctx.key.RLock()
				if key, ok := ctx.key.(StringKey); ok {
					val = key.Get()
				} else {
					ctx.key.RUnlock()
					return nil, errInvalidKeyType
				}
				ctx.key.RUnlock()

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
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(StringKey); ok {
				old := key.Get()
				key.Set(ctx.s1)
				return old, nil
			}
			return nil, errInvalidKeyType
		}), 2, 2)

var cmdSet = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			CreateKey(ctx)
			return nil, nil
		},
		func(ctx *Ctx) (interface{}, error) {
			// TODO : Remove expiration on SET
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(StringKey); ok {
				key.Set(ctx.s1)
				return nil, nil
			}
			return nil, errInvalidKeyType
		}), 2, 2)

var cmdStrLen = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(StringKey); ok {
				return int64(len(key.Get())), nil
			}
			return nil, errInvalidKeyType
		}, int64(0), nil), 1, 1)
