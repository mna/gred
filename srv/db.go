package srv

import (
	"fmt"
	"sync"
	"time"

	"github.com/PuerkitoBio/gred/vals"
)

// NoKeyFlag indicates what the DB should do when a non-existing key
// is requested.
type NoKeyFlag int

const (
	// NoKeyNone indicates that a nil key should be returned if the requested
	// key does not exist.
	NoKeyNone NoKeyFlag = iota

	// NoKeyDefaultVal indicates that a the default key should be returned if the requested
	// key does not exist. This key implements all value types, returning default values
	// for each.
	NoKeyDefaultVal

	// NoKeyCreateString indicates that a key holding an empty string should be created if the
	// requested key does not exist.
	NoKeyCreateString

	// NoKeyCreateStringInt indicates that a key holding string value of "0" should be created if the
	// requested key does not exist.
	NoKeyCreateStringInt

	// NoKeyCreateHash indicates that a key holding an empty hash should be created if the
	// requested key does not exist.
	NoKeyCreateHash

	// NoKeyCreateList indicates that a key holding an empty list should be created if the
	// requested key does not exist.
	NoKeyCreateList

	// NoKeyCreateSet indicates that a key holding an empty set should be created if the
	// requested key does not exist.
	NoKeyCreateSet

	// NoKeyCreateSortedSet indicates that a key holding an empty sorted set should be created if the
	// requested key does not exist.
	NoKeyCreateSortedSet
)

// WaitChan is the channel type required for the blocking operations on Lists.
type WaitChan <-chan chan<- [2]string

// DB represents a Database, and defines the methods required to manipulate
// its keys.
type DB interface {
	// Sync mutex interface
	RWLocker

	// DB-level commands
	Del(...string) int
	Exists(string) bool
	Expire(string, int64, func()) bool
	ExpireAt(string, int64, func()) bool
	Persist(string) bool
	PExpire(string, int64, func()) bool
	PExpireAt(string, int64, func()) bool
	PSetEx(string, int64, string, func())
	PTTL(string) int64
	SetEx(string, int64, string, func())
	TTL(string) int64
	Type(string) string

	// Keys access
	Keys() map[string]Key
	DelKey(string)
	LockGetKey(string, NoKeyFlag) (Key, func())
	XLockGetKey(string, NoKeyFlag) (Key, func())

	// Blocking list waiters
	WaitLPop(string, WaitChan)
	WaitRPop(string, WaitChan)
}

// Static check to make sure *db implements the DB interface.
var _ DB = (*db)(nil)

// db is the implementation of the DB interface.
type db struct {
	sync.RWMutex

	// the database index
	ix int

	// the keys held by the database
	keys map[string]Key

	// Block list waiters
	waitersChans  map[string][]WaitChan
	waitersPopPos map[string][]bool
}

// NewDB creates a new DB value, with the specified index.
func NewDB(ix int) DB {
	return &db{
		ix:            ix,
		keys:          make(map[string]Key),
		waitersChans:  make(map[string][]WaitChan),
		waitersPopPos: make(map[string][]bool),
	}
}

func (d *db) WaitLPop(key string, ch WaitChan) {
	d.waitPop(key, ch, false)
}

func (d *db) WaitRPop(key string, ch WaitChan) {
	d.waitPop(key, ch, true)
}

func (d *db) waitPop(key string, ch WaitChan, rpop bool) {
	slch := d.waitersChans[key]
	slch = append(slch, ch)
	d.waitersChans[key] = slch
	slbl := d.waitersPopPos[key]
	slbl = append(slbl, rpop)
	d.waitersPopPos[key] = slbl
}

func (d *db) Del(names ...string) int {
	var cnt int
	for _, nm := range names {
		if k, ok := d.keys[nm]; ok {
			k.Lock()
			k.Abort()
			delete(d.keys, nm)
			cnt++
			k.Unlock()
		}
	}
	return cnt
}

func (d *db) Exists(name string) bool {
	_, ok := d.keys[name]
	return ok
}

func (d *db) Expire(name string, secs int64, fn func()) bool {
	return d.expireDuration(name, time.Duration(secs)*time.Second, fn)
}

func (d *db) ExpireAt(name string, uxts int64, fn func()) bool {
	secs := uxts - time.Now().Unix()
	return d.expireDuration(name, time.Duration(secs)*time.Second, fn)
}

