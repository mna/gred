package db

import "sync"

var _ ListKey = (*listKey)(nil)

const defaultListCap int = 10

type list []string

func (l *list) lpush(vals ...string) {
	*l = append(vals, *l...)
}

func (l *list) rpush(vals ...string) {
	*l = append(*l, vals...)
}

func (l *list) lpop() (string, bool) {
	if len(*l) == 0 {
		return "", false
	}
	val := (*l)[0]
	*l = (*l)[1:]
	return val, true
}

func (l *list) rpop() (string, bool) {
	if len(*l) == 0 {
		return "", false
	}
	val := (*l)[len(*l)-1]
	*l = (*l)[:len(*l)-1]
	return val, true
}

type ListKey interface {
	Key
	Get() *list
}

type listKey struct {
	sync.RWMutex
	Expirer

	name string
	l    *list
}

func (l *listKey) Name() string {
	return l.name
}

func (l *listKey) Get() *list {
	return l.l
}
