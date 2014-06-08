package connection

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
)

func init() {
	cmd.Register("echo", echo)
	cmd.Register("ping", ping)
	cmd.Register("quit", quit)
	cmd.Register("select", selct)
}

var echo = cmd.NewSrvCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	echoFn)

func echoFn(args []string, ints []int64, floats []float64) (interface{}, error) {
	return args[0], nil
}

var ping = cmd.NewSrvCmd(
	&cmd.ArgDef{},
	pingFn)

func pingFn(args []string, ints []int64, floats []float64) (interface{}, error) {
	return cmd.PongVal, nil
}

var quit = cmd.NewSrvCmd(
	&cmd.ArgDef{},
	quitFn)

func quitFn(args []string, ints []int64, floats []float64) (interface{}, error) {
	return nil, cmd.ErrQuit
}

var selct = cmd.NewConnCmd(
	&cmd.ArgDef{
		MinArgs:    1,
		MaxArgs:    1,
		IntIndices: []int{0},
	},
	selctFn)

func selctFn(conn srv.Conn, args []string, ints []int64, floats []float64) (interface{}, error) {
	srv.DefaultServer.RLock()
	defer srv.DefaultServer.RUnlock()

	_, ok := srv.DefaultServer.GetDB(int(ints[0]))
	if !ok {
		return nil, cmd.ErrInvalidDBIndex
	}
	conn.Select(int(ints[0]))
	return cmd.OKVal, nil
}
