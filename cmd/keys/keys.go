package dbcmds

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
)

func init() {
	cmd.Register("del", del)
}

var del = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: -1,
	},
	delFn)

func delFn(db srv.DB, args []string, ints []int, floats []float64) (interface{}, error) {
	db.Lock()
	defer db.Unlock()

	var cnt int64
	for _, nm := range args {
		if db.Del(nm) {
			cnt++
		}
	}
	return cnt, nil
}
