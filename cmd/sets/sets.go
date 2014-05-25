package sets

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/vals"
)

func init() {
	cmd.Register("sadd", sadd)
	cmd.Register("scard", scard)
	cmd.Register("sdiff", sdiff)
}

var sadd = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: -1,
	},
	srv.NoKeyCreateSet,
	saddFn)

func saddFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(vals.Set); ok {
		return v.SAdd(args[1:]...), nil
	}
	return nil, cmd.ErrInvalidValType
}

var scard = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	scardFn)

func scardFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(vals.Set); ok {
		return v.SCard(), nil
	}
	return nil, cmd.ErrInvalidValType
}

var sdiff = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: -1,
	},
	sdiffFn)

func sdiffFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	db.RLock()
	defer db.RUnlock()

	// Get and rlock all keys
	keys := db.Keys()
	diffSets := make([]vals.Set, 0, len(args))
	first := true
	for _, nm := range args {
		// Check if key exists
		if k, ok := keys[nm]; ok {
			// It does, rlock the key
			k.RLock()
			defer k.RUnlock()

			// Get the value, make sure it is a Set
			v := k.Val()
			if v, ok := v.(vals.Set); ok {
				diffSets = append(diffSets, v)
			} else {
				return nil, cmd.ErrInvalidValType
			}
		} else if first {
			// If first key does not exist, insert an empty set
			diffSets = append(diffSets, vals.NewSet())
		}
		first = false
	}

	return diffSets[0].SDiff(diffSets[1:]...), nil
}
