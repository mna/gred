package server

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
)

func init() {
	cmd.Register("flushdb", flushdb)
	cmd.Register("flushall", flushall)
}

var flushdb = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 0,
		MaxArgs: 0,
	},
	flushdbFn)

func flushdbFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.Lock()
	defer db.Unlock()

	db.FlushDB()
	return cmd.OKVal, nil
}

var flushall = cmd.NewSrvCmd(
	&cmd.ArgDef{
		MinArgs: 0,
		MaxArgs: 0,
	},
	flushallFn)

func flushallFn(args []string, ints []int64, floats []float64) (interface{}, error) {
	srv.DefaultServer.Lock()
	defer srv.DefaultServer.Unlock()

	srv.DefaultServer.FlushAll()
	return cmd.OKVal, nil
}
