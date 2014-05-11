package db

import (
	"sync"
	"time"
)

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

	name string
	val  string

	// Expiration fields
	tmr   *time.Timer
	expAt time.Time
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

// Expire sets the key to expire at time t, at which point function fn will
// be executed in its own goroutine. It returns true if expiration was
// successfully set, false otherwise.
func (s *stringKey) Expire(t time.Time, fn func()) bool {
	dur := t.Sub(time.Now())
	if s.tmr != nil {
		set := s.tmr.Reset(dur)
		if set {
			s.expAt = t
		}
		return set
	}

	s.tmr = time.AfterFunc(dur, fn)
	s.expAt = t
	return true
}

// Abort aborts execution of an expiration timer. It returns true if a timer
// was successfully aborted, false otherwise.
func (s *stringKey) Abort() bool {
	if s.tmr != nil {
		ok := s.tmr.Stop()
		if ok {
			s.tmr = nil
			s.expAt = time.Time{}
		}
		return ok
	}
	return false
}

// TTL returns the time-to-live of the key before an expiration is triggered.
// It returns -1 if there is no expiration associated with the key.
func (s *stringKey) TTL() time.Duration {
	if s.tmr == nil {
		return -1
	}
	dur := s.expAt.Sub(time.Now())
	if dur > 0 {
		return dur
	}
	return 0
}
