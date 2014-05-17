package cmds

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

var Cmds = make(map[string]Cmd)

func Register(name string, c Cmd) {
	if name == "" {
		panic("cmds: call Register with empty command name")
	}
	if _, ok := Cmds[name]; ok {
		panic(fmt.Sprintf("cmds: command %s already registered", name))
	}
	Cmds[name] = c
}

type Cmd interface {
	IntArgIndices() []int
	FloatArgIndices() []int
	NumArgs() (int, int)
}

type DBCmd interface {
	Cmd

	ExecWithDB(srv.DB, []string, []int, []float64) (interface{}, error)
}
