package db

// Lindex is actually O(1) here, better than Redis' O(n)
var cmdLindex = CheckArgCount(
	ParseIntArgs(
		RLockExistBranch(
			func(ctx *Ctx) (interface{}, error) {
				ctx.key.RLock()
				defer ctx.key.RUnlock()

				if key, ok := ctx.key.(ListKey); ok {
					l := key.Get()
					ln := len(*l)
					if ctx.i0 < 0 {
						ctx.i0 = ln + ctx.i0
					}
					if ctx.i0 >= 0 && ctx.i0 < ln {
						return (*l)[ctx.i0], nil
					}
					return nil, errNilSuccess
				}
				return nil, errInvalidKeyType
			}, nil, errNilSuccess), 1), 2, 2)

var cmdLlen = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(ListKey); ok {
				l := key.Get()
				return int64(len(*l)), nil
			}
			return nil, errInvalidKeyType
		}, int64(0), nil), 1, 1)

var cmdLpop = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(ListKey); ok {
				l := key.Get()
				v, ok := l.lpop()
				if !ok {
					// Should never happen...
					return nil, errNilSuccess
				}
				if len(*l) == 0 {
					// Delete the key, no more values (have to upgrade db lock)
					ctx.db.mu.RUnlock()
					ctx.db.mu.Lock()
					// Abort any expiration
					key.Abort()
					// Delete from DB
					delete(ctx.db.keys, key.Name())
					// Downgrade db lock (so that RLockExistBranch is happy on exit)
					ctx.db.mu.Unlock()
					ctx.db.mu.RLock()
				}
				return v, nil
			}
			return nil, errInvalidKeyType
		}, nil, errNilSuccess), 1, 1)

var cmdLpush = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			l := make(list, 0, defaultListCap)
			key := &listKey{
				Expirer: &expirer{},
				name:    ctx.s0,
				l:       &l,
			}
			l.lpushr(ctx.raw[1:]...)
			ctx.db.keys[ctx.s0] = key
			return int64(len(l)), nil
		},
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(ListKey); ok {
				l := key.Get()
				l.lpushr(ctx.raw[1:]...)
				return int64(len(*l)), nil
			}
			return nil, errInvalidKeyType
		}), 2, -1)

var cmdRpop = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(ListKey); ok {
				l := key.Get()
				v, ok := l.rpop()
				if !ok {
					// Should never happen...
					return nil, errNilSuccess
				}
				if len(*l) == 0 {
					// Delete the key, no more values (have to upgrade db lock)
					ctx.db.mu.RUnlock()
					ctx.db.mu.Lock()
					// Abort any expiration
					key.Abort()
					// Delete from DB
					delete(ctx.db.keys, key.Name())
					// Downgrade db lock (so that RLockExistBranch is happy on exit)
					ctx.db.mu.Unlock()
					ctx.db.mu.RLock()
				}
				return v, nil
			}
			return nil, errInvalidKeyType
		}, nil, errNilSuccess), 1, 1)

var cmdRpush = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			l := make(list, 0, defaultListCap)
			key := &listKey{
				Expirer: &expirer{},
				name:    ctx.s0,
				l:       &l,
			}
			l.rpush(ctx.raw[1:]...)
			ctx.db.keys[ctx.s0] = key
			return int64(len(l)), nil
		},
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(ListKey); ok {
				l := key.Get()
				l.rpush(ctx.raw[1:]...)
				return int64(len(*l)), nil
			}
			return nil, errInvalidKeyType
		}), 2, -1)
