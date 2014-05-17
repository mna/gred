package strings

import (
	"github.com/PuerkitoBio/gred/cmds"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
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

func init() {
	cmds.Register("get", get)
	cmds.Register("set", set)
}

var get = &dbCmd{
	noKey:   srv.NoKeyNone,
	minArgs: 1,
	maxArgs: 1,
	fn:      getFn,
}

func getFn(k srv.Key, args []string, ints []int, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		return v.Get(), nil
	}
	return nil, cmds.ErrInvalidValType
}

var set = &dbCmd{
	noKey:   srv.NoKeyCreateString,
	minArgs: 2,
	maxArgs: 2,
	fn:      setFn,
}

func setFn(k srv.Key, args []string, ints []int, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.String); ok {
		k.Abort()
		v.Set(args[1])
		return nil, nil
	}
	return nil, cmds.ErrInvalidValType
}
