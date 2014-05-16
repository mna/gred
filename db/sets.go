package db

var cmdSadd = CheckArgCount(
	LockBothBranches(
		func(ctx *Ctx) (interface{}, error) {
			key := &setKey{
				Expirer: &expirer{},
				name:    ctx.s0,
				s:       make(set),
			}
			for _, v := range ctx.raw[1:] {
				key.s[v] = struct{}{}
			}
			ctx.db.keys[ctx.s0] = key
			return int64(len(key.s)), nil
		},
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if key, ok := ctx.key.(SetKey); ok {
				s := key.Get()
				var cnt int64
				for _, v := range ctx.raw[1:] {
					if _, ok := s[v]; !ok {
						s[v] = struct{}{}
						cnt++
					}
				}
				return cnt, nil
			}
			return nil, errInvalidKeyType
		}), 2, -1)

var cmdScard = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			if key, ok := ctx.key.(SetKey); ok {
				return int64(len(key.Get())), nil
			}
			return nil, errInvalidKeyType
		}, int64(0), nil), 1, 1)
