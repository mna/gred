package db

import "sync"

// Static type-check that the hashKey implements HashKey.
var _ HashKey = (*hashKey)(nil)

// hash is the internal type of the hash key's value.
type hash map[string]string

// HashKey is the interface that represents a key holding a
// hash.
type HashKey interface {
	Key
	Get() hash
}

// hashKey is the internal implementation of a hash key.
type hashKey struct {
	sync.RWMutex
	Expirer

	name string
	h    hash
}

// Name returns the name of the key. It assumes at least a read-lock
// is held by the caller.
func (h *hashKey) Name() string {
	return h.name
}

// Get returns the hash map of the key. It assumes at least a read-lock
// is held by the caller.
func (h *hashKey) Get() hash {
	return h.h
}
