package db

import "time"

// Expirer is the interface that defines the methods to manage expiration.
type Expirer interface {
	Expire(time.Time, func()) bool
	TTL() time.Duration
	Abort() bool
}

// expirer is the internal implementation of an Expirer for a Key.
type expirer struct {
	tmr   *time.Timer
	expAt time.Time
}

// Expire sets the key to expire at time t, at which point function fn will
// be executed in its own goroutine. It returns true if expiration was
// successfully set, false otherwise.
func (e *expirer) Expire(t time.Time, fn func()) bool {
	dur := t.Sub(time.Now())
	if e.tmr != nil {
		set := e.tmr.Reset(dur)
		if set {
			e.expAt = t
		}
		return set
	}

	e.tmr = time.AfterFunc(dur, fn)
	e.expAt = t
	return true
}

// Abort aborts execution of an expiration timer. It returns true if a timer
// was successfully aborted, false otherwise.
func (e *expirer) Abort() bool {
	if e.tmr != nil {
		ok := e.tmr.Stop()
		if ok {
			e.tmr = nil
			e.expAt = time.Time{}
		}
		return ok
	}
	return false
}

// TTL returns the time-to-live of the key before an expiration is triggered.
// It returns -1 if there is no expiration associated with the key.
func (e *expirer) TTL() time.Duration {
	if e.tmr == nil {
		return -1
	}
	dur := e.expAt.Sub(time.Now())
	if dur > 0 {
		return dur
	}
	return 0
}
