package db

import "time"

var cmdDel = CheckArgCount(
	LockEachKey(
		func(ctx *Ctx) (interface{}, error) {
			ctx.key.Lock()
			defer ctx.key.Unlock()

			// TODO : Remove expiration
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

				ctx.key.SetTimer(time.AfterFunc(ctx.i0 * time.Second))
				return int64(1), nil
			}, int64(0), nil), 1), 2, 2)
