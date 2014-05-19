package connection

import "github.com/PuerkitoBio/gred/cmd"

func init() {
	cmd.Register("echo", echo)
	cmd.Register("ping", ping)
	cmd.Register("quit", quit)
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
