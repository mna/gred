package srv

import (
	"fmt"
	"sync"

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
	switch flag {
	case NoKeyCreateString:
		k := newKey(name, vals.NewString())
		d.keys[name] = k
		return k, ret
	default:
		panic(fmt.Sprintf("db.Key NoKeyFlag not implemented: %d", flag))
	}
}
