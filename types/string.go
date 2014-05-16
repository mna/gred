package types

type String interface {
	Append(string) int
	Get() string
	GetRange(int, int) string
	GetSet(string) string
	Set(string)
	StrLen() int
}
