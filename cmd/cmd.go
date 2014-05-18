package cmd

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/gred/srv"
)

var (
	ErrInvalidValType = errors.New("ERR Operation against a key holding the wrong kind of value")
	ErrNilSuccess     = errors.New("nil")
	ErrPong           = errors.New("pong")
)

var Commands = make(map[string]Cmd)

func Register(name string, c Cmd) {
	if name == "" {
		panic("cmds: call Register with empty command name")
	}
	if _, ok := Commands[name]; ok {
		panic(fmt.Sprintf("cmds: command %s already registered", name))
	}
	Commands[name] = c
}

type Cmd interface {
	GetArgDef() *ArgDef
}

type DBCmd interface {
	Cmd
	ExecWithDB(srv.DB, []string, []int64, []float64) (interface{}, error)
}

type ArgDef struct {
	FloatIndices     []int
	IntIndices       []int
	MinArgs, MaxArgs int
}

func (a *ArgDef) GetArgDef() *ArgDef { return a }

var _ DBCmd = (*singleKeyCmd)(nil)

type KeyFn func(srv.Key, []string, []int64, []float64) (interface{}, error)

func NewSingleKeyCmd(arg *ArgDef, noKeyFlag srv.NoKeyFlag, fn KeyFn) DBCmd {
	return &singleKeyCmd{
		ArgDef: arg,
		noKey:  noKeyFlag,
		fn:     fn,
	}
}

type singleKeyCmd struct {
	*ArgDef
	noKey srv.NoKeyFlag
	fn    KeyFn
}

func (c *singleKeyCmd) ExecWithDB(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	k, def := db.LockGetKey(args[0], c.noKey)
	defer def()

	return c.fn(k, args, ints, floats)
}

var _ DBCmd = (*dbCmd)(nil)

type DBFn func(srv.DB, []string, []int64, []float64) (interface{}, error)

func NewDBCmd(arg *ArgDef, fn DBFn) DBCmd {
	return &dbCmd{
		ArgDef: arg,
		fn:     fn,
	}
}

type dbCmd struct {
	*ArgDef
	fn DBFn
}

func (d *dbCmd) ExecWithDB(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	return d.fn(db, args, ints, floats)
}
