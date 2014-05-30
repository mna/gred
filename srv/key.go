package srv

import (
	"sync"

	"github.com/PuerkitoBio/gred/vals"
)

// RWLocker defines the methods required to implement a multi-reader,
// single-writer lock.
type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

// Key defines the methods required to implement a database Key.
type Key interface {
	// Read-Write locker
	RWLocker

	// Expirer behaviour
	Expirer

	// Val returns the underlying value
	Val() vals.Value

	// Name returns the name of the key
	Name() string
}

// key implements the Key interface.
type key struct {
	sync.RWMutex
	*expirer

	v    vals.Value
	name string
}

// NewKey creates a new Key with the specified name and value.
func NewKey(name string, v vals.Value) Key {
	return &key{
		expirer: &expirer{},
		v:       v,
		name:    name,
	}
}

// Name returns the name of the key.
func (k *key) Name() string { return k.name }

// Val returns the value of the key.
func (k *key) Val() vals.Value { return k.v }
