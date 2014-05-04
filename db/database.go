package db

import (
	"errors"
	"sync"

	"github.com/PuerkitoBio/gred/resp"
)

var ErrInvalidCommand = errors.New("db: invalid command")

type Database struct {
	ix   int
	mu   sync.RWMutex
	keys map[resp.BulkString]*Key
}

func New(index int) *Database {
	return &Database{
		ix:   index,
		keys: make(map[resp.BulkString]*Key),
	}
}

func (d *Database) Do(cmd resp.BulkString, args ...interface{}) (interface{}, error) {
	switch string(cmd) {
	case "set":
		d.mu.Lock()
		defer d.mu.Unlock()
		d.keys[args[0].(resp.BulkString)].Set(args[1].(resp.BulkString))
		return nil, nil
	case "get":
		d.mu.RLock()
		defer d.mu.RUnlock()
		return d.keys[args[0].(resp.BulkString)].Get(), nil
	default:
		return nil, ErrInvalidCommand
	}
}

type Key struct {
	mu  sync.RWMutex
	val resp.BulkString
}

func (k *Key) Get() resp.BulkString {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.val
}

func (k *Key) Set(v resp.BulkString) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.val = v
}
