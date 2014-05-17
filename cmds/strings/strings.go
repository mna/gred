package strings

import (
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
)

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
	return nil, ErrInvalidValType
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
	return nil, ErrInvalidValType
}