func (d *db) PExpire(name string, ms int64, fn func()) bool {
	return d.expireDuration(name, time.Duration(ms)*time.Millisecond, fn)
}

func (d *db) PExpireAt(name string, uxts int64, fn func()) bool {
	dur := (time.Duration(uxts) * time.Millisecond) - time.Duration(time.Now().UnixNano())
	return d.expireDuration(name, dur, fn)
}

func (d *db) expireDuration(name string, dur time.Duration, fn func()) bool {
	if k, ok := d.keys[name]; ok {
		k.Lock()
		defer k.Unlock()
		return k.Expire(dur, fn)
	}
	return false
}

func (d *db) PSetEx(name string, ms int64, v string, fn func()) {
	d.setExDuration(name, time.Duration(ms)*time.Millisecond, v, fn)
}

func (d *db) SetEx(name string, secs int64, v string, fn func()) {
	d.setExDuration(name, time.Duration(secs)*time.Second, v, fn)
}

func (d *db) setExDuration(name string, dur time.Duration, v string, fn func()) {
	// Get or create the key
	k, def := d.LockGetKey(name, NoKeyCreateString)
	defer def()

	// Set its value
	k.Lock()
	defer k.Unlock()
	kv := k.Val().(vals.String)
	kv.Set(v)

	// Expire the key
	k.Expire(dur, fn)
}

func (d *db) Persist(name string) bool {
	if k, ok := d.keys[name]; ok {
		k.Lock()
		defer k.Unlock()
		return k.Abort()
	}
	return false
}

func (d *db) PTTL(name string) int64 {
	if k, ok := d.keys[name]; ok {
		k.RLock()
		defer k.RUnlock()
		ttl := k.TTL()
		if ttl < 0 {
			return int64(ttl)
		}
		return int64(ttl / time.Millisecond)
	}
	return -2
}

func (d *db) TTL(name string) int64 {
	if k, ok := d.keys[name]; ok {
		k.RLock()
		defer k.RUnlock()
		ttl := k.TTL()
		if ttl < 0 {
			return int64(ttl)
		}
		return int64(ttl / time.Second)
	}
	return -2
}

func (d *db) Type(name string) string {
	if k, ok := d.keys[name]; ok {
		k.RLock()
		defer k.RUnlock()
		return k.Val().Type()
	}
	return "none"
}

func (d *db) Keys() map[string]Key {
	return d.keys
}

// DelKey deletes the specified key. It is assumed the caller has an exclusive lock
// for both the DB and the key to delete.
func (d *db) DelKey(name string) {
	k, ok := d.keys[name]
	if ok {
		k.Abort()
		delete(d.keys, name)
	}
}

func (d *db) XLockGetKey(name string, flag NoKeyFlag) (Key, func()) {
	return d.lockGetKey(true, name, flag)
}

func (d *db) LockGetKey(name string, flag NoKeyFlag) (Key, func()) {
	return d.lockGetKey(false, name, flag)
}

func (d *db) lockGetKey(excl bool, name string, flag NoKeyFlag) (Key, func()) {
	var ret func()

	if excl {
		d.Lock()
		ret = d.Unlock
	} else {
		d.RLock()
		ret = d.RUnlock
	}
	if k, ok := d.keys[name]; ok {
		return k, ret
	}

	// Key does not exist, what to do?
	switch flag {
	case NoKeyNone:
		return nil, ret
	case NoKeyDefaultVal:
		return defKey(name), ret
	}

	// Otherwise, upgrade lock if it wasn't already exclusive
	if !excl {
		d.RUnlock()
		d.Lock()
		ret = d.Unlock

		// Check if key now exists (added during the lock upgrade)
		if k, ok := d.keys[name]; ok {
			return k, ret
		}
	}

	// Still no chance, create as requested
	var k Key
	switch flag {
	case NoKeyCreateString:
		k = NewKey(name, vals.NewIncString(""))
	case NoKeyCreateStringInt:
		k = NewKey(name, vals.NewIncString("0"))
	case NoKeyCreateHash:
		k = NewKey(name, vals.NewIncHash())
	case NoKeyCreateList:
		k = NewKey(name, vals.NewList())
	case NoKeyCreateSet:
		k = NewKey(name, vals.NewSet())
	default:
		panic(fmt.Sprintf("db.Key NoKeyFlag not implemented: %d", flag))
	}
	d.keys[name] = k
	return k, ret
}
