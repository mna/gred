package db

var cmdHdel = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			// TODO : Removing the last field should delete the key
			if key, ok := ctx.key.(HashKey); ok {
				h := key.Get()
				var cnt int64
				for _, nm := range ctx.raw[1:] {
					if _, ok := h[nm]; ok {
						cnt++
						delete(h, nm)
					}
				}
				return cnt, nil
			}
			return nil, errInvalidKeyType
		}, int64(0), nil), 2, -1)

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
