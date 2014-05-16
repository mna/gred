package db

import "sync"

var _ SetKey = (*setKey)(nil)

type set map[string]struct{}

type SetKey interface {
	Key
	Get() set
}

type setKey struct {
	sync.RWMutex
	Expirer

	name string
	s    set
}

func (s *setKey) Name() string {
	return s.name
}

func (s *setKey) Get() set {
	return s.s
}
