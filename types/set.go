package types

type Set interface {
	SAdd(...string) int
	SCard() int
	SDiff(Set, ...Set) []string
	SInter(Set, ...Set) []string
	SIsMember(string) bool
	SMembers() []string
	SMove(Set, string) bool
	SPop() (string, bool)
	SRem(string, ...string) int
	SUnion(Set, ...Set) []string

	// Private method to get the internal map
	get() set
}

type set map[string]struct{}
