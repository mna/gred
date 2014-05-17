package srv

import (
	"sync"

	"github.com/PuerkitoBio/gred/vals"
)

type Key interface {
	RWLocker
	Expirer

	Val() vals.Value
	Name() string
}

type key struct {
	sync.RWMutex
	*expirer

	v    vals.Value
	name string
}

func newKey(name string, v vals.Value) Key {
	return &key{
		expirer: &expirer{},
		v:       v,
		name:    name,
	}
}

func (k *key) Name() string    { return k.name }
func (k *key) Val() vals.Value { return k.v }
