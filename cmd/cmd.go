// Package cmd defines the common command interfaces. Actual command implementations
// are registered in this package and then exposed by the server.
package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/gred/resp"
	"github.com/PuerkitoBio/gred/srv"
)

const (
	// WrongNumberOfArgsFmt is a string that holds the normalized error message
	// for when the number of arguments of a command is invalid.
	WrongNumberOfArgsFmt = "ERR wrong number of arguments for '%s' command"
)

var (
	// PongVal is a sentinel value used to return the PONG response for the PING command.
	PongVal = resp.Pong{}

	// OKVal is a sentinel value used to indicate that the standard OK simple string
	// response should be returned.
	OKVal = resp.OK{}

	// ErrNotInteger is returned when an argument that is expected to be an integer
	// cannot be parsed as an integer.
	ErrNotInteger = errors.New("ERR value is not an integer or out of range")

	// ErrNotFloat is returned when an argument that is expected to be a float
	// cannot be parsed as a float.
	ErrNotFloat = errors.New("ERR value is not a valid float")

	// ErrInvalidValType is returned when the key's value is not of the expected type.
	ErrInvalidValType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")

	// ErrSyntax is returned when an argument doesn't have the expected allowed syntax.
	ErrSyntax = errors.New("ERR syntax error")

	// ErrNoSuchKey is returned when a command is attempted against a non-existing key.
	ErrNoSuchKey = errors.New("ERR no such key")

	// ErrOutOfRange is returned when an argument is out of range.
	ErrOutOfRange = errors.New("ERR index out of range")

	// ErrOfsOutOfRange is returned when an offset argument is out of range.
	ErrOfsOutOfRange = errors.New("ERR offset is out of range")

	// ErrHashFieldNotInt is returned when an increment operation is attempted on
	// a hash field that does not contain an integer value.
	ErrHashFieldNotInt = errors.New("ERR hash value is not an integer")

	// ErrHashFieldNotFloat is returned when an increment operation is attempted on
	// a hash field that does not contain a float value.
	ErrHashFieldNotFloat = errors.New("ERR hash value is not a valid float")

	// ErrQuit is a sentinel error value to indicate that the network connection
	// should be closed, as requested by the client.
	ErrQuit = errors.New("quit")

	// ErrInvalidDBIndex is returned when a DB index outside the bounds of
	// available DBs is requested.
	ErrInvalidDBIndex = errors.New("ERR invalid DB index")
)

// Commands holds the list of registered commands.
var Commands = make(map[string]Cmd)

// Register registers a command name with an implementation.
func Register(name string, c Cmd) {
	if name == "" {
		panic("cmds: call Register with empty command name")
	}
	if _, ok := Commands[name]; ok {
		panic(fmt.Sprintf("cmds: command %s already registered", name))
	}
	Commands[name] = c
}

// Cmd defines the common method required to implement a basic command.
// It is insufficient to implement this sole interface. A command must also
// implement one of the more specific {Srv,DB}Cmd interfaces.
type Cmd interface {
	Parse(string, []string) ([]string, []int64, []float64, error)
}

// SrvFn defines the function signature required for the SrvCmd implementation.
type SrvFn func([]string, []int64, []float64) (interface{}, error)

// SrvCmd defines the methods required to implement a server command.
type SrvCmd interface {
	Cmd
	Exec([]string, []int64, []float64) (interface{}, error)
}

// NewSrvCmd creates a new SrvCmd value with the specified argument definition
// and execution function.
func NewSrvCmd(arg *ArgDef, fn SrvFn) SrvCmd {
	return &srvCmd{
		arg,
		fn,
	}
}

// srvCmd implements a server command, which is a type of command that doesn't
// act on a specific DB or Key.
type srvCmd struct {
	*ArgDef
	fn SrvFn
}

// Exec executes the server command with the provided arguments.
func (s *srvCmd) Exec(args []string, ints []int64, floats []float64) (interface{}, error) {
	return s.fn(args, ints, floats)
}

// ConnFn defines the function signature required for the ConnCmd implementation.
type ConnFn func(srv.Conn, []string, []int64, []float64) (interface{}, error)

// ConnCmd defines the methods required to implement a connection command.
type ConnCmd interface {
	Cmd
	ExecWithConn(srv.Conn, []string, []int64, []float64) (interface{}, error)
}

// NewConnCmd creates a new ConnCmd value with the specified argument definition
// and execution function.
func NewConnCmd(arg *ArgDef, fn ConnFn) ConnCmd {
	return &connCmd{
		arg,
		fn,
	}
}

