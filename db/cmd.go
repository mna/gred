package db

type Ctx struct {
	conn *Conn
	db   *Database
	key  *key

	s0, s1, s2 string
	i0, i1, i2 int

	ss []string
	is []int
}

type CmdFunc func(*Ctx) (interface{}, error)

func (fn CmdFunc) Do(c *Ctx) (interface{}, error) {
	return fn(c)
}

type Cmd interface {
	Do(c *Ctx) (interface{}, error)
}

func CreateKey(ctx *Ctx) (interface{}, error) {
	ky := &key{name: ctx.s0, val: ctx.s1}
	ctx.db.keys[ctx.s0] = ky

	return ctx.s1, nil
}

func LockKey(cmd Cmd) Cmd {
	return CmdFunc(func(ctx *Ctx) (interface{}, error) {
		ctx.key.mu.Lock()
		defer ctx.key.mu.Unlock()
		return cmd.Do(ctx)
	})
}

func RLockExistBranch(cmd Cmd, defRes interface{}, defErr error) Cmd {
	return CmdFunc(func(ctx *Ctx) (interface{}, error) {
		ctx.db.mu.RLock()
		defer ctx.db.mu.RUnlock()

		if key, ok := ctx.db.keys[ctx.s0]; ok {
			ctx.key = key
			return cmd.Do(ctx)
		}
		return defRes, defErr
	})
}

func LockTwoBranches(cmdNotExist, cmdExist Cmd) Cmd {
	return CmdFunc(func(ctx *Ctx) (interface{}, error) {
		ctx.db.mu.RLock()
		if ky, ok := ctx.db.keys[ctx.s0]; !ok {
			// Key does not exist yet
			ctx.db.mu.RUnlock()
			ctx.db.mu.Lock()
			defer ctx.db.mu.Unlock()
			return cmdNotExist.Do(ctx)

		} else {
			// Key already exists
			defer ctx.db.mu.RUnlock()
			ctx.key = ky
			return cmdExist.Do(ctx)
		}
	})
}
