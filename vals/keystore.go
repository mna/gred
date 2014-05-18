package vals

// KeyStore is the value that holds keys in a database. It is not
// a vals.Value, so it cannot be the value of a key.
type KeyStore interface {
	Del(string, ...string) int
	Exists(string) bool
	Expire(string, int) bool
	ExpireAt(string, int64) bool
	Keys(string) []string
	Persist(string) bool
	PExpire(string, int) bool
	PExpireAt(string, int64) bool
	PTTL(string) int64
	RandomKey() (string, bool)
	Rename(string, string) error
	RenameNX(string, string) (bool, error)
	TTL(string) int64
	Type(string) string
}
