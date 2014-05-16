package vals

type List interface {
	Value

	LIndex(int) (string, bool)
	LInsertBefore(string, string) int
	LInsertAfter(string, string) int
	LLen() int
	LPop() (string, bool)
	LPush(string, ...string) int
	LRange(int, int) []string
	LRem(int, string) int
	LSet(int, string) bool
	RPop() (string, bool)
	RPush(string, ...string) int
}
