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

// RWLocker is the interface that defines the methods to lock/unlock
// and read-lock/read-unlock.
type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

// Key is the interface that defines the methods to represent a Key.
type Key interface {
	RWLocker
	Expirer
	Name() string
}
