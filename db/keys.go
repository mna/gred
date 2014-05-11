package db

import "time"

var cmdDel = CheckArgCount(
	LockEachKey(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			ctx.key.Abort()
			delete(ctx.db.keys, ctx.key.Name())
			return nil, nil
		}), 1, -1)

var cmdExists = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			return int64(1), nil
		}, int64(0), nil), 1, 1)

var cmdExpire = CheckArgCount(
	ParseIntArgs(
		RLockExistBranch(
			func(ctx *Ctx) (interface{}, error) {
				ctx.key.Lock()
				defer ctx.key.Unlock()

				// Avoid closing over the original ctx (and its Conn).
				ctxDel := &Ctx{
					db: ctx.db,
					s0: ctx.s0,
				}
				if ctx.key.Expire(
					time.Now().Add(time.Duration(ctx.i0)*time.Second),
					func() {
						LockExistBranch(
							func(ctx *Ctx) (interface{}, error) {
								ctx.key.Lock()
								defer ctx.key.Unlock()

								delete(ctx.db.keys, ctx.key.Name())
								return nil, nil
							}, nil, nil)(ctxDel)
					}) {
					return int64(1), nil
				}
				return int64(0), nil
			}, int64(0), nil), 1), 2, 2)

var cmdPersist = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			if ctx.key.Abort() {
				return int64(1), nil
			}
			return int64(0), nil
		}, int64(0), nil), 1, 1)

var cmdTTL = CheckArgCount(
	RLockExistBranch(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.RLock()
			defer ctx.key.RUnlock()

			dur := ctx.key.TTL()
			if dur == -1 {
				return int64(-1), nil
			}
			return int64(dur.Seconds()), nil
		}, int64(-2), nil), 1, 1)