// connCmd is the interal implementation of a ConnCmd.
type connCmd struct {
	*ArgDef
	fn ConnFn
}

// Exec executes the connection command with the provided arguments.
func (c *connCmd) ExecWithConn(conn srv.Conn, args []string, ints []int64, floats []float64) (interface{}, error) {
	return c.fn(conn, args, ints, floats)
}

// DBCmd defines the methods required to implement a database command.
type DBCmd interface {
	Cmd
	ExecWithDB(srv.DB, []string, []int64, []float64) (interface{}, error)
}

// ArgFn defines the function signature required for the argument validation function.
type ArgFn func([]string, []int64, []float64) error

// ArgDef holds the argument definition for a command.
type ArgDef struct {
	// Indices of arguments to parse as floats.
	FloatIndices []int

	// Indices of arguments to parse as integers.
	IntIndices []int

	// Min and Max number of arguments. -1 can be specified for
	// unbounded (variadic) maximum number of arguments.
	MinArgs, MaxArgs int

	// ValidateFn is a function that is called (if set) to provide custom
	// argument validation. It is called after the parsing of arguments as
	// floats and integers, if applicable.
	ValidateFn ArgFn
}

// Parse parses the provided list of arguments according to the argument
// definition specs. It returns the list of arguments, the parsed integers,
// the parsed floats, and an error if the arguments are invalid.
func (a *ArgDef) Parse(name string, args []string) ([]string, []int64, []float64, error) {
	l := len(args)
	if l < a.MinArgs || (l > a.MaxArgs && a.MaxArgs >= 0) {
		return nil, nil, nil, fmt.Errorf(WrongNumberOfArgsFmt, name)
	}

	// Parse integers
	intix := a.IntIndices
	ints := make([]int64, len(intix))
	for i, ix := range intix {
		if ix < 0 {
			ix = l + ix
		}
		val, err := strconv.ParseInt(args[ix], 10, 64)
		if err != nil {
			return nil, nil, nil, ErrNotInteger
		}
		ints[i] = val
	}

	// Parse floats
	fix := a.FloatIndices
	floats := make([]float64, len(fix))
	for i, ix := range fix {
		if ix < 0 {
			ix = l + ix
		}
		val, err := strconv.ParseFloat(args[ix], 64)
		if err != nil {
			return nil, nil, nil, ErrNotFloat
		}
		floats[i] = val
	}

	if a.ValidateFn != nil {
		err := a.ValidateFn(args, ints, floats)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return args, ints, floats, nil
}

// Static type check that *singleKeyCmd implements the DBCmd interface.
var _ DBCmd = (*singleKeyCmd)(nil)

// KeyFn defines the function signature required for a Key command function.
type KeyFn func(srv.Key, []string, []int64, []float64) (interface{}, error)

// NewSingleKeyCmd creates a DBCmd that acts on a single key. The noKeyFlag
// argument specifies what the command should do if the key does not exist.
func NewSingleKeyCmd(arg *ArgDef, noKeyFlag srv.NoKeyFlag, fn KeyFn) DBCmd {
	return &singleKeyCmd{
		ArgDef: arg,
		noKey:  noKeyFlag,
		fn:     fn,
	}
}

// singleKeyCmd implements DBCmd for commands that act on a single key.
type singleKeyCmd struct {
	*ArgDef
	noKey srv.NoKeyFlag
	fn    KeyFn
}

// ExecWithDB executes the command with the provided database and arguments.
func (c *singleKeyCmd) ExecWithDB(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	k, def := db.LockGetKey(args[0], c.noKey)
	defer def()

	return c.fn(k, args, ints, floats)
}

// Static type check to validate that *dbCmd implements DBCmd.
var _ DBCmd = (*dbCmd)(nil)

// DBFn defines the function signature required for DB command functions.
type DBFn func(srv.DB, []string, []int64, []float64) (interface{}, error)

// NewDBCmd creates a new command that acts on a database.
func NewDBCmd(arg *ArgDef, fn DBFn) DBCmd {
	return &dbCmd{
		ArgDef: arg,
		fn:     fn,
	}
}

// dbCmd implements a DBCmd that acts on a database.
type dbCmd struct {
	*ArgDef
	fn DBFn
}

// ExecWithDB executes the command with the specified database and arguments.
func (d *dbCmd) ExecWithDB(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	return d.fn(db, args, ints, floats)
}
