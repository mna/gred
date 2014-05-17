package cmds

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/gred/srv"
)

var (
	ErrInvalidValType = errors.New("ERR Operation against a key holding the wrong kind of value")
)

var Cmds = make(map[string]Cmd)

func Register(name string, c Cmd) {
	if name == "" {
		panic("cmds: call Register with empty command name")
	}
	if _, ok := Cmds[name]; ok {
		panic(fmt.Sprintf("cmds: command %s already registered", name))
	}
	Cmds[name] = Cmd
}

type Cmd interface {
	ParseAndExec([]string) (interface{}, error)
	IntArgIndices() []int
	FloatArgIndices() []int
	NumArgs() (int, int)
}

type DBCmd interface {
	Cmd

	ExecWithDB(srv.DB, []string, []int, []float64) (interface{}, error)
}

var _ DBCmd = (*dbCmd)(nil)

type dbCmd struct {
	fn func(srv.Key, []string, []int, []float64) (interface{}, error)

	noKey            NoKeyFlag
	floatIndices     []int
	intIndices       []int
	minArgs, maxArgs int
}

func (c *dbCmd) IntArgIndices() []int       { return c.intIndices }
func (c *dbCmd) FloatArgIndices() []float64 { return c.floatIndices }
func (c *dbCmd) NumArgs() (int, int)        { return c.minArgs, c.maxArgs }

func (c *dbCmd) ExecWithDB(db srv.DB, args []string, ints []int, floats []float64) (interface{}, error) {
	k, def := c.db.Key(args[0], c.noKey)
	defer def() // Unlock the db

	return c.fn(k, args, ints, floats)
}
