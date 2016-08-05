package sets

import (
	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/types"
)

func init() {
	cmd.Register("sadd", sadd)
	cmd.Register("scard", scard)
	cmd.Register("sdiff", sdiff)
	cmd.Register("sdiffstore", sdiffstore)
	cmd.Register("sismember", sismember)
	cmd.Register("smembers", smembers)
	cmd.Register("srem", srem)
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
	if v, ok := v.(types.Set); ok {
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
	if v, ok := v.(types.Set); ok {
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
	diffSets := make([]types.Set, 0, len(args))
	first := true
	for _, nm := range args {
		// Check if key exists
		if k, ok := keys[nm]; ok {
			// It does, rlock the key
			k.RLock()
			defer k.RUnlock()

			// Get the value, make sure it is a Set
			v := k.Val()
			if v, ok := v.(types.Set); ok {
				diffSets = append(diffSets, v)
			} else {
				return nil, cmd.ErrInvalidValType
			}
		} else if first {
			// If first key does not exist, insert an empty set
			diffSets = append(diffSets, types.NewSet())
		}
		first = false
	}

	return diffSets[0].SDiff(diffSets[1:]...), nil
}

var sdiffstore = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: -1,
	},
	sdiffstoreFn)

func sdiffstoreFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	// In every case, a new key is created at destination, so must have a lock
	db.Lock()
	defer db.Unlock()

	keys := db.Keys()

	diffSets := make([]types.Set, 0, len(args)-1)
	first := true
	for _, nm := range args[1:] {
		// Check if key exists
		if k, ok := keys[nm]; ok {
			// It does, rlock the key
			k.RLock()
			defer k.RUnlock()

			// Get the value, make sure it is a Set
			v := k.Val()
			if v, ok := v.(types.Set); ok {
				diffSets = append(diffSets, v)
			} else {
				return nil, cmd.ErrInvalidValType
			}
		} else if first {
			// If first key does not exist, insert an empty set
			diffSets = append(diffSets, types.NewSet())
		}
		first = false
	}

	val := diffSets[0].SDiff(diffSets[1:]...)
	// If destination exists, remove any expiration and delete
	if dst, ok := keys[args[0]]; ok {
		dst.Lock()
		dst.Abort()
		delete(keys, args[0])
		dst.Unlock()
	}
	// Then create the destination key
	newSet := types.NewSet()
	dst := srv.NewKey(args[0], newSet)
	keys[args[0]] = dst
	return newSet.SAdd(val...), nil
}

var sismember = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	srv.NoKeyDefaultVal,
	sismemberFn)

func sismemberFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Set); ok {
		return v.SIsMember(args[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var smembers = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	smembersFn)

func smembersFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Set); ok {
		return v.SMembers(), nil
	}
	return nil, cmd.ErrInvalidValType
}

var srem = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: -1,
	},
	srv.NoKeyDefaultVal,
	sremFn)

func sremFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Set); ok {
		return v.SRem(args[1:]...), nil
	}
	return nil, cmd.ErrInvalidValType
}
