package db

import "sync"

// Database represents a Redis database, identified by its index.
type Database struct {
	ix   int
	mu   sync.RWMutex
	keys map[string]Key
}

// NewDB creates a new Database identified by the specified index.
func NewDB(index int) *Database {
	return &Database{
		ix:   index,
		keys: make(map[string]Key),
	}
}

type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

type Key interface {
	RWLocker
	Name() string
}

type StringKey interface {
	Key
	Get() string
	Set(string)
}
