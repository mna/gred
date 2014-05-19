package srv

import (
	"fmt"
	"sync"
	"time"

	"github.com/PuerkitoBio/gred/vals"
	"github.com/golang/glog"
)

type NoKeyFlag int

const (
	NoKeyNone NoKeyFlag = iota
	NoKeyDefaultVal
	NoKeyCreateString
	NoKeyCreateHash
	NoKeyCreateList
	NoKeyCreateSet
	NoKeyCreateSortedSet
)

type DB interface {
	RWLocker

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

	Key(string) Key
	LockGetKey(string, NoKeyFlag) (Key, func())
}

type db struct {
	sync.RWMutex

	ix   int
	keys map[string]Key
}

func NewDB(ix int) DB {
	return &db{
		ix:   ix,
		keys: make(map[string]Key),
	}
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

func (d *db) Key(name string) Key {
	return d.keys[name]
}

func (d *db) LockGetKey(name string, flag NoKeyFlag) (Key, func()) {
	d.RLock()
	ret := d.RUnlock
	if k, ok := d.keys[name]; ok {
		glog.V(2).Infof("db %d: found key %s", d.ix, name)
		return k, ret
	}

	glog.V(2).Infof("db %d: key %s does not exist", d.ix, name)
	// Key does not exist, what to do?
	switch flag {
	case NoKeyNone:
		return nil, ret
	case NoKeyDefaultVal:
		return defKey(name), ret
	}

	// Otherwise, upgrade lock
	d.RUnlock()
	d.Lock()
	ret = d.Unlock
	var k Key
	switch flag {
	case NoKeyCreateString:
		k = newKey(name, vals.NewString())
	case NoKeyCreateHash:
		k = newKey(name, vals.NewHash())
	default:
		panic(fmt.Sprintf("db.Key NoKeyFlag not implemented: %d", flag))
	}
	d.keys[name] = k
	return k, ret
}
