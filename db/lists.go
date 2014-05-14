package db

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
			l.lpush(ctx.raw[1:]...)
			ctx.db.keys[ctx.s0] = key
			return int64(len(l)), nil
		},
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(ListKey); ok {
				l := key.Get()
				l.lpush(ctx.raw[1:]...)
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
