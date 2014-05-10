package db

var cmdExists = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			return int64(1), nil
		}, int64(0), nil), 1, 1)

var cmdDel = CheckArgCount(
	LockEachKey(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.mu.Lock()
			defer ctx.key.mu.Unlock()

			// TODO : Remove expiration
			delete(ctx.db.keys, ctx.key.name)
			return nil, nil
		}), 1, -1)
