package db

import "github.com/PuerkitoBio/gred/resp"

var cmdHdel = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(HashKey); ok {
				h := key.Get()
				var cnt int64
				for _, nm := range ctx.raw[1:] {
					if _, ok := h[nm]; ok {
						cnt++
						delete(h, nm)
					}
				}
				if len(h) == 0 {
					// Delete the key, no more fields (have to upgrade db lock)
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
				return cnt, nil
			}
			return nil, errInvalidKeyType
		}, int64(0), nil), 2, -1)

var cmdHexists = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			if key, ok := ctx.key.(HashKey); ok {
				var ex int64

				h := key.Get()
				if _, ok := h[ctx.s1]; ok {
					ex = 1
				}
				return ex, nil
			}
			return nil, errInvalidKeyType
		}, int64(0), nil), 2, 2)

var cmdHget = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(HashKey); ok {
				h := key.Get()
				if v, ok := h[ctx.s1]; ok {
					return v, nil
				}
				return nil, errNilSuccess
			}
			return nil, errInvalidKeyType
		}, nil, errNilSuccess), 2, 2)

var cmdHgetAll = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(HashKey); ok {
				h := key.Get()
				ar := make(resp.Array, 2*len(h))
				i := 0
				for k, v := range h {
					ar[i] = k
					ar[i+1] = v
					i += 2
				}
				return ar, nil
			}
			return nil, errInvalidKeyType
		}, emptyArray, nil), 1, 1)

var cmdHkeys = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(HashKey); ok {
				h := key.Get()
				ar := make(resp.Array, len(h))
				i := 0
				for k := range h {
					ar[i] = k
					i++
				}
				return ar, nil
			}
			return nil, errInvalidKeyType
		}, emptyArray, nil), 1, 1)

var cmdHlen = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(HashKey); ok {
				h := key.Get()
				return int64(len(h)), nil
			}
			return nil, errInvalidKeyType
		}, int64(0), nil), 1, 1)

var cmdHmget = CheckArgCount(
	RLockExistOrNot(
		func(ctx *Ctx) (interface{}, error) {
			if ctx.key == nil {
				// Return nils for each requested field
				return make(resp.Array, len(ctx.raw)-1), nil
			}

			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(HashKey); ok {
				h := key.Get()
				ar := make(resp.Array, len(ctx.raw)-1)
				for i := 1; i < len(ctx.raw); i++ {
					if v, ok := h[ctx.raw[i]]; ok {
						ar[i-1] = v
					}
				}
				return ar, nil
			}
			return nil, errInvalidKeyType
		}), 2, -1)

var cmdHset = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			key := &hashKey{
				Expirer: &expirer{},
				name:    ctx.s0,
				h:       hash{ctx.s1: ctx.s2},
			}
			ctx.db.keys[ctx.s0] = key
			return int64(1), nil
		},
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(HashKey); ok {
				var ret int64

				h := key.Get()
				if _, ok := h[ctx.s1]; !ok {
					ret = 1
				}
				h[ctx.s1] = ctx.s2
				return ret, nil
			}
			return nil, errInvalidKeyType
		}), 3, 3)

var cmdHvals = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(HashKey); ok {
				h := key.Get()
				ar := make(resp.Array, len(h))
				i := 0
				for _, v := range h {
					ar[i] = v
					i++
				}
				return ar, nil
			}
			return nil, errInvalidKeyType
		}, emptyArray, nil), 1, 1)
