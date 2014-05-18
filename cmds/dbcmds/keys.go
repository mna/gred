package dbcmds

import "github.com/PuerkitoBio/gred/srv"

var del = &dbCmd{
	noKey:   srv.NoKeyNone,
	minArgs: 1,
	maxArgs: -1,
	fn:      delFn,
}

func delFn(db srv.DB, args []string, ints []int, floats []float64) (interface{}, error) {
	db.Lock()
	defer db.Unlock()

}
