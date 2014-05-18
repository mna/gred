package dbcmd

import (
	"github.com/PuerkitoBio/gred/cmds"
	"github.com/PuerkitoBio/gred/srv"
)

var _ cmds.DBCmd = (*dbCmd)(nil)

type dbCmd struct {
	fn func(srv.DB, []string, []int, []float64) (interface{}, error)

	noKey            srv.NoKeyFlag
	floatIndices     []int
	intIndices       []int
	minArgs, maxArgs int
}

func (c *dbCmd) IntArgIndices() []int   { return c.intIndices }
func (c *dbCmd) FloatArgIndices() []int { return c.floatIndices }
func (c *dbCmd) NumArgs() (int, int)    { return c.minArgs, c.maxArgs }

func (c *dbCmd) ExecWithDB(db srv.DB, args []string, ints []int, floats []float64) (interface{}, error) {
	return c.fn(db, args, ints, floats)
}
