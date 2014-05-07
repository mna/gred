package db

type Ctx struct {
	conn *Conn
	db   *Database
	key  *key

	s0, s1, s2 string
	i0, i1, i2 int

	raw []string
}

type Cmd func(*Ctx) (interface{}, error)

func CreateKey(ctx *Ctx) (interface{}, error) {
	ky := &key{name: ctx.s0, val: ctx.s1}
	ctx.db.keys[ctx.s0] = ky

	return ctx.s1, nil
}

func CheckArgCount(cmd Cmd, min, max int) Cmd {
	return func(ctx *Ctx) (interface{}, error) {
		l := len(ctx.raw)
		if l < min || l > max {
			return nil, ErrWrongNumberOfArgs
		}
		// Store in s0..s2
		if l > 0 {
			ctx.s0 = ctx.raw[0]
		}
		if l > 1 {
			ctx.s1 = ctx.raw[1]
		}
		if l > 2 {
			ctx.s2 = ctx.raw[2]
		}
		return cmd(ctx)
	}
}

func LockKey(cmd Cmd) Cmd {
	return func(ctx *Ctx) (interface{}, error) {
		ctx.key.mu.Lock()
		defer ctx.key.mu.Unlock()
		return cmd(ctx)
	}
}

func RLockExistBranch(cmd Cmd, defRes interface{}, defErr error) Cmd {
	return Cmd(func(ctx *Ctx) (interface{}, error) {
		ctx.db.mu.RLock()
		defer ctx.db.mu.RUnlock()

		if key, ok := ctx.db.keys[ctx.s0]; ok {
			ctx.key = key
			return cmd(ctx)
		}
		return defRes, defErr
	})
}

func LockBothBranches(cmdNotExist, cmdExist Cmd) Cmd {
	return func(ctx *Ctx) (interface{}, error) {
		ctx.db.mu.RLock()
		if ky, ok := ctx.db.keys[ctx.s0]; !ok {
			// Key does not exist yet
			ctx.db.mu.RUnlock()
			ctx.db.mu.Lock()
			defer ctx.db.mu.Unlock()
			return cmdNotExist(ctx)

		} else {
			// Key already exists
			defer ctx.db.mu.RUnlock()
			ctx.key = ky
			return cmdExist(ctx)
		}
	}
}
