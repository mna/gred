package db

import "sync"

// Static type check that stringKey implements StringKey
var _ StringKey = (*stringKey)(nil)

// StringKey is the interface that represents a key holding
// a string value.
type StringKey interface {
	Key
	Get() string
	Set(string)
}

// stringKey is the internal implementation of a StringKey.
type stringKey struct {
	sync.RWMutex
	Expirer

	name string
	val  string
}

// Name returns the name of the key. It assumes at least a read lock
// is held by the caller.
func (s *stringKey) Name() string {
	return s.name
}

// Get returns the string value of the key. It assumes at least a read lock
// is held by the caller.
func (s *stringKey) Get() string {
	return s.val
}

// Set sets the string value of the key. It assumes a lock is held
// by the caller.
func (s *stringKey) Set(v string) {
	s.val = v
}
