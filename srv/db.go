package srv

type NoKeyFlag int

const (
	NoKeyNone NoKeyFlag = iota
	NoKeyCreateString
	NoKeyCreateHash
	NoKeyCreateList
	NoKeyCreateSet
	NoKeyCreateSortedSet
)

type DB interface {
	RWLocker

	Key(string, NoKeyFlag) (Key, func())
}
