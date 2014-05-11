package db

import "sync"

var _ HashKey = (*hashKey)(nil)

type hash map[string]string

type HashKey interface {
	Key
	Get() hash
}

type hashKey struct {
	sync.RWMutex
	Expirer

	name string
	h    hash
}

func (h *hashKey) Name() string {
	return h.name
}

func (h *hashKey) Get() hash {
	return h.h
}
