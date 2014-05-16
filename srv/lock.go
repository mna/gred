package srv

type RWLocker interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
	RequireExclusive()
}
