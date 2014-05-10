package db

import (
	"sync"
	"time"
)

type expirer struct {
	db    *Database
	key   string
	setAt time.Time
	dur   time.Duration
	stop  chan struct{}

	mu   sync.RWMutex
	done bool
}

func newExpirer(db *Database, key string, dur time.Duration) *expirer {
	return &expirer{
		db:   db,
		key:  key,
		dur:  dur,
		stop: make(chan struct{}),
	}
}

func (e *expirer) start() {
	e.setAt = time.Now()
	ch := time.After(e.dur)
	go func() {
		select {
		case <-ch:
		case <-e.stop:
		}
		e.mu.Lock()
		defer e.mu.Unlock()
		e.done = true
	}()
}

func (e *expirer) abort() {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if !e.done {
		close(e.stop)
	}
}

func (e *expirer) ttl() time.Duration {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if !e.done {
		return time.Now().Sub(e.setAt)
	}
	return 0
}
