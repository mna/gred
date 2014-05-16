package vals

type Hash interface {
	Value

	HDel(string, ...string) int
	HExists(string) bool
	HGet(string) (string, bool)
	HGetAll() []string
	HKeys() []string
	HLen() int
	HMGet(string, ...string) []string
	HSet(string, string) bool
	HVals() []string
}
