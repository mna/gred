package srv

import "sync"

type RWLocker interface {
	sync.Locker

	RLock()
	RUnlock()
}
