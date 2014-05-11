package db

import (
	"errors"
	"strconv"
)

var (
	// errWrongNumberOfArgs is returned if there are not enough or too many arguments to call
	// the specified command.
	errWrongNumberOfArgs = errors.New("db: wrong number of arguments")

	// errNotAnInt is returned is the value is not an integer when an integer
	// argument is expected.
	errNotAnInt = errors.New("db: value is not an integer")
)

// Ctx holds the execution context for a command.
type Ctx struct {
	conn *Conn
	db   *Database
	key  Key

	s0, s1, s2 string
	i0, i1, i2 int

	raw []string
}

// Cmd represents a command's function.
type Cmd func(*Ctx) (interface{}, error)

// CheckArgCount verifies the argument count based on the
// specified requirements. It stores the string values in
// the predefined s0, s1 and s2 Ctx registers.
func CheckArgCount(cmd Cmd, min, max int) Cmd {
	return func(ctx *Ctx) (interface{}, error) {
		l := len(ctx.raw)
		if l < min || (l > max && max >= 0) {
			return nil, errWrongNumberOfArgs
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

// ParseIntArgs parses as integers the arguments of the Ctx identified by the indices
// and stores the first 3 integers in i0, i1 and i2 on the Ctx.
func ParseIntArgs(cmd Cmd, indices ...int) Cmd {
	return func(ctx *Ctx) (interface{}, error) {
		for i, ix := range indices {
			if ix >= 0 && ix < len(ctx.raw) {
				j, err := strconv.Atoi(ctx.raw[ix])
				if err != nil {
					return nil, errNotAnInt
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

// LockEachKey acquires a Lock on the database and calls cmd for each key identified
// by the raw arguments.
func LockEachKey(cmd Cmd) Cmd {
	return Cmd(func(ctx *Ctx) (interface{}, error) {
		ctx.db.mu.Lock()
		defer ctx.db.mu.Unlock()

		n := 0
		for _, k := range ctx.raw {
			if ky, ok := ctx.db.keys[k]; ok {
				ctx.key = ky
				_, err := cmd(ctx)
				if err != nil {
					return n, err
				}
				n++
			}
		}
		return int64(n), nil
	})
}

// RLockExistBranch read-locks the database and calls cmd with the key identified
// by ctx.s0 if the key exists. It returns defRes and defErr is no such key exists.
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

// LockExistBranch acquires a Lock on the database and calls cmd for the key
// identified by ctx.s0 if such key exists. It returns defRes and defErr is no
// such key exists.
func LockExistBranch(cmd Cmd, defRes interface{}, defErr error) Cmd {
	return Cmd(func(ctx *Ctx) (interface{}, error) {
		ctx.db.mu.Lock()
		defer ctx.db.mu.Unlock()

		if key, ok := ctx.db.keys[ctx.s0]; ok {
			ctx.key = key
			return cmd(ctx)
		}
		return defRes, defErr
	})
}

// LockBothBranches read-locks the database, and upgrades the lock to a read-write
// lock if the key identified by ctx.s0 does not exist (and calls cmdNotExist), or it
// keeps the read-only lock and calls cmdExist it such key exists.
func LockBothBranches(cmdNotExist, cmdExist Cmd) Cmd {
	return func(ctx *Ctx) (interface{}, error) {
		ctx.db.mu.RLock()

		ky, ok := ctx.db.keys[ctx.s0]
		if !ok {
			// Key does not exist yet
			ctx.db.mu.RUnlock()
			ctx.db.mu.Lock()
			defer ctx.db.mu.Unlock()
			return cmdNotExist(ctx)
		}

		// Key already exists
		defer ctx.db.mu.RUnlock()
		ctx.key = ky
		return cmdExist(ctx)
	}
}
