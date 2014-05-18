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
	if k, ok := d.keys[name]; ok {
		k.Lock()
		defer k.Unlock()
		return k.Expire(time.Now().Add(time.Duration(secs)*time.Second), fn)
	}
	return false
}

func (d *db) ExpireAt(name string, uxts int64, fn func()) bool {
	secs := uxts - time.Now().Unix()
	return d.Expire(name, secs, fn)
}

func (d *db) PExpire(name string, ms int64, fn func()) bool {

}

func (d *db) PExpireAt(name string, uxts int64, fn func()) bool {

}

func (d *db) Persist(name string) bool {
	if k, ok := d.keys[name]; ok {
		k.Lock()
		defer k.Unlock()
		return k.Abort()
	}
	return false
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
	default:
		panic(fmt.Sprintf("db.Key NoKeyFlag not implemented: %d", flag))
	}
	d.keys[name] = k
	return k, ret
}
