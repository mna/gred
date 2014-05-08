package db

import "strconv"

// Ctx holds the execution context for a command.
type Ctx struct {
	conn *Conn
	db   *Database
	key  *key

	s0, s1, s2 string
	i0, i1, i2 int

	raw []string
}

// Cmd represents a command's function.
type Cmd func(*Ctx) (interface{}, error)

// CreateKey is a Cmd function that creates a string key
// using the current Ctx values.
func CreateKey(ctx *Ctx) (interface{}, error) {
	ky := &key{name: ctx.s0, val: ctx.s1}
	ctx.db.keys[ctx.s0] = ky

	return ctx.s1, nil
}

// CheckArgCount verifies the argument count based on the
// specified requirements. It stores the string values in
// the predefined s0, s1 and s2 Ctx registers.
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

func ParseIntArgs(cmd Cmd, indices ...int) Cmd {
	return func(ctx *Ctx) (interface{}, error) {
		for i, ix := range indices {
			if ix >= 0 && ix < len(ctx.raw) {
				j, err := strconv.Atoi(ctx.raw[ix])
				if err != nil {
					return nil, ErrNotAnInt
				}
				switch i {
				case 0:
					ctx.i0 = j
				case 1:
					ctx.i1 = j
				case 2:
					ctx.i2 = j
				}
			}
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
