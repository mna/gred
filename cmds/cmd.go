package cmds

import "github.com/PuerkitoBio/gred/srv"

type NoKeyFlag int

const (
	NoKeyNone NoKeyFlag = iota
	NoKeyCreateString
	NoKeyCreateHash
	NoKeyCreateList
	NoKeyCreateSet
	NoKeyCreateSortedSet
)

type Cmd interface {
	IntArgIndices() []int
	FloatArgIndices() []int
	NumArgs() (int, int)
	Do([]string, []int, []float64) (interface{}, error)
}

type dbCmd struct {
	db srv.DB

	noKey            NoKeyFlag
	floatIndices     []int
	intIndices       []int
	minArgs, maxArgs int
}

func (c *cmd) Do(args []string, ints []int, floats []float64) (interface{}, error) {

	k := c.db.Key(args[0])
}
