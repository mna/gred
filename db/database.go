package db

import (
	"errors"
	"sync"

	"github.com/PuerkitoBio/gred/resp"
)

var (
	ErrInvalidCommand = errors.New("db: invalid command")
	ErrMissingArg     = errors.New("db: missing argument")
	ErrNilSuccess     = errors.New("db: (nil)")
)

type Database struct {
	ix   int
	mu   sync.RWMutex
	keys map[string]*Key
}

func NewDB(index int) *Database {
	return &Database{
		ix:   index,
		keys: make(map[string]*Key),
	}
}

func (d *Database) Do(cmd string, args ...string) (interface{}, error) {
	switch cmd {
	case "set":
		if len(args) < 2 {
			return nil, ErrMissingArg
		}
		d.Set(args[0], args[1])
		return nil, nil

	case "get":
		if len(args) < 1 {
			return nil, ErrMissingArg
		}
		val, err := d.Get(args[0])
		return resp.BulkString(val), err

	default:
		return nil, ErrInvalidCommand
	}
}

func (d *Database) Get(k string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if key, ok := d.keys[k]; ok {
		return key.Get(), nil
	}
	return "", ErrNilSuccess
}

func (d *Database) Set(k, v string) {
	d.mu.RLock()
	if key, ok := d.keys[k]; !ok {
		// Key does not exist yet, must create the key
		d.mu.RUnlock()
		d.mu.Lock()
		defer d.mu.Unlock()
		key = &Key{val: v}
		d.keys[k] = key

	} else {
		// Key already exists, set the new value
		defer d.mu.RUnlock()
		key.Set(v)
	}
}

type Key struct {
	mu  sync.RWMutex
	val string
}

func (k *Key) Get() string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.val
}

func (k *Key) Set(v string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.val = v
}
