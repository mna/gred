package db

import (
	"errors"
	"strings"
	"sync"
)

var ErrInvalidCommand = errors.New("db: invalid command")

type Database struct {
	ix   int
	mu   sync.RWMutex
	keys map[string]*Key
}

func New(index int) *Database {
	return &Database{
		ix:   index,
		keys: make(map[string]*Key),
	}
}

func (d *Database) Do(cmd string, args ...string) (interface{}, error) {
	switch strings.ToLower(cmd) {
	case "set":
		d.mu.Lock()
		defer d.mu.Unlock()
		d.keys[args[0]].Set(args[1])
		return nil, nil
	case "get":
		d.mu.RLock()
		defer d.mu.RUnlock()
		return d.keys[args[0]].Get(), nil
	default:
		return nil, ErrInvalidCommand
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
