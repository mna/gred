package db

import (
	"errors"
	"sync"
)

type cmdFunc func(db *Database, args ...string) (interface{}, error)

type cmdDef struct {
	fn    cmdFunc
	nArgs int
}

var cmds = map[string]cmdDef{
	"get":      cmdDef{(*Database).get, 1},
	"set":      cmdDef{(*Database).set, 2},
	"append":   cmdDef{(*Database).append, 2},
	"getrange": cmdDef{(*Database).getRange, 3},
	"substr":   cmdDef{(*Database).getRange, 3},
	"getset":   cmdDef{(*Database).getset, 2},
	"strlen":   cmdDef{(*Database).strlen, 1},
	"exists":   cmdDef{(*Database).exists, 1},
	"del":      cmdDef{(*Database).del, 1},
}

var (
	// ErrInvalidCommand is returned when a malformed command is received.
	ErrInvalidCommand = errors.New("db: invalid command")

	// ErrMissingArg is returned if there are not enough arguments to call
	// the specified command.
	ErrWrongNumberOfArgs = errors.New("db: wrong number of arguments")

	// ErrNilSuccess is a sentinel value to indicate the success of a command,
	// and that the nil value should be returned.
	errNilSuccess = errors.New("db: (nil)")
)

// Database represents a Redis database, identified by its index.
type Database struct {
	ix   int
	mu   sync.RWMutex
	keys map[string]*key
}

// NewDB creates a new Database identified by the specified index.
func NewDB(index int) *Database {
	return &Database{
		ix:   index,
		keys: make(map[string]*key),
	}
}

// Do executes the command cmd with the specified arguments args.
func (d *Database) Do(cmd string, args ...string) (interface{}, error) {
	if def, ok := cmds[cmd]; ok {
		if len(args) < def.nArgs {
			return nil, ErrMissingArg
		}
		return def.fn(d, args...)
	}
	return nil, ErrInvalidCommand
}
