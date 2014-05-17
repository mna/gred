package dbcmd

import (
	"github.com/PuerkitoBio/gred/cmds"
	"github.com/PuerkitoBio/gred/srv"
)

var _ cmds.DBCmd = (*dbCmd)(nil)

type dbCmd struct {
	fn func(srv.Key, []string, []int, []float64) (interface{}, error)

	noKey            srv.NoKeyFlag
	floatIndices     []int
	intIndices       []int
	minArgs, maxArgs int
}

func (c *dbCmd) IntArgIndices() []int   { return c.intIndices }
func (c *dbCmd) FloatArgIndices() []int { return c.floatIndices }
func (c *dbCmd) NumArgs() (int, int)    { return c.minArgs, c.maxArgs }

func (c *dbCmd) ExecWithDB(db srv.DB, args []string, ints []int, floats []float64) (interface{}, error) {
	k, def := db.Key(args[0], c.noKey)
	defer def() // Unlock the db

	return c.fn(k, args, ints, floats)
}
